package ndn

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	//"fmt"
)

/*
	Base TLV encoding
*/

type TLV struct {
	Type     uint64
	Value    []byte
	Children []TLV
}

func readByte(buf *bytes.Buffer) (result uint64, err error) {
	b, err := buf.ReadByte()
	if err != nil {
		return
	}
	switch b {
	case 0xFF:
		var v64 uint64
		err = binary.Read(buf, binary.BigEndian, &v64)
		if err != nil {
			return
		}
		result = v64
	case 0xFE:
		var v32 uint32
		err = binary.Read(buf, binary.BigEndian, &v32)
		if err != nil {
			return
		}
		result = uint64(v32)
	case 0xFD:
		var v16 uint16
		err = binary.Read(buf, binary.BigEndian, &v16)
		if err != nil {
			return
		}
		result = uint64(v16)
	default:
		result = uint64(b)
	}
	return
}

func writeByte(buf *bytes.Buffer, v uint64) (err error) {
	switch {
	case v > math.MaxUint32:
		buf.WriteByte(0xFF)
		err = binary.Write(buf, binary.BigEndian, v)
	case v > math.MaxUint16:
		buf.WriteByte(0xFE)
		err = binary.Write(buf, binary.BigEndian, uint32(v))
	case v > math.MaxUint8-3:
		buf.WriteByte(0xFD)
		err = binary.Write(buf, binary.BigEndian, uint16(v))
	default:
		err = binary.Write(buf, binary.BigEndian, uint8(v))
	}
	return
}

func (this *TLV) Decode(raw []byte) ([]byte, error) {
	buf := bytes.NewBuffer(raw)
	t, err := readByte(buf)
	if err != nil {
		return nil, err
	}
	this.Type = t
	l, err := readByte(buf)
	if err != nil {
		return nil, err
	}
	this.Value = buf.Next(int(l))
	return buf.Bytes(), nil
}

func (this *TLV) Len() (length uint64) {
	if len(this.Value) == 0 {
		for _, c := range this.Children {
			l := c.Len()
			length += countBytes(c.Type) + countBytes(l) + l
		}
	} else {
		length = uint64(len(this.Value))
	}
	return
}

func (this *TLV) Add(n TLV) {
	this.Children = append(this.Children, n)
}

func countBytes(v uint64) (c uint64) {
	c = 1
	switch {
	case v > math.MaxUint32:
		c += 8
	case v > math.MaxUint16:
		c += 4
	case v > math.MaxUint8-3:
		c += 2
	}
	return
}

func (this *TLV) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := writeByte(buf, this.Type)
	if err != nil {
		return nil, err
	}

	err = writeByte(buf, this.Len())
	if err != nil {
		return nil, err
	}
	if len(this.Value) != 0 && len(this.Children) != 0 {
		return nil, errors.New(VALUE_CHILDREN_COEXIST)
	}
	if len(this.Value) == 0 {
		for _, c := range this.Children {
			b, err := c.Encode()
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
	} else {
		buf.Write(this.Value)
	}
	return buf.Bytes(), nil
}
