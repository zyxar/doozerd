package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/soundcloud/doozerd/consensus"
	"github.com/soundcloud/doozerd/store"
)

type conn struct {
	c        io.ReadWriter
	wl       sync.Mutex // write lock
	addr     string
	p        consensus.Proposer
	st       *store.Store
	canWrite bool
	rwsk     string
	rosk     string
	waccess  bool
	raccess  bool
	self     string
	closed   chan struct{}
}

func (c *conn) serve() {
	defer close(c.closed)

	for {
		var (
			start = time.Now()

			t txn
		)

		t.c = c
		err := c.read(&t.req)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		t.run()

		path := t.req.GetPath()
		if path == "" {
			path = "<empty>"
		}

		statusString := "success"
		if t.resp.ErrCode != nil {
			statusString = Response_Err_name[int32(*t.resp.ErrCode)]
		}

		// Log the transaction
		fmt.Printf(
			"%s %s %s %d %s\n",
			c.addr,
			Request_Verb_name[int32(t.req.GetVerb())],
			path,
			time.Since(start),
			statusString,
		)
	}
}

func (c *conn) read(r *Request) error {
	var size int32
	err := binary.Read(c.c, binary.BigEndian, &size)
	if err != nil {
		return err
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(c.c, buf)
	if err != nil {
		return err
	}

	return proto.Unmarshal(buf, r)
}

func (c *conn) write(r *Response) error {
	buf, err := proto.Marshal(r)
	if err != nil {
		return err
	}

	c.wl.Lock()
	defer c.wl.Unlock()

	err = binary.Write(c.c, binary.BigEndian, int32(len(buf)))
	if err != nil {
		return err
	}

	_, err = c.c.Write(buf)
	return err
}

// Grant compares sk against c.rwsk and c.rosk and
// updates c.waccess and c.raccess as necessary.
// It returns true if sk matched either password.
func (c *conn) grant(sk string) bool {
	switch sk {
	case c.rwsk:
		c.waccess = true
		c.raccess = true
		return true
	case c.rosk:
		c.raccess = true
		return true
	}
	return false
}
