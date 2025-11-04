package codec

import (
    "encoding/gob"
    "io"
    "log"
)

type GobCodec struct {
    conn io.ReadWriteCloser
    dec  *gob.Decoder
    enc  *gob.Encoder
  }

func NewGobCodec(conn io.ReadWriteCloser) Codec {
    return &GobCodec{
        conn: conn,
        dec:  gob.NewDecoder(conn),
        enc:  gob.NewEncoder(conn),
    }
}

func (c *GobCodec) ReadHeader(h *Header) error {
    return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body any) error {
    return c.dec.Decode(body)
  }

func (c *GobCodec) Write(h *Header, body any) (err error) {
    defer func() {
        if e := recover(); e != nil {
            log.Printf("rpc: panic in Write: %v", e)
        }
    }()
    if err = c.enc.Encode(h); err != nil {
        return err
    }
    return c.enc.Encode(body)
}

func (c *GobCodec) Close() error {
    return c.conn.Close()
}
