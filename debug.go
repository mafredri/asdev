package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/mafredri/cdp/rpcc"
)

func newLogCodec(prefix string) rpcc.DialOption {
	logger := log.New(os.Stdout, fmt.Sprintf("rpcc(%s) ", prefix), log.LstdFlags)
	return rpcc.WithCodec(func(conn io.ReadWriter) rpcc.Codec {
		return &rpccLogCodec{conn: conn, log: logger}
	})
}

// rpccLogCodec captures the output from writing RPC requests and reading
// responses on the connection. It implements rpcc.Codec via
// WriteRequest and ReadResponse.
type rpccLogCodec struct {
	conn io.ReadWriter
	log  *log.Logger
}

// WriteRequest marshals v into a buffer, writes its contents onto the
// connection and logs it.
func (c *rpccLogCodec) WriteRequest(req *rpcc.Request) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return err
	}

	c.log.Println("=>", "rpc_id="+strconv.Itoa(int(req.ID)), "rpc_method="+req.Method, "data="+buf.String())

	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// ReadResponse unmarshals from the connection into v whilst echoing
// what is read into a buffer for logging.
func (c *rpccLogCodec) ReadResponse(resp *rpcc.Response) error {
	var buf bytes.Buffer
	if err := json.NewDecoder(io.TeeReader(c.conn, &buf)).Decode(resp); err != nil {
		return err
	}
	c.log.Println("<=", "rpc_id="+strconv.Itoa(int(resp.ID)), "rpc_event="+strconv.FormatBool(resp.Method != ""), "data="+buf.String())

	return nil
}
