package consensus

type coordinator struct {
	size int
	quor int

	begun  bool
	target string
	crnd   int64
	cval   string
	rsvp   []bool
	nrsvp  int
	vr     int64
	vv     string

	sched bool
}

func (co *coordinator) update(p *packet, from int) (m *Msg, wantTick bool) {
	in := &p.Msg
	switch *in.Cmd {
	case Msg_PROPOSE:
		if co.begun {
			break
		}

		co.begun = true
		co.target = string(in.Value)
		co.vr = 0
		co.vv = ""
		co.rsvp = make([]bool, co.size)
		co.cval = ""
		return &Msg{Cmd: invite, Crnd: &co.crnd}, true
	case Msg_RSVP:
		if !co.begun {
			break
		}

		if in.Crnd == nil || in.Vrnd == nil {
			break
		}

		i, vrnd, vval := *in.Crnd, *in.Vrnd, in.Value

		if co.cval != "" {
			break
		}

		if i != co.crnd {
			break
		}

		if vrnd > co.vr {
			co.vr = vrnd
			co.vv = string(vval)
		}

		if !co.rsvp[from] {
			co.rsvp[from] = true
			co.nrsvp++
		}
		if co.nrsvp >= co.quor {
			var v string

			if co.vr > 0 {
				v = co.vv
			} else {
				v = co.target
			}
			co.cval = v

			return &Msg{Cmd: nominate, Crnd: &co.crnd, Value: []byte(v)}, false
		}
	case Msg_TICK:
		co.crnd += int64(co.size)
		co.vr = 0
		co.vv = ""
		co.rsvp = make([]bool, co.size)
		co.nrsvp = 0
		co.cval = ""
		co.sched = false
		return &Msg{Cmd: invite, Crnd: &co.crnd}, true
	}

	return
}
