package codec

import(
  "io"
)

type Header struct{
  ServiceMethod string
  Seq uint64
  Error any
}
//Codec 是所有编解码器必须实现的接口
type Codec interface{
  io.Close
  ReadHeader(*Header) error
  ReadBody(any) error
  Write(*Header,any) error
}
//定义Codec的构造函数签名
type NewCodecFunc funcf(io.ReadWriteClose) Codec

//使用字符串区分不同的codec
type CodecType string

const(
  GobType CodecType = "gob"
  JSONType CodecType = "json"
)

//注册表
var NewCodecFuncMap map[CodecType]NewCodecFunc

func init(){
  NewCodecFuncMap := make(map[CodecType]NewCodecFunc)
  NewCodecFuncMap[GobType] = NewGobCodec
  NewCodecFuncMap[JSONType] = NewJSONCodec
}
