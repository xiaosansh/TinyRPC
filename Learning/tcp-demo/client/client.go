package main

import(
  "net"
  "fmt"
)

func main(){
  conn,err := net.Dial("tcp","localhost:1234")
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  message := "Hello"
  conn.Write([]byte(message))
  fmt.Println("Sent to server:",message)

  buf := make([]byte,1024)
  n,err := conn.Read(buf)
  fmt.Println("Received from server:",string(buf[:n]))
}
