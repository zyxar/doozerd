// Code generated by protoc-gen-go.
// source: msg.proto
// DO NOT EDIT!

/*
Package server is a generated protocol buffer package.

It is generated from these files:
	msg.proto

It has these top-level messages:
	Request
	Response
*/
package server

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Request_Verb int32

const (
	Request_GET    Request_Verb = 1
	Request_SET    Request_Verb = 2
	Request_DEL    Request_Verb = 3
	Request_REV    Request_Verb = 5
	Request_WAIT   Request_Verb = 6
	Request_NOP    Request_Verb = 7
	Request_WALK   Request_Verb = 9
	Request_GETDIR Request_Verb = 14
	Request_STAT   Request_Verb = 16
	Request_SELF   Request_Verb = 20
	Request_ACCESS Request_Verb = 99
)

var Request_Verb_name = map[int32]string{
	1:  "GET",
	2:  "SET",
	3:  "DEL",
	5:  "REV",
	6:  "WAIT",
	7:  "NOP",
	9:  "WALK",
	14: "GETDIR",
	16: "STAT",
	20: "SELF",
	99: "ACCESS",
}
var Request_Verb_value = map[string]int32{
	"GET":    1,
	"SET":    2,
	"DEL":    3,
	"REV":    5,
	"WAIT":   6,
	"NOP":    7,
	"WALK":   9,
	"GETDIR": 14,
	"STAT":   16,
	"SELF":   20,
	"ACCESS": 99,
}

func (x Request_Verb) Enum() *Request_Verb {
	p := new(Request_Verb)
	*p = x
	return p
}
func (x Request_Verb) String() string {
	return proto.EnumName(Request_Verb_name, int32(x))
}
func (x *Request_Verb) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Request_Verb_value, data, "Request_Verb")
	if err != nil {
		return err
	}
	*x = Request_Verb(value)
	return nil
}

type Response_Err int32

const (
	// don't use value 0
	Response_OTHER        Response_Err = 127
	Response_TAG_IN_USE   Response_Err = 1
	Response_UNKNOWN_VERB Response_Err = 2
	Response_READONLY     Response_Err = 3
	Response_TOO_LATE     Response_Err = 4
	Response_REV_MISMATCH Response_Err = 5
	Response_BAD_PATH     Response_Err = 6
	Response_MISSING_ARG  Response_Err = 7
	Response_RANGE        Response_Err = 8
	Response_NOTDIR       Response_Err = 20
	Response_ISDIR        Response_Err = 21
	Response_NOENT        Response_Err = 22
)

var Response_Err_name = map[int32]string{
	127: "OTHER",
	1:   "TAG_IN_USE",
	2:   "UNKNOWN_VERB",
	3:   "READONLY",
	4:   "TOO_LATE",
	5:   "REV_MISMATCH",
	6:   "BAD_PATH",
	7:   "MISSING_ARG",
	8:   "RANGE",
	20:  "NOTDIR",
	21:  "ISDIR",
	22:  "NOENT",
}
var Response_Err_value = map[string]int32{
	"OTHER":        127,
	"TAG_IN_USE":   1,
	"UNKNOWN_VERB": 2,
	"READONLY":     3,
	"TOO_LATE":     4,
	"REV_MISMATCH": 5,
	"BAD_PATH":     6,
	"MISSING_ARG":  7,
	"RANGE":        8,
	"NOTDIR":       20,
	"ISDIR":        21,
	"NOENT":        22,
}

func (x Response_Err) Enum() *Response_Err {
	p := new(Response_Err)
	*p = x
	return p
}
func (x Response_Err) String() string {
	return proto.EnumName(Response_Err_name, int32(x))
}
func (x *Response_Err) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Response_Err_value, data, "Response_Err")
	if err != nil {
		return err
	}
	*x = Response_Err(value)
	return nil
}

// see doc/proto.md
type Request struct {
	Tag              *int32        `protobuf:"varint,1,opt,name=tag" json:"tag,omitempty"`
	Verb             *Request_Verb `protobuf:"varint,2,opt,name=verb,enum=server.Request_Verb" json:"verb,omitempty"`
	Path             *string       `protobuf:"bytes,4,opt,name=path" json:"path,omitempty"`
	Value            []byte        `protobuf:"bytes,5,opt,name=value" json:"value,omitempty"`
	OtherTag         *int32        `protobuf:"varint,6,opt,name=other_tag" json:"other_tag,omitempty"`
	Offset           *int32        `protobuf:"varint,7,opt,name=offset" json:"offset,omitempty"`
	Rev              *int64        `protobuf:"varint,9,opt,name=rev" json:"rev,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}

func (m *Request) GetTag() int32 {
	if m != nil && m.Tag != nil {
		return *m.Tag
	}
	return 0
}

func (m *Request) GetVerb() Request_Verb {
	if m != nil && m.Verb != nil {
		return *m.Verb
	}
	return Request_GET
}

func (m *Request) GetPath() string {
	if m != nil && m.Path != nil {
		return *m.Path
	}
	return ""
}

func (m *Request) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Request) GetOtherTag() int32 {
	if m != nil && m.OtherTag != nil {
		return *m.OtherTag
	}
	return 0
}

func (m *Request) GetOffset() int32 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return 0
}

func (m *Request) GetRev() int64 {
	if m != nil && m.Rev != nil {
		return *m.Rev
	}
	return 0
}

// see doc/proto.md
type Response struct {
	Tag              *int32        `protobuf:"varint,1,opt,name=tag" json:"tag,omitempty"`
	Flags            *int32        `protobuf:"varint,2,opt,name=flags" json:"flags,omitempty"`
	Rev              *int64        `protobuf:"varint,3,opt,name=rev" json:"rev,omitempty"`
	Path             *string       `protobuf:"bytes,5,opt,name=path" json:"path,omitempty"`
	Value            []byte        `protobuf:"bytes,6,opt,name=value" json:"value,omitempty"`
	Len              *int32        `protobuf:"varint,8,opt,name=len" json:"len,omitempty"`
	ErrCode          *Response_Err `protobuf:"varint,100,opt,name=err_code,enum=server.Response_Err" json:"err_code,omitempty"`
	ErrDetail        *string       `protobuf:"bytes,101,opt,name=err_detail" json:"err_detail,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}

func (m *Response) GetTag() int32 {
	if m != nil && m.Tag != nil {
		return *m.Tag
	}
	return 0
}

func (m *Response) GetFlags() int32 {
	if m != nil && m.Flags != nil {
		return *m.Flags
	}
	return 0
}

func (m *Response) GetRev() int64 {
	if m != nil && m.Rev != nil {
		return *m.Rev
	}
	return 0
}

func (m *Response) GetPath() string {
	if m != nil && m.Path != nil {
		return *m.Path
	}
	return ""
}

func (m *Response) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Response) GetLen() int32 {
	if m != nil && m.Len != nil {
		return *m.Len
	}
	return 0
}

func (m *Response) GetErrCode() Response_Err {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return Response_OTHER
}

func (m *Response) GetErrDetail() string {
	if m != nil && m.ErrDetail != nil {
		return *m.ErrDetail
	}
	return ""
}

func init() {
	proto.RegisterEnum("server.Request_Verb", Request_Verb_name, Request_Verb_value)
	proto.RegisterEnum("server.Response_Err", Response_Err_name, Response_Err_value)
}
