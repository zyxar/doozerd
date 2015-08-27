package solo

import (
	"net"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/soundcloud/doozer"
	"github.com/zyxar/doozerd/store"
)

func TestDoozerNop(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())
	err := cl.Nop()
	assert.Equal(t, nil, err)
}

func TestDoozerGet(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())

	_, err := cl.Set("/x", store.Missing, []byte{'a'})
	assert.Equal(t, nil, err)

	ents, rev, err := cl.Get("/x", nil)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, store.Dir, rev)
	assert.Equal(t, []byte{'a'}, ents)
}

func TestDoozerSet(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())

	for i := byte(0); i < 10; i++ {
		_, err := cl.Set("/x", store.Clobber, []byte{'0' + i})
		assert.Equal(t, nil, err)
	}

	_, err := cl.Set("/x", 0, []byte{'X'})
	assert.Equal(t, &doozer.Error{doozer.ErrOldRev, ""}, err)
}

func TestDoozerGetWithRev(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())

	rev1, err := cl.Set("/x", store.Missing, []byte{'a'})
	assert.Equal(t, nil, err)

	v, rev, err := cl.Get("/x", &rev1) // Use the snapshot.
	assert.Equal(t, nil, err)
	assert.Equal(t, rev1, rev)
	assert.Equal(t, []byte{'a'}, v)

	rev2, err := cl.Set("/x", rev, []byte{'b'})
	assert.Equal(t, nil, err)

	v, rev, err = cl.Get("/x", nil) // Read the new value.
	assert.Equal(t, nil, err)
	assert.Equal(t, rev2, rev)
	assert.Equal(t, []byte{'b'}, v)

	v, rev, err = cl.Get("/x", &rev1) // Read the saved value again.
	assert.Equal(t, nil, err)
	assert.Equal(t, rev1, rev)
	assert.Equal(t, []byte{'a'}, v)
}

func TestDoozerWaitSimple(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())
	var rev int64 = 1

	cl.Set("/test/foo", store.Clobber, []byte("bar"))
	ev, err := cl.Wait("/test/**", rev)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/test/foo", ev.Path)
	assert.Equal(t, []byte("bar"), ev.Body)
	assert.T(t, ev.IsSet())
	rev = ev.Rev + 1

	cl.Set("/test/fun", store.Clobber, []byte("house"))
	ev, err = cl.Wait("/test/**", rev)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/test/fun", ev.Path)
	assert.Equal(t, []byte("house"), ev.Body)
	assert.T(t, ev.IsSet())
	rev = ev.Rev + 1

	cl.Del("/test/foo", store.Clobber)
	ev, err = cl.Wait("/test/**", rev)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/test/foo", ev.Path)
	assert.T(t, ev.IsDel())
}

func TestDoozerWaitWithRev(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())

	// Create some history
	cl.Set("/test/foo", store.Clobber, []byte("bar"))
	cl.Set("/test/fun", store.Clobber, []byte("house"))

	ev, err := cl.Wait("/test/**", 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/test/foo", ev.Path)
	assert.Equal(t, []byte("bar"), ev.Body)
	assert.T(t, ev.IsSet())
	rev := ev.Rev + 1

	ev, err = cl.Wait("/test/**", rev)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/test/fun", ev.Path)
	assert.Equal(t, []byte("house"), ev.Body)
	assert.T(t, ev.IsSet())
}

func TestDoozerStat(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())

	cl.Set("/test/foo", store.Clobber, []byte("bar"))
	setRev, _ := cl.Set("/test/fun", store.Clobber, []byte("house"))

	ln, rev, err := cl.Stat("/test", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, store.Dir, rev)
	assert.Equal(t, int(2), ln)

	ln, rev, err = cl.Stat("/test/fun", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, setRev, rev)
	assert.Equal(t, int(5), ln)
}

func TestDoozerGetdirOnDir(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())

	cl.Set("/test/a", store.Clobber, []byte("1"))
	cl.Set("/test/b", store.Clobber, []byte("2"))
	cl.Set("/test/c", store.Clobber, []byte("3"))

	rev, err := cl.Rev()
	if err != nil {
		panic(err)
	}

	got, err := cl.Getdir("/test", rev, 0, -1)
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"a", "b", "c"}, got)
}

func TestDoozerGetdirOnFile(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())

	cl.Set("/test/a", store.Clobber, []byte("1"))

	rev, err := cl.Rev()
	if err != nil {
		panic(err)
	}

	names, err := cl.Getdir("/test/a", rev, 0, -1)
	assert.Equal(t, &doozer.Error{doozer.ErrNotDir, ""}, err)
	assert.Equal(t, []string(nil), names)
}

func TestDoozerGetdirMissing(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())

	rev, err := cl.Rev()
	if err != nil {
		panic(err)
	}

	names, err := cl.Getdir("/not/here", rev, 0, -1)
	assert.Equal(t, &doozer.Error{doozer.ErrNoEnt, ""}, err)
	assert.Equal(t, []string(nil), names)
}

func TestDoozerGetdirOffsetLimit(t *testing.T) {
	var (
		l  = mustListen()
		st = store.New(store.DefaultInitialRev)
	)
	defer l.Close()

	go Main("a", "X", nil, l, nil, st, 1e6)

	cl := dial(l.Addr().String())
	cl.Set("/test/a", store.Clobber, []byte("1"))
	cl.Set("/test/b", store.Clobber, []byte("2"))
	cl.Set("/test/c", store.Clobber, []byte("3"))
	cl.Set("/test/d", store.Clobber, []byte("4"))

	rev, err := cl.Rev()
	if err != nil {
		panic(err)
	}

	names, err := cl.Getdir("/test", rev, 1, 2)
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"b", "c"}, names)
}

func dial(addr string) *doozer.Conn {
	c, err := doozer.Dial(addr)
	if err != nil {
		panic(err)
	}
	return c
}

func mustListen() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	return l
}
