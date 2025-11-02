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


