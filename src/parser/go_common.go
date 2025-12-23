package parser

const GoCommon string = `package minibin

import (
	"encoding/binary"
	"math"
)

func PackString(s string) []byte {
	result := []byte{
		1,
	}
	lenB := make([]byte, 4)
	binary.BigEndian.PutUint32(lenB, uint32(len(s)))
	result = append(result, lenB...)
	result = append(result, []byte(s)...)
	return result
}

func PackInt32(num int32) []byte {
	result := []byte{
		2,
	}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(num))
	result = append(result, b...)
	return result
}

func PackInt64(num int64) []byte {
	result := []byte{
		3,
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(num))
	result = append(result, b...)
	return result
}

func PackUint32(num uint32) []byte {
	result := []byte{
		4,
	}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, num)
	result = append(result, b...)
	return result
}

func PackUint64(num uint64) []byte {
	result := []byte{
		5,
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	result = append(result, b...)
	return result
}

func PackFloat32(num float32) []byte {
	result := []byte{
		6,
	}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(num))
	result = append(result, b...)
	return result
}

func PackFloat64(num float64) []byte {
	result := []byte{
		7,
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(num))
	result = append(result, b...)
	return result
}

func PackBool(num bool) []byte {
	if num {
		return []byte{8, 1}
	}
	return []byte{8, 0}
}

func PackObject(b []byte) []byte {
	result := []byte{
		9,
	}
	lenB := make([]byte, 4)
	binary.BigEndian.PutUint32(lenB, uint32(len(b)))
	result = append(result, lenB...)
	result = append(result, b...)
	return result
}
	`
