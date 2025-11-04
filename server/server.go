package server

import(
  "encoding/json"
  "io"
  "fmt"
  "log"
  "net"
  "reflect"
  "sync"
  "Tinyrpc/codec"
)

//MagicNumber用于标识TinyRPC协议
const MagicNumber = 0x3bef5c

//Option用于客户端和服务端协商通信方式
type Option struct{
  MagicNumber int
  CodecType string //e.g. gob json
}

//Server表示一个RPC服务器实例
type Server struct{
  serviceMap sync.Map //储存注册的服务
}

func NewServer() *Server{
  return &Server{}
}

//Accept等待并并发处理所有进入的连接
func (s *Server) Accept(lis net.Listener){
  for{
    conn,err := lis.Accept()
    if err != nil{
      log.Println("rpc servevr:accept error:",err)
      return
    }
    go s.ServerConn(conn)
  }
}

//处理单个连接
func (s *Server) ServerConn(conn io.ReadWriteCloser){
  defer conn.Close()

  var opt Option
  if err := json.NewDecoder(conn).Decode(&opt); err != nil{
    log.Println("rpc server: decode option error:",err)
    return
  }

  //根据Option决定使用哪种CodecType
  f := codec.NewCodecFuncMap[opt.CodecType]
  if f == nil{
    log.Printf("rpc server: invalid codec type %s\n",opt.CodecType)
    return
  }
  s.serverCodec(f(conn))
}

//并发处理
func (s *Server) serverCodec(cc codec.Codec){
  sending := new(sync.Mutex)
  wg := new(sync.WaitGroup)

  for{
    req,err := s.readRequest(cc)
    if err != nil{
      if err == io.EOF{
        break
      }
      log.Println("rpc server:read request error:",err)
      break
    }
    wg.Add(1)
    go s.handleRequest(cc,req,sending,wg)
  }

  wg.Wait()
  _ = cc.Close()
}

type request struct {
  h            *codec.Header  //请求头部
  argv,replyv  reflect.Value
  svc          *service
  mtype        *methodType
}

func (s *Server) readRequest(cc codec.Codec) (*request,error){
  h,err := cc.ReadHeader() //1.读取头部
  if err != nil {
    return nil,err
  }

  req := &request{h:h} //2.创建上下文
  req.svc,req.mtype,err = s.findService(h.ServiceMethod) //3.查找服务
  if err != nil{
    return req,err
  }

  req.argv = req.mtype.newArgv() //4.创建参数容器
  req.replyv = req.mtype.newReplyv() //5.创建响应容器

  //6.准备解码
  argvi := req.argv.Interface()
  if req.argv.Kind() != reflect.Ptr{
    argvi = req.argv.Addr().Interface()
  }

  //7.读取体部
  if err = cc.ReadBody(argvi);err != nil{
    log.Println("rpc server:read body error:",err)
    return req,err
  }

  return req,nil
}

//findService解析请求中的service.method并且定位已经注册的服务
func (s *Server) findService(ServiceMethod string) (*service,*methodType,error){
  //分别获取服务名和方法名
  serviceName := serviceMethod[:dot]
  methodName := serviceMethod[dot+1:]

  svcInterface,ok := s.serviceMap.Load(serviceName) //读取sync.Map
  if !ok {
    return nil,nil,fmt.Errorf("rpc server:can't find server %s",serviceName)
  }

  svc := svcInterface.(*service) //类型断言
  mtype := svc.method[methodName] //读取普通map，获取mtype，因为后续注册后不被更改所以并发读取没什么问题
  if mtype == nil {
    return nil,nil,fmt.Errorf("rpc server:can't find method %s", methodName)
  }

  return svc,mtype,nil
}

//handleRequest 调用目标方法并返回结果
func (s *Server) handleRequest(cc codec.Codec,req *request,sending *sync.Mutex,wg *sync.WaitGroup){
  defer wg.Done()  //确保WaitGroup的计数器在减少

  err := req.svc.call(req.mtype,req.argv,rep.replyv) //执行真正的业务逻辑
  if err != nil{
    req.h.Error = err.Error()                            //错误处理
    s.sendResponse(cc,req.h,invalidRequest{},sending)
    return
  }

  s.sendResponse(cc,req.h,req.replyv.Interface(),sending) //成功处理
}

//sendResponse发送响应给客户端
func (s *Server) sendResponse(cc codec.Codec,h *codec.Header,body interface{},sending *sync.Mutex){
  //确保同一时刻只有一个goroutine能写入数据
  sending.Lock()
  defer sending.Unlock()

  if err := cc.Write(h,body); err != nil{
    log.Println("rpc server: write response error:", err)
  }
}

type invalidRequest struct{}
