package consensus

var (
	nop      = Msg_NOP.Enum()
	invite   = Msg_INVITE.Enum()
	rsvp     = Msg_RSVP.Enum()
	nominate = Msg_NOMINATE.Enum()
	vote     = Msg_VOTE.Enum()
	tick     = Msg_TICK.Enum()
	propose  = Msg_PROPOSE.Enum()
	learn    = Msg_LEARN.Enum()
)

const nmsg = 8

var (
	msgTick = &Msg{Cmd: tick}
)
