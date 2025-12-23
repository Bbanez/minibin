package parser

const GoCommon string = `package minibin

import (
	"encoding/binary"
	"errors"
	"math"
)

func Compress(data []byte) ([]byte, error) {
	return data, nil
}

func Decompress(compressed []byte) ([]byte, error) {
	return compressed, nil
}

type UnpackabeEntry interface {
	SetPropAtPos(pos int, value any)
}

func Unpack[T UnpackabeEntry](o T, b []byte) error {
	bytes, err := Decompress(b)
	if err != nil {
		return err
	}
	atByte := 0
	for atByte < len(bytes) {
		dataType := bytes[atByte]
		switch dataType {
		case 0:
			atByte++
			continue
		case 1:
			data, pos, next := UnpackString(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 2:
			data, pos, next := UnpackInt32(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 3:
			data, pos, next := UnpackInt64(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 4:
			data, pos, next := UnpackUint32(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 5:
			data, pos, next := UnpackUint64(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 6:
			data, pos, next := UnpackFloat32(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 7:
			data, pos, next := UnpackFloat64(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 8:
			data, pos, next := UnpackBool(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 9:
			data, pos, next := UnpackObject(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 10:
			data, pos, next := UnpackString(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		default:
			return errors.New("Invalid datatype")
		}
	}
	return nil
}

func PackString(s string, pos int) []byte {
	data := []byte(s)
	dataLen := len(data)
	var lenB []byte
	lenD := 1
	if dataLen < 256 {
		lenB = []byte{byte(dataLen)}
	} else if dataLen < 65536 {
		lenD = 2
		lenB := make([]byte, lenD)
		binary.LittleEndian.PutUint16(lenB, uint16(dataLen))
	} else if dataLen < 16777216 {
		lenD = 3
		lenB = []byte{
			byte((0xFF0000 & dataLen) >> 16),
			byte((0x00FF00 & dataLen) >> 8),
			byte(0x0000FF & dataLen),
		}
	} else {
		lenD = 4
		lenB := make([]byte, lenD)
		binary.LittleEndian.PutUint32(lenB, uint32(dataLen))
	}
	result := []byte{1, byte(pos), byte(lenD)}
	result = append(result, lenB...)
	result = append(result, data...)
	return result
}
func UnpackString(b []byte, atByte int) (string, int, int) {
	// Skip byte 0, it is type
	atByte += 1
	pos := int(b[atByte])
	atByte += 1
	lenD := int(b[atByte])
	atByte += 1
	dataLenBytes := b[atByte : atByte+lenD]
	atByte += lenD
	var dataLen int
	if lenD == 4 {
		dataLen = int(binary.LittleEndian.Uint32(dataLenBytes))
	} else if lenD == 3 {
		dataLen = int(dataLenBytes[0])<<16 +
			int(dataLenBytes[1])<<8 +
			int(dataLen)
	} else if lenD == 2 {
		dataLen = int(binary.LittleEndian.Uint16(dataLenBytes))
	} else {
		dataLen = int(dataLenBytes[0])
	}
	dataBytes := b[atByte : atByte+dataLen]
	return string(dataBytes), pos, atByte + dataLen
}

func PackInt32(num int32, pos int) []byte {
	result := []byte{2, byte(pos)}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(num))
	result = append(result, b...)
	return result
}
func UnpackInt32(b []byte, atByte int) (int32, int, int) {
	atByte += 1
	pos := b[atByte]
	atByte += 1
	dataBytes := b[atByte : atByte+4]
	atByte += 4
	data := binary.LittleEndian.Uint32(dataBytes)
	return int32(data), int(pos), atByte
}

func PackInt64(num int64, pos int) []byte {
	result := []byte{3, byte(pos)}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(num))
	result = append(result, b...)
	return result
}
func UnpackInt64(b []byte, atByte int) (int64, int, int) {
	atByte += 1
	pos := b[atByte]
	atByte += 1
	dataBytes := b[atByte : atByte+8]
	atByte += 8
	data := binary.LittleEndian.Uint64(dataBytes)
	return int64(data), int(pos), atByte
}

func PackUint32(num uint32, pos int) []byte {
	result := []byte{4, byte(pos)}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, num)
	result = append(result, b...)
	return result
}
func UnpackUint32(b []byte, atByte int) (uint32, int, int) {
	atByte += 1
	pos := b[atByte]
	atByte += 1
	dataBytes := b[atByte : atByte+4]
	atByte += 4
	data := binary.LittleEndian.Uint32(dataBytes)
	return data, int(pos), atByte
}

func PackUint64(num uint64, pos int) []byte {
	result := []byte{5, byte(pos)}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	result = append(result, b...)
	return result
}
func UnpackUint64(b []byte, atByte int) (uint64, int, int) {
	atByte += 1
	pos := b[atByte]
	atByte += 1
	dataBytes := b[atByte : atByte+8]
	atByte += 8
	data := binary.LittleEndian.Uint64(dataBytes)
	return data, int(pos), atByte
}

func PackFloat32(num float32, pos int) []byte {
	result := []byte{6, byte(pos)}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(num))
	result = append(result, b...)
	return result
}
func UnpackFloat32(b []byte, atByte int) (float32, int, int) {
	atByte += 1
	pos := b[atByte]
	atByte += 1
	dataBytes := b[atByte : atByte+4]
	atByte += 4
	data := math.Float32frombits(
		binary.LittleEndian.Uint32(dataBytes),
	)
	return data, int(pos), atByte
}

func PackFloat64(num float64, pos int) []byte {
	result := []byte{7, byte(pos)}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(num))
	result = append(result, b...)
	return result
}
func UnpackFloat64(b []byte, atByte int) (float64, int, int) {
	atByte += 1
	pos := b[atByte]
	atByte += 1
	dataBytes := b[atByte : atByte+8]
	atByte += 8
	data := math.Float64frombits(
		binary.LittleEndian.Uint64(dataBytes),
	)
	return data, int(pos), atByte
}

func PackBool(num bool, pos int) []byte {
	result := []byte{8, byte(pos)}
	if num {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}
	return result
}
func UnpackBool(b []byte, atByte int) (bool, int, int) {
	atByte += 1
	pos := b[atByte]
	atByte += 1
	dataBytes := b[atByte]
	atByte += 1
	var data bool
	if dataBytes > 0 {
		data = true
	} else {
		data = false
	}
	return data, int(pos), atByte
}

func PackObject(b []byte, pos int) []byte {
	result := []byte{9, byte(pos)}
	lenB := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenB, uint32(len(b)))
	result = append(result, lenB...)
	result = append(result, b...)
	return result
}
func UnpackObject(b []byte, atByte int) ([]byte, int, int) {
	// Skip byte 0, it is type
	atByte += 1
	pos := b[atByte]
	atByte += 1
	dataLenBytes := b[atByte : atByte+4]
	atByte += 4
	dataLen := int(binary.LittleEndian.Uint32(dataLenBytes))
	dataBytes := b[atByte : atByte+dataLen]
	return dataBytes, int(pos), atByte + dataLen
}
`
