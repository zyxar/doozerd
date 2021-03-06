// Code generated by protoc-gen-go.
// source: m.proto
// DO NOT EDIT!

/*
Package consensus is a generated protocol buffer package.

It is generated from these files:
	m.proto

It has these top-level messages:
	Msg
*/
package consensus

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Msg_Cmd int32

const (
	Msg_NOP      Msg_Cmd = 0
	Msg_INVITE   Msg_Cmd = 1
	Msg_RSVP     Msg_Cmd = 2
	Msg_NOMINATE Msg_Cmd = 3
	Msg_VOTE     Msg_Cmd = 4
	Msg_TICK     Msg_Cmd = 5
	Msg_PROPOSE  Msg_Cmd = 6
	Msg_LEARN    Msg_Cmd = 7
)

var Msg_Cmd_name = map[int32]string{
	0: "NOP",
	1: "INVITE",
	2: "RSVP",
	3: "NOMINATE",
	4: "VOTE",
	5: "TICK",
	6: "PROPOSE",
	7: "LEARN",
}
var Msg_Cmd_value = map[string]int32{
	"NOP":      0,
	"INVITE":   1,
	"RSVP":     2,
	"NOMINATE": 3,
	"VOTE":     4,
	"TICK":     5,
	"PROPOSE":  6,
	"LEARN":    7,
}

func (x Msg_Cmd) Enum() *Msg_Cmd {
	p := new(Msg_Cmd)
	*p = x
	return p
}
func (x Msg_Cmd) String() string {
	return proto.EnumName(Msg_Cmd_name, int32(x))
}
func (x *Msg_Cmd) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Msg_Cmd_value, data, "Msg_Cmd")
	if err != nil {
		return err
	}
	*x = Msg_Cmd(value)
	return nil
}

type Msg struct {
	Cmd              *Msg_Cmd `protobuf:"varint,1,opt,name=cmd,enum=consensus.Msg_Cmd" json:"cmd,omitempty"`
	Seqn             *int64   `protobuf:"varint,2,opt,name=seqn" json:"seqn,omitempty"`
	Crnd             *int64   `protobuf:"varint,3,opt,name=crnd" json:"crnd,omitempty"`
	Vrnd             *int64   `protobuf:"varint,4,opt,name=vrnd" json:"vrnd,omitempty"`
	Value            []byte   `protobuf:"bytes,5,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Msg) Reset()         { *m = Msg{} }
func (m *Msg) String() string { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()    {}

func (m *Msg) GetCmd() Msg_Cmd {
	if m != nil && m.Cmd != nil {
		return *m.Cmd
	}
	return Msg_NOP
}

func (m *Msg) GetSeqn() int64 {
	if m != nil && m.Seqn != nil {
		return *m.Seqn
	}
	return 0
}

func (m *Msg) GetCrnd() int64 {
	if m != nil && m.Crnd != nil {
		return *m.Crnd
	}
	return 0
}

func (m *Msg) GetVrnd() int64 {
	if m != nil && m.Vrnd != nil {
		return *m.Vrnd
	}
	return 0
}

func (m *Msg) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterEnum("consensus.Msg_Cmd", Msg_Cmd_name, Msg_Cmd_value)
}
