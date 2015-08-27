package server

import (
	"io"
	"log"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soundcloud/doozerd/consensus"
	"github.com/soundcloud/doozerd/store"
)

type txn struct {
	c    *conn
	req  Request
	resp Response
}

var ops = map[int32]func(*txn){
	int32(Request_DEL):    (*txn).del,
	int32(Request_GET):    (*txn).get,
	int32(Request_GETDIR): (*txn).getdir,
	int32(Request_NOP):    (*txn).nop,
	int32(Request_REV):    (*txn).rev,
	int32(Request_SET):    (*txn).set,
	int32(Request_STAT):   (*txn).stat,
	int32(Request_SELF):   (*txn).self,
	int32(Request_WAIT):   (*txn).wait,
	int32(Request_WALK):   (*txn).walk,
	int32(Request_ACCESS): (*txn).access,
}

// response flags
const (
	_ = 1 << iota
	_
	set
	del
)

func (t *txn) run() {
	verb := int32(t.req.GetVerb())

	if f, ok := ops[verb]; ok {
		f(t)
	} else {
		t.respondErrCode(Response_UNKNOWN_VERB)
	}
}

func (t *txn) get() {
	trace := t.instrumentVerb()

	if !t.c.raccess {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
		return
	}

	if t.req.Path == nil {
		t.respondErrCode(Response_MISSING_ARG)
		trace("emissingarg")
		return
	}

	go func() {
		g, err := t.getter()
		if err != nil {
			t.respondOsError(err)
			trace("eunknown")
			return
		}

		v, rev := g.Get(*t.req.Path)
		if rev == store.Dir {
			t.respondErrCode(Response_ISDIR)
			trace("eisdir")
			return
		}

		t.resp.Rev = &rev
		if len(v) == 1 { // not missing
			t.resp.Value = []byte(v[0])
		}
		t.respond()
		trace("success")
	}()
}

func (t *txn) set() {
	trace := t.instrumentVerb()

	if !t.c.waccess {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
		return
	}

	if !t.c.canWrite {
		t.respondErrCode(Response_READONLY)
		trace("ereadonly")
		return
	}

	if t.req.Path == nil || t.req.Rev == nil {
		t.respondErrCode(Response_MISSING_ARG)
		trace("emissingarg")
		return
	}

	go func() {
		ev := consensus.Set(t.c.p, *t.req.Path, t.req.Value, *t.req.Rev)
		if ev.Err != nil {
			t.respondOsError(ev.Err)
			trace("eunknown")
			return
		}
		t.resp.Rev = &ev.Seqn
		t.respond()
		trace("success")
	}()
}

func (t *txn) del() {
	trace := t.instrumentVerb()

	if !t.c.waccess {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
		return
	}

	if !t.c.canWrite {
		t.respondErrCode(Response_READONLY)
		trace("ereadonly")
		return
	}

	if t.req.Path == nil || t.req.Rev == nil {
		t.respondErrCode(Response_MISSING_ARG)
		trace("emissingarg")
		return
	}

	go func() {
		ev := consensus.Del(t.c.p, *t.req.Path, *t.req.Rev)
		if ev.Err != nil {
			t.respondOsError(ev.Err)
			trace("eunknown")
			return
		}
		t.respond()
		trace("success")
	}()
}

func (t *txn) nop() {
	trace := t.instrumentVerb()

	if !t.c.waccess {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
		return
	}

	if !t.c.canWrite {
		t.respondErrCode(Response_READONLY)
		trace("ereadonly")
		return
	}

	go func() {
		t.c.p.Propose([]byte(store.Nop))
		t.respond()
		trace("success")
	}()
}

func (t *txn) rev() {
	trace := t.instrumentVerb()

	rev := <-t.c.st.Seqns
	t.resp.Rev = &rev
	t.respond()
	trace("success")
}

func (t *txn) self() {
	trace := t.instrumentVerb()
	t.resp.Value = []byte(t.c.self)
	t.respond()
	trace("success")
}

func (t *txn) stat() {
	trace := t.instrumentVerb()

	if !t.c.raccess {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
		return
	}

	go func() {
		g, err := t.getter()
		if err != nil {
			t.respondOsError(err)
			trace("eunknown")
			return
		}

		len, rev := g.Stat(t.req.GetPath())
		t.resp.Len = &len
		t.resp.Rev = &rev
		t.respond()
		trace("success")
	}()
}

func (t *txn) getdir() {
	trace := t.instrumentVerb()

	if !t.c.raccess {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
		return
	}

	if t.req.Path == nil || t.req.Offset == nil {
		t.respondErrCode(Response_MISSING_ARG)
		trace("emissingarg")
		return
	}

	go func() {
		g, err := t.getter()
		if err != nil {
			t.respondOsError(err)
			trace("eunknown")
			return
		}

		ents, rev := g.Get(*t.req.Path)
		if rev == store.Missing {
			t.respondErrCode(Response_NOENT)
			trace("enoent")
			return
		}
		if rev != store.Dir {
			t.respondErrCode(Response_NOTDIR)
			trace("enotdir")
			return
		}

		sort.Strings(ents)
		offset := int(*t.req.Offset)
		if offset < 0 || offset >= len(ents) {
			t.respondErrCode(Response_RANGE)
			trace("erange")
			return
		}

		t.resp.Path = &ents[offset]
		t.respond()
		trace("success")
	}()
}

func (t *txn) wait() {
	trace := t.instrumentVerb()

	if !t.c.raccess {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
		return
	}

	if t.req.Path == nil || t.req.Rev == nil {
		t.respondErrCode(Response_MISSING_ARG)
		trace("emissingarg")
		return
	}

	glob, err := store.CompileGlob(*t.req.Path)
	if err != nil {
		t.respondOsError(err)
		trace("eunknown")
		return
	}

	ch, err := t.c.st.Wait(glob, *t.req.Rev)
	if err != nil {
		t.respondOsError(err)
		trace("eunknown")
		return
	}

	go func() {
		var ev store.Event

		select {
		case ev = <-ch:
		case <-t.c.closed:
			t.c.st.Cancel(ch)
			return
		}

		t.resp.Path = &ev.Path
		t.resp.Value = []byte(ev.Body)
		t.resp.Rev = &ev.Seqn
		switch {
		case ev.IsSet():
			t.resp.Flags = proto.Int32(set)
		case ev.IsDel():
			t.resp.Flags = proto.Int32(del)
		default:
			t.resp.Flags = proto.Int32(0)
		}
		t.respond()
		trace("success")
	}()
}

func (t *txn) walk() {
	trace := t.instrumentVerb()

	if !t.c.raccess {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
		return
	}

	if t.req.Path == nil || t.req.Offset == nil {
		t.respondErrCode(Response_MISSING_ARG)
		trace("emissingarg")
		return
	}

	glob, err := store.CompileGlob(*t.req.Path)
	if err != nil {
		t.respondOsError(err)
		trace("eunknown")
		return
	}

	offset := *t.req.Offset
	if offset < 0 {
		t.respondErrCode(Response_RANGE)
		trace("erange")
		return
	}

	go func() {
		g, err := t.getter()
		if err != nil {
			t.respondOsError(err)
			trace("eunknown")
			return
		}

		f := func(path, body string, rev int64) (stop bool) {
			if offset == 0 {
				t.resp.Path = &path
				t.resp.Value = []byte(body)
				t.resp.Rev = &rev
				t.resp.Flags = proto.Int32(set)
				t.respond()
				trace("success")
				return true
			}
			offset--
			return false
		}
		if !store.Walk(g, glob, f) {
			t.respondErrCode(Response_RANGE)
			trace("erange")
		}
	}()
}

func (t *txn) access() {
	trace := t.instrumentVerb()

	if t.c.grant(string(t.req.Value)) {
		t.respond()
		trace("success")
	} else {
		t.respondOsError(syscall.EACCES)
		trace("eacces")
	}
}

func (t *txn) respondOsError(err error) {
	switch err {
	case store.ErrBadPath:
		t.respondErrCode(Response_BAD_PATH)
	case store.ErrRevMismatch:
		t.respondErrCode(Response_REV_MISMATCH)
	case store.ErrTooLate:
		t.respondErrCode(Response_TOO_LATE)
	case syscall.EISDIR:
		t.respondErrCode(Response_ISDIR)
	case syscall.ENOTDIR:
		t.respondErrCode(Response_NOTDIR)
	default:
		t.resp.ErrDetail = proto.String(err.Error())
		t.respondErrCode(Response_OTHER)
	}
}

func (t *txn) respondErrCode(e Response_Err) {
	t.resp.ErrCode = &e
	t.respond()
}

func (t *txn) respond() {
	t.resp.Tag = t.req.Tag
	err := t.c.write(&t.resp)
	if err != nil && err != io.EOF {
		log.Println(err)
	}
}

func (t *txn) getter() (store.Getter, error) {
	if t.req.Rev == nil {
		_, g := t.c.st.Snap()
		return g, nil
	}

	ch, err := t.c.st.Wait(store.Any, *t.req.Rev)
	if err != nil {
		return nil, err
	}
	return <-ch, nil
}

func (t *txn) instrumentVerb() func(string) {
	start := time.Now()

	return func(result string) {
		verb := t.req.GetVerb().String()
		if verb == "" {
			verb = "unknown"
		} else {
			verb = strings.ToLower(verb)
		}

		host := t.c.addr
		parts := strings.Split(host, ":")
		if len(parts) != 2 {
			host = "unknown"
		} else {
			host = parts[0]
		}

		dur := float64(time.Since(start) / time.Microsecond)

		txnLatencies.WithLabelValues(verb, host, result).Observe(dur)
	}
}

var txnLatencies = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace: PrometheusNamespace,
		Name:      "txn_latency_microseconds",
		Help:      "Transaction latencies in microseconds partitioned by operation, host, and outcome.",
	},
	[]string{"verb", "host", "result"},
)

func init() {
	prometheus.MustRegister(txnLatencies)
}
