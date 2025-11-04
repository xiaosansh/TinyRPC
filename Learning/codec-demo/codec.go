package main

import(
  "encoding/json"
  "io"
)

type Request struct{
  ServiceMethod string
  Seq uint64
  Argv any
}

type Response struct{
  Seq uint64
  Error string
  Reply any
}

type Codec interface{
  ReadRequest(*Request) error
  WriteResponse(*Response) error
  Close() error
}

type JSONCodec struct{
  conn io.ReadWriteCloser
  dec *json.Decoder
  enc *json.Encoder
}

func NewJSONCodec(conn io.ReadWriteCloser) *JSONCodec {
  return &JSONCodec{
    conn: conn,
    dec: json.NewDecoder(conn),
    enc: json.NewEncoder(conn),
  }
}

func(j *JSONCodec) ReadRequest(req *Request) error{
  return j.dec.Decode(&req)
}

func(j *JSONCodec) WriteResponse(resp *Response) error{
  return j.enc.Encode(resp)
}

func(j *JSONCodec) Close() error{
  return j.conn.Close()
}
