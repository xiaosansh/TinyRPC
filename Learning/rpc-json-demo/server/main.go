package main

import(
  "fmt"
  "encoding/json"
  "net"
  "time"
)

type Request struct{
  Method string `json:"method"`
  Params []interface{} `json:"params"`
}

type Response struct{
  Result interface{} `json:"result"`
  Error string `json:"error"`
}

func main(){
  l,e := net.Listen("tcp",":1234")
  if e != nil {
    panic(e)
  }
  fmt.Println("Reserve is Linstening on 1234...")

  for{
    conn,err := l.Accept()
    if err != nil {
      fmt.Println("Error acception:",err)
      continue
    }

    go handlConn(conn)
  }
}

func handlConn(conn net.Conn){
  defer conn.Close()

  decoder := json.NewDecoder(conn)
  encoder := json.NewEncoder(conn)

  var req Request
  if err := decoder.Decode(&req);err != nil {
    fmt.Println("Decode error:",err)
    return
  }

 fmt.Println("Received request:", req.Method)

  var resp Response
  switch req.Method {
  case "GetTime":
    resp.Result =  time.Now().Format(time.RFC3339)
  default:
    resp.Error = "Unknown method: " + req.Method
  }

  if err := encoder.Encode(resp);err != nil {
    fmt.Println("Encode error:",err)
  }
}
