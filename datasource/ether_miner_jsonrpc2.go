package datasource

//stripped down and lenient jsonrpc2 client supporting
//ethminer and claymore (which does not support standard
//jsonrpc2)

//heavily inspired by github.com/powerman/rpc-codec/jsonrpc2

import (
	"io"
	"net"
	"net/rpc"
	"encoding/json"
	"errors"
	"fmt"
)

type clientCodec struct {
	dec  *json.Decoder
	enc  *json.Encoder
	c    io.Closer
	req  clientRequest
	resp clientResponse
}

func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{
		dec:     json.NewDecoder(conn),
		enc:     json.NewEncoder(conn),
		c:       conn,
		req:	clientRequest{"2.0", "", nil},
	}
}

//client request supporting no parameters (not needed by miners)
type clientRequest struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      *uint64     `json:"id"`
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) error {
	c.req.Method = r.ServiceMethod
	c.req.ID = &r.Seq
	return c.enc.Encode(&c.req)
}

type clientResponse struct {
	//allow omission of Version and ID for claymore
	Version string                      `json:"jsonrpc,omitempty"`
	ID      *uint64                     `json:"id,omitempty"`
	Result  *json.RawMessage            `json:"result,omitempty"`
	Error   map[string]*json.RawMessage `json:"error,omitempty"`
}

func (r *clientResponse) reset() {
	r.Version = ""
	r.ID = nil
	r.Result = nil
	r.Error = nil
}

func (r *clientResponse) UnmarshalJSON(raw []byte) error {
	r.reset()
	type resp *clientResponse
	return json.Unmarshal(raw, resp(r))
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	if err := c.dec.Decode(&c.resp); err != nil {
		return err
	}

	if c.resp.ID == nil {
		return errors.New("id is nil: " + fmt.Sprintf("%v", c.resp.Error))
	}

	r.Error = ""
	r.Seq = *c.resp.ID
	if c.resp.Error != nil {
		r.Error = fmt.Sprintf("%v", c.resp.Error)
	}
	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	return json.Unmarshal(*c.resp.Result, x)
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}

func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn))
}

func Dial(network, address string) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), err
}
