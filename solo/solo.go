package solo

import (
	"net"
	"os"
	"time"

	"github.com/soundcloud/doozer"
	"github.com/zyxar/doozerd/gc"
	"github.com/zyxar/doozerd/server"
	"github.com/zyxar/doozerd/store"
	"github.com/zyxar/doozerd/web"
)

type proposal struct {
	eventc chan store.Event
	mut    string
}

type proposer struct {
	proposalc chan proposal
	store     *store.Store
}

func (p *proposer) Propose(v []byte) store.Event {
	eventc := make(chan store.Event)

	p.proposalc <- proposal{
		eventc: eventc,
		mut:    string(v),
	}

	return <-eventc
}

func (p *proposer) process() {
	for prop := range p.proposalc {
		op := store.Op{
			Seqn: 1 + <-p.store.Seqns,
			Mut:  prop.mut,
		}

		p.store.Ops <- op

		waitc, err := p.store.Wait(store.Any, op.Seqn)
		if err != nil {
			panic(err) // can't happen, but happened before.
		}
		ev := <-waitc

		// This is a safety measure if the sequential nature of solo mode missed a
		// corner-case.
		if ev.Mut == prop.mut {
			prop.eventc <- ev
			continue
		}

		panic("not reachable")
	}
}

// Main takes care of essentail setup of proposer, gc and server.
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
			proposalc: make(chan proposal),
			store:     st,
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

	// sequential handling of mutations
	go p.process()

	go gc.Clean(st, history, time.Tick(1e9))

	canWrite <- true
	server.ListenAndServe(listener, canWrite, st, p, "", "", name)
}

func set(st *store.Store, path, body string, rev int64) {
	mut := store.MustEncodeSet(path, body, rev)
	st.Ops <- store.Op{
		Mut:  mut,
		Seqn: 1 + <-st.Seqns,
	}
}
