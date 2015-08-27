package server

import (
	"bytes"
	"io"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/golang/protobuf/proto"
	"github.com/soundcloud/doozerd/store"
)

var (
	fooPath = "/foo"
)

type bchan chan []byte

func (b bchan) Write(buf []byte) (int, error) {
	b <- buf
	return len(buf), nil
}

func (b bchan) Read(buf []byte) (int, error) {
	return 0, io.EOF // not implemented
}

func mustUnmarshal(b []byte) (r *Response) {
	r = new(Response)
	err := proto.Unmarshal(b, r)
	if err != nil {
		panic(err)
	}
	return
}

func assertResponseErrCode(t *testing.T, exp Response_Err, c *conn) {
	b := c.c.(*bytes.Buffer).Bytes()
	assert.T(t, len(b) > 4, b)
	assert.Equal(t, &exp, mustUnmarshal(b[4:]).ErrCode)
}

func TestDelNilFields(t *testing.T) {
	c := &conn{
		c:        &bytes.Buffer{},
		canWrite: true,
		waccess:  true,
	}
	tx := &txn{
		c:   c,
		req: Request{Tag: proto.Int32(1)},
	}
	tx.del()
	assertResponseErrCode(t, Response_MISSING_ARG, c)
}

func TestSetNilFields(t *testing.T) {
	c := &conn{
		c:        &bytes.Buffer{},
		canWrite: true,
		waccess:  true,
	}
	tx := &txn{
		c:   c,
		req: Request{Tag: proto.Int32(1)},
	}
	tx.set()
	assertResponseErrCode(t, Response_MISSING_ARG, c)
}

func TestServerNoAccess(t *testing.T) {
	b := make(bchan, 2)
	c := &conn{
		c:        b,
		canWrite: true,
		st:       store.New(store.DefaultInitialRev),
	}
	tx := &txn{
		c:   c,
		req: Request{Tag: proto.Int32(1)},
	}

	for i, op := range ops {
		if i != int32(Request_ACCESS) {
			op(tx)
			var exp Response_Err = Response_OTHER
			assert.Equal(t, 4, len(<-b), Request_Verb_name[i])
			assert.Equal(t, &exp, mustUnmarshal(<-b).ErrCode, Request_Verb_name[i])
		}
	}
}

func TestServerRo(t *testing.T) {
	b := make(bchan, 2)
	c := &conn{
		c:        b,
		canWrite: true,
		st:       store.New(store.DefaultInitialRev),
	}
	tx := &txn{
		c:   c,
		req: Request{Tag: proto.Int32(1)},
	}

	wops := []int32{int32(Request_DEL), int32(Request_NOP), int32(Request_SET)}

	for _, i := range wops {
		op := ops[i]
		op(tx)
		var exp Response_Err = Response_OTHER
		assert.Equal(t, 4, len(<-b), Request_Verb_name[i])
		assert.Equal(t, &exp, mustUnmarshal(<-b).ErrCode, Request_Verb_name[i])
	}
}
