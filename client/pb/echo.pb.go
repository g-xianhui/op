// Code generated by protoc-gen-go.
// source: echo.proto
// DO NOT EDIT!

package pb

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type MQEcho struct {
	Data             *string `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MQEcho) Reset()         { *m = MQEcho{} }
func (m *MQEcho) String() string { return proto.CompactTextString(m) }
func (*MQEcho) ProtoMessage()    {}

func (m *MQEcho) GetData() string {
	if m != nil && m.Data != nil {
		return *m.Data
	}
	return ""
}

type MREcho struct {
	Data             *string `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MREcho) Reset()         { *m = MREcho{} }
func (m *MREcho) String() string { return proto.CompactTextString(m) }
func (*MREcho) ProtoMessage()    {}

func (m *MREcho) GetData() string {
	if m != nil && m.Data != nil {
		return *m.Data
	}
	return ""
}

func init() {
}
