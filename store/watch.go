package store

type watch struct {
	glob *Glob
	rev  int64
	c    chan Event
}

func (w *watch) notify(e Event) bool {
	if e.Seqn >= w.rev && w.glob.Match(e.Path) {
		w.c <- e
		return true
	}

	return false
}

type watches []*watch

func (ws *watches) notify(e Event) {
	var i int
	for i < len(*ws) {
		if (*ws)[i].notify(e) {
			ws.delete(i)
		} else {
			i++
		}
	}
}

func (ws *watches) delete(i int) {
	old := *ws
	copy(old[i:], old[i+1:])
	old[len(old)-1] = nil
	*ws = old[:len(old)-1]
}

func (ws watches) close() {
	for _, w := range ws {
		close(w.c)
	}
}

func (ws *watches) cancel(ch <-chan Event) {
	for i, watch := range *ws {
		if watch.c == ch {
			ws.delete(i)
			return
		}
	}
}
