package consensus

import (
	"testing"

	"github.com/bmizerany/assert"
	_ "github.com/zyxar/doozerd/quiet"
)

func TestIgnoreOldMessages(t *testing.T) {
	tests := [][]*Msg{
		{newInviteSeqn1(11), newNominateSeqn1(1, "v")},
		{newNominateSeqn1(11, "v"), newInviteSeqn1(1)},
		{newInviteSeqn1(11), newInviteSeqn1(1)},
		{newNominateSeqn1(11, "v"), newNominateSeqn1(1, "v")},
	}

	for _, test := range tests {
		ac := acceptor{}

		ac.update(test[0])

		got := ac.update(test[1])
		assert.Equal(t, (*Msg)(nil), got)
	}
}

func TestAcceptsInvite(t *testing.T) {
	ac := acceptor{}
	got := ac.update(newInviteSeqn1(1))
	assert.Equal(t, newRsvp(1, 0, ""), got)
}

func TestItVotes(t *testing.T) {
	totest := [][]*Msg{
		{newNominateSeqn1(1, "foo"), newVote(1, "foo")},
		{newNominateSeqn1(1, "bar"), newVote(1, "bar")},
	}

	for _, test := range totest {
		ac := acceptor{}
		got := ac.update(test[0])
		assert.Equal(t, test[1], got, test)
	}
}

func TestItVotesWithAnotherRound(t *testing.T) {
	ac := acceptor{}
	val := "bar"

	// According to paxos, we can omit Phase 1 in the first round
	got := ac.update(newNominateSeqn1(2, val))
	assert.Equal(t, newVote(2, val), got)
}

func TestItVotesWithAnotherSelf(t *testing.T) {
	ac := acceptor{}
	val := "bar"

	// According to paxos, we can omit Phase 1 in the first round
	got := ac.update(newNominateSeqn1(2, val))
	assert.Equal(t, newVote(2, val), got)
}

func TestVotedRoundsAndValuesAreTracked(t *testing.T) {
	ac := acceptor{}

	ac.update(newNominateSeqn1(1, "v"))

	got := ac.update(newInviteSeqn1(2))
	assert.Equal(t, newRsvp(2, 1, "v"), got)
}

func TestVotesOnlyOncePerRound(t *testing.T) {
	ac := acceptor{}

	got := ac.update(newNominateSeqn1(1, "v"))
	assert.Equal(t, newVote(1, "v"), got)

	got = ac.update(newNominateSeqn1(1, "v"))
	assert.Equal(t, (*Msg)(nil), got)
}

func TestAcceptorIgnoresBadMessages(t *testing.T) {
	ac := acceptor{}

	got := ac.update(&Msg{Cmd: invite}) // missing Crnd
	assert.Equal(t, (*Msg)(nil), got)

	got = ac.update(&Msg{Cmd: nominate}) // missing Crnd
	assert.Equal(t, (*Msg)(nil), got)
}
