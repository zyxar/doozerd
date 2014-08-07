package solo

import (
	"net"
	"os"
	"time"

	"github.com/soundcloud/doozer"
	"github.com/soundcloud/doozerd/gc"
	"github.com/soundcloud/doozerd/server"
	"github.com/soundcloud/doozerd/store"
	"github.com/soundcloud/doozerd/web"
)

type proposer struct {
	store *store.Store
	seq   int64
}

func (p *proposer) Propose(v []byte) store.Event {
	var (
		op = store.Op{
			Seqn: 1 + <-p.store.Seqns,
			Mut:  string(v),
		}

		ev store.Event
	)

	p.store.Ops <- op

	for ev.Mut != string(v) {
		p.seq += 1
		w, err := p.store.Wait(store.Any, p.seq)
		if err != nil {
			panic(err) // can't happen
		}
		ev = <-w
	}

	return ev
}

func Main(
	clusterName string,
	name string,
	cl *doozer.Conn,
	listener net.Listener,
	webListener net.Listener,
	st *store.Store,
	history int64,
) {
	var (
		canWrite = make(chan bool, 1)
		p        = &proposer{
			seq:   0,
			store: st,
		}
	)

	listenAddr := listener.Addr().String()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	set(st, "/ctl/name", clusterName, store.Missing)
	set(st, "/ctl/node/"+name+"/addr", listenAddr, store.Missing)
	set(st, "/ctl/node/"+name+"/hostname", hostname, store.Missing)
	set(st, "/ctl/cal/0", name, store.Missing)

	if webListener != nil {
		web.Store = st
		web.ClusterName = clusterName
		go web.Serve(webListener)
	}

	go gc.Clean(st, history, time.Tick(1e9))

	canWrite <- true
	server.ListenAndServe(listener, canWrite, st, p, "", "", name)
}

func set(st *store.Store, path, body string, rev int64) {
	mut := store.MustEncodeSet(path, body, rev)
	st.Ops <- store.Op{1 + <-st.Seqns, mut}
}
