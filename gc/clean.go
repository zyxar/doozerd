package gc

import (
	"github.com/soundcloud/doozerd/store"
	"time"
)

func Clean(st *store.Store, keep int64, ticker <-chan time.Time) {
	for _ = range ticker {
		last := (<-st.Seqns) - keep
		st.Clean(last)
	}
}
