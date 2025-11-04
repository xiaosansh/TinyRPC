package main

import (
    "fmt"
    "io"
    "log"
    "net"
)

// 每个连接都会启动一个协程来处理
func handleConn(conn net.Conn) {
    defer conn.Close()
    codec := NewJSONCodec(conn)

    for {
        var req Request
        if err := codec.ReadRequest(&req); err != nil {
            if err == io.EOF {
                log.Println("连接关闭")
                return
            }
            log.Println("读取请求出错：", err)
            return
        }

        // 这里是模拟服务逻辑
        log.Printf("收到请求: %#v\n", req)

        // 构造响应
        resp := Response{
            Seq:   req.Seq,
            Reply: fmt.Sprintf("服务端收到: %v", req.Argv),
        }

        if err := codec.WriteResponse(&resp); err != nil {
            log.Println("发送响应出错：", err)
            return
        }
    }
}

func main() {
    listener, err := net.Listen("tcp", ":9999")
    if err != nil {
        log.Fatal("监听失败：", err)
    }
    log.Println("TinyRPC 服务端启动，监听端口 9999")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("连接失败：", err)
            continue
        }
        go handleConn(conn)
    }
}
