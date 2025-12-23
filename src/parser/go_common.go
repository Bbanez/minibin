package parser

const GoCommon string = `package minibin

import (
	"encoding/binary"
	"math"
)

func PosToBytes(pos int) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(pos))
	return b
}

func BytesToPos(b []byte) int {
	result := binary.LittleEndian.Uint16(b)
	return int(result)
}

func PackString(s string, pos int) []byte {
	result := []byte{1}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	lenB := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenB, uint32(len(s)))
	result = append(result, lenB...)
	result = append(result, []byte(s)...)
	return result
}
func UnpackString(b []byte, atByte int) (string, int, int) {
	// Skip byte 0, it is type
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataLenBytes := b[atByte : atByte+4]
	atByte += 4
	dataLen := int(binary.LittleEndian.Uint32(dataLenBytes))
	dataBytes := b[atByte : atByte+dataLen]
	return string(dataBytes), pos, atByte + dataLen
}

func PackInt32(num int32, pos int) []byte {
	result := []byte{2}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(num))
	result = append(result, b...)
	return result
}
func UnpackInt32(b []byte, atByte int) (int32, int, int) {
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataBytes := b[atByte : atByte+4]
	atByte += 4
	data := binary.LittleEndian.Uint32(dataBytes)
	return int32(data), pos, atByte
}

func PackInt64(num int64, pos int) []byte {
	result := []byte{3}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(num))
	result = append(result, b...)
	return result
}
func UnpackInt64(b []byte, atByte int) (int64, int, int) {
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataBytes := b[atByte : atByte+8]
	atByte += 8
	data := binary.LittleEndian.Uint64(dataBytes)
	return int64(data), pos, atByte
}

func PackUint32(num uint32, pos int) []byte {
	result := []byte{4}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, num)
	result = append(result, b...)
	return result
}
func UnpackUint32(b []byte, atByte int) (uint32, int, int) {
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataBytes := b[atByte : atByte+4]
	atByte += 4
	data := binary.LittleEndian.Uint32(dataBytes)
	return data, pos, atByte
}

func PackUint64(num uint64, pos int) []byte {
	result := []byte{5}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	result = append(result, b...)
	return result
}
func UnpackUint64(b []byte, atByte int) (uint64, int, int) {
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataBytes := b[atByte : atByte+8]
	atByte += 8
	data := binary.LittleEndian.Uint64(dataBytes)
	return data, pos, atByte
}

func PackFloat32(num float32, pos int) []byte {
	result := []byte{6}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(num))
	result = append(result, b...)
	return result
}
func UnpackFloat32(b []byte, atByte int) (float32, int, int) {
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataBytes := b[atByte : atByte+4]
	atByte += 4
	data := math.Float32frombits(
		binary.LittleEndian.Uint32(dataBytes),
	)
	return data, pos, atByte
}

func PackFloat64(num float64, pos int) []byte {
	result := []byte{7}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(num))
	result = append(result, b...)
	return result
}
func UnpackFloat64(b []byte, atByte int) (float64, int, int) {
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataBytes := b[atByte : atByte+8]
	atByte += 8
	data := math.Float64frombits(
		binary.LittleEndian.Uint64(dataBytes),
	)
	return data, pos, atByte
}

func PackBool(num bool, pos int) []byte {
	result := []byte{8}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	if num {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}
	return result
}
func UnpackBool(b []byte, atByte int) (bool, int, int) {
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataBytes := b[atByte]
	atByte += 1
	var data bool
	if dataBytes > 0 {
		data = true
	} else {
		data = false
	}
	return data, pos, atByte
}

func PackObject(b []byte, pos int) []byte {
	result := []byte{9}
	posB := PosToBytes(pos)
	result = append(result, posB...)
	lenB := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenB, uint32(len(b)))
	result = append(result, lenB...)
	result = append(result, b...)
	return result
}
func UnpackObject(b []byte, atByte int) ([]byte, int, int) {
	// Skip byte 0, it is type
	atByte += 1
	pos := BytesToPos(b[atByte : atByte+2])
	atByte += 2
	dataLenBytes := b[atByte : atByte+4]
	atByte += 4
	dataLen := int(binary.LittleEndian.Uint32(dataLenBytes))
	dataBytes := b[atByte : atByte+dataLen]
	return dataBytes, pos, atByte + dataLen
}
	`
