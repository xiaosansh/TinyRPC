package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:9999")
    if err != nil {
        log.Fatal("连接失败：", err)
    }
    codec := NewJSONCodec(conn)

    // 构造请求
    req := Request{
        ServiceMethod: "DemoService.Hello",
        Seq:           1,
        Argv:          "Hello TinyRPC",
    }

    // 发送请求
    if err := codec.enc.Encode(&req); err != nil {
        log.Fatal("发送请求失败：", err)
    }

    // 读取响应
    var resp Response
    if err := codec.dec.Decode(&resp); err != nil {
        log.Fatal("读取响应失败：", err)
    }

    fmt.Println("收到响应:", resp.Reply)
}

