package codec

import(
  "encoding/json"
  "log"
  "io"
)

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

func(j *JSONCodec) ReadHeader(h *Header) error{
  return j.dec.Decode(&h)
}

func(j *JSONCodec) ReadBody(body any) error{
  return j.dec.Decode(body)
}

func(j *JSONCodec) Write(h *Header,body any) (err error){
  defer func(){
    if err := recover();err != nil{
      log.Printf("rpc:panic in Write %v",err)
    }()
    if err = j.enc.Encode(h);err != nil{
      return err
    }
    return j.enc.Encode(body)
  }
}


func(j *JSONCodec) Close() error{
  return j.conn.Close()
}
