package server

import (
	"log"
	"net"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/soundcloud/doozerd/consensus"
	"github.com/soundcloud/doozerd/store"
)

const PrometheusNamespace = "doozerd"

var (
	conEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Name:      "connection_events_total",
			Help:      "Received connection events total count.",
		},
		[]string{"ev"},
	)

	conCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: PrometheusNamespace,
		Name:      "num_connections",
		Help:      "Current number of open network connections.",
	})

	expvarCollector = prometheus.NewExpvarCollector(map[string]*prometheus.Desc{
		"memstats": prometheus.NewDesc(
			"process_memstats",
			"All numeric memstats as one metric family. This is temporary until these metrics are exported properly by default.",
			[]string{"type"}, nil,
		),
	})
)

func init() {
	prometheus.MustRegister(conEvents)
	prometheus.MustRegister(conCount)
	prometheus.MustRegister(expvarCollector)
}

// ListenAndServe listens on l, accepts network connections, and
// handles requests according to the doozer protocol.
func ListenAndServe(l net.Listener, canWrite chan bool, st *store.Store, p consensus.Proposer, rwsk, rosk string, self string) {
	var w bool
	for {
		c, err := l.Accept()
		if err != nil {
			if err == syscall.EINVAL {
				conEvents.WithLabelValues("invalid").Inc()
				break
			}
			if e, ok := err.(*net.OpError); ok && !e.Temporary() {
				conEvents.WithLabelValues("invalid").Inc()
				break
			}
			log.Println(err)
			conEvents.WithLabelValues("unknown").Inc()

			continue
		}

		// has this server become writable?
		select {
		case w = <-canWrite:
			canWrite = nil
		default:
		}

		conEvents.WithLabelValues("accepted").Inc()
		go serve(c, st, p, w, rwsk, rosk, self)
	}
}

func serve(nc net.Conn, st *store.Store, p consensus.Proposer, w bool, rwsk, rosk string, self string) {
	conCount.Inc()
	defer conCount.Dec()

	c := &conn{
		c:        nc,
		addr:     nc.RemoteAddr().String(),
		st:       st,
		p:        p,
		canWrite: w,
		rwsk:     rwsk,
		rosk:     rosk,
		self:     self,
		closed:   make(chan struct{}),
	}

	c.grant("") // start as if the client supplied a blank password
	c.serve()
	nc.Close()
}
