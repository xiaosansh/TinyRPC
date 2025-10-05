package main

import(
  "fmt"
  "net"
)

func main(){
  //在本地启动一个TCP服务器，监听1234端口
  ln,err := net.Listen("tcp",":1234")
  if err != nil {
    panic(err)
  }
  fmt.Println("server is Listening on port 1234...")

  for {
  conn,err := ln.Accept()
  if err != nil {
    fmt.Println("Error acception connection:",err)
    continue
  }

  go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  defer conn.Close()

  buf := make([]byte,1024)
  n,err := conn.Read(buf)
  if err != nil {
    fmt.Println("Error Reading:",err)
    return
  }

  message := string(buf[:n])
  fmt.Println("Received from client:",message)

  response := "Hello from Server"
  conn.Write([]byte(response))
}
