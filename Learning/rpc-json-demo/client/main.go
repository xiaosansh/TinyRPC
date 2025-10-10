package main

import(
  "fmt"
  "encoding/json"
  "net"
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
  conn,err := net.Dial("tcp","localhost:1234")
  if err != nil {
    panic(err)
  }

  decoder := json.NewDecoder(conn)
  encoder := json.NewEncoder(conn)

  req := Request{
    Method: "GetTime",
    Params: []interface{}{},
  }

  
  if err := encoder.Encode(req); err != nil {
    panic(err)
  }

  var resp Response
  if err := decoder.Decode(&resp); err != nil{
    panic(err)
  }

  fmt.Println("Result:",resp.Result)
  if resp.Error != "" {
        fmt.Println("Error:", resp.Error)
    }
}
