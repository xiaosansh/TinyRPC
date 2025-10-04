package main

import(
  "time"
  "log"
  "net"
  "net/http"
  "net/rpc"
)

type Args struct{}
type TimeServer int64

func (t *TimeServer) GiveServerTime(args *Args,reply *int64) error{
  //通过改变指针指向内容来实现效果，而不是return
  *reply = time.Now().Unix()
  //返回nil表示没有错误
  return nil
}

func main(){
  //1.创建一个新的RPC服务器实例
  timeserver := new(TimeServer)

  //2.注册这个实例
  rpc.Register(timeserver)

  //3.将RPC服务器绑定到HTTP协议上
  rpc.HandleHTTP()

  //4.监听TCP端口
  l,e := net.Listen("tcp",":2233")
  if e != nil {
    log.Fatal("Listen error:",e)
  }

  //5.启动HTTP服务，开始接收请求
  http.Serve(l,nil)
}
