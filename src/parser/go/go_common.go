package parser_go

const Common string = `package minibin

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
		pos := int(bytes[atByte])
		atByte += 1
		typ, lenD := unmergeDataTypeAndLenDataLen(bytes[atByte])
		atByte += 1
		switch typ {
		case 0:
			atByte++
			continue
		case 1:
			data, next := UnpackString(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 2:
			data, next := UnpackInt32(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 3:
			data, next := UnpackInt64(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 4:
			data, next := UnpackUint32(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 5:
			data, next := UnpackUint64(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 6:
			data, next := UnpackFloat32(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 7:
			data, next := UnpackFloat64(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 8:
			data, next := UnpackBool(bytes, atByte)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 9:
			data, next := UnpackObject(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 10:
			data, next := UnpackString(bytes, atByte, lenD)
			atByte = next
			o.SetPropAtPos(pos, data)
		case 11:
			data, next := UnpackBytes(bytes, atByte, lenD)
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
	lenD, dataLenBytes := splitUint32(uint32(len(data)))
	typLenD := mergeDataTypeAndLenDataLen(1, byte(lenD))
	result := []byte{byte(pos), typLenD}
	result = append(result, dataLenBytes...)
	result = append(result, data...)
	return result
}
func UnpackString(b []byte, atByte int, lenD int) (string, int) {
	lenD++
	dataLen := int(mergeUint32(lenD, b[atByte:atByte+lenD]))
	atByte += lenD
	dataBytes := b[atByte : atByte+dataLen]
	return string(dataBytes), atByte + dataLen
}

func PackInt32(num int32, pos int) []byte {
	lenD, data := splitUint32(uint32(num))
	typLenD := mergeDataTypeAndLenDataLen(2, byte(lenD))
	result := []byte{byte(pos), typLenD}
	result = append(result, data...)
	return result
}
func UnpackInt32(b []byte, atByte int, lenD int) (int32, int) {
	// This is required because lenD == 0 represents 1 byte of data
	lenD++
	data := mergeUint32(lenD, b[atByte:atByte+lenD])
	atByte += lenD
	return int32(data), atByte
}

func PackInt64(num int64, pos int) []byte {
	lenD, data := splitUint64(uint64(num))
	typLenD := mergeDataTypeAndLenDataLen(3, byte(lenD))
	result := []byte{byte(pos), typLenD}
	result = append(result, data...)
	return result
}
func UnpackInt64(b []byte, atByte int, lenD int) (int64, int) {
	// This is required because lenD == 0 represents 1 byte of data
	lenD++
	data := mergeUint64(lenD, b[atByte:atByte+lenD])
	atByte += lenD
	return int64(data), atByte
}

func PackUint32(num uint32, pos int) []byte {
	lenD, data := splitUint32(num)
	typLenD := mergeDataTypeAndLenDataLen(4, byte(lenD))
	result := []byte{byte(pos), typLenD}
	result = append(result, data...)
	return result
}
func UnpackUint32(b []byte, atByte int, lenD int) (uint32, int) {
	// This is required because lenD == 0 represents 1 byte of data
	lenD++
	data := mergeUint32(lenD, b[atByte:atByte+lenD])
	atByte += lenD
	return data, atByte
}

func PackUint64(num uint64, pos int) []byte {
	lenD, data := splitUint64(num)
	typLenD := mergeDataTypeAndLenDataLen(5, byte(lenD))
	result := []byte{byte(pos), typLenD}
	result = append(result, data...)
	return result
}
func UnpackUint64(b []byte, atByte int, lenD int) (uint64, int) {
	// This is required because lenD == 0 represents 1 byte of data
	lenD++
	data := mergeUint64(lenD, b[atByte:atByte+lenD])
	atByte += lenD
	return data, atByte
}

func PackFloat32(num float32, pos int) []byte {
	typLenD := mergeDataTypeAndLenDataLen(6, 3)
	result := []byte{byte(pos), typLenD}
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, math.Float32bits(num))
	result = append(result, data...)
	return result
}
func UnpackFloat32(b []byte, atByte int, lenD int) (float32, int) {
	lenD++
	dataBytes := b[atByte : atByte+lenD]
	atByte += lenD
	data := math.Float32frombits(
		binary.LittleEndian.Uint32(dataBytes),
	)
	return data, atByte
}

func PackFloat64(num float64, pos int) []byte {
	typLenD := mergeDataTypeAndLenDataLen(7, 7)
	result := []byte{byte(pos), typLenD}
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, math.Float64bits(num))
	result = append(result, data...)
	return result
}
func UnpackFloat64(b []byte, atByte int, lenD int) (float64, int) {
	lenD++
	dataBytes := b[atByte : atByte+lenD]
	atByte += lenD
	data := math.Float64frombits(
		binary.LittleEndian.Uint64(dataBytes),
	)
	return data, atByte
}

func PackBool(num bool, pos int) []byte {
	typLenD := mergeDataTypeAndLenDataLen(8, 0)
	result := []byte{byte(pos), typLenD}
	if num {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}
	return result
}
func UnpackBool(b []byte, atByte int) (bool, int) {
	dataBytes := b[atByte]
	atByte += 1
	var data bool
	if dataBytes > 0 {
		data = true
	} else {
		data = false
	}
	return data, atByte
}

func PackObject(data []byte, pos int) []byte {
	lenD, dataLenBytes := splitUint32(uint32(len(data)))
	typLenD := mergeDataTypeAndLenDataLen(9, byte(lenD))
	result := []byte{byte(pos), typLenD}
	result = append(result, dataLenBytes...)
	result = append(result, data...)
	return result
}
func UnpackObject(b []byte, atByte int, lenD int) ([]byte, int) {
	lenD++
	dataLen := int(mergeUint32(lenD, b[atByte:atByte+lenD]))
	atByte += lenD
	dataBytes := b[atByte : atByte+dataLen]
	return dataBytes, atByte + dataLen
}

func PackEnum(s string, pos int) []byte {
	data := []byte(s)
	lenD, dataLenBytes := splitUint32(uint32(len(data)))
	typLenD := mergeDataTypeAndLenDataLen(10, byte(lenD))
	result := []byte{byte(pos), typLenD}
	result = append(result, dataLenBytes...)
	result = append(result, data...)
	return result
}
func UnpackEnum(b []byte, atByte int, lenD int) (string, int) {
	lenD++
	dataLen := int(mergeUint32(lenD, b[atByte:atByte+lenD]))
	atByte += lenD
	dataBytes := b[atByte : atByte+dataLen]
	return string(dataBytes), atByte + dataLen
}

func PackBytes(s []byte, pos int) []byte {
	lenD, dataLenBytes := splitUint32(uint32(len(s)))
	typLenD := mergeDataTypeAndLenDataLen(11, byte(lenD))
	result := []byte{byte(pos), typLenD}
	result = append(result, dataLenBytes...)
	result = append(result, s...)
	return result
}
func UnpackBytes(b []byte, atByte int, lenD int) ([]byte, int) {
	lenD++
	dataLen := int(mergeUint32(lenD, b[atByte:atByte+lenD]))
	atByte += lenD
	dataBytes := b[atByte : atByte+dataLen]
	return dataBytes, atByte + dataLen
}

func mergeDataTypeAndLenDataLen(typ byte, lenD byte) byte {
	return lenD + (typ << 4)
}

func unmergeDataTypeAndLenDataLen(b byte) (int, int) {
	lenD := b & 0b00001111
	typ := (b & 0b11110000) >> 4
	return int(typ), int(lenD)
}

func splitUint32(unum uint32) (int, []byte) {
	var lenD int
	var b []byte
	if unum < 0xFF {
		lenD = 0
		b = []byte{byte(unum)}
	} else if unum < 0xFFFF {
		lenD = 1
		b = make([]byte, lenD+1)
		binary.LittleEndian.PutUint16(b, uint16(unum))
	} else if unum < 0xFFFFFF {
		lenD = 2
		b = []byte{
			byte((0xFF0000 & unum) >> 16),
			byte((0x00FF00 & unum) >> 8),
			byte(0x0000FF & unum),
		}
	} else {
		lenD = 3
		b = make([]byte, lenD+1)
		binary.LittleEndian.PutUint32(b, uint32(unum))
	}
	return lenD, b
}

func mergeUint32(lenD int, bytes []byte) uint32 {
	if lenD == 1 {
		return uint32(bytes[0])
	} else if lenD == 2 {
		return uint32(binary.LittleEndian.Uint16(bytes))
	} else if lenD == 3 {
		return uint32(bytes[0])<<16 +
			uint32(bytes[1])<<8 +
			uint32(bytes[2])
	} else {
		return binary.LittleEndian.Uint32(bytes)
	}
}

func splitUint64(unum uint64) (int, []byte) {
	var lenD int
	var b []byte
	if unum < 0xFF {
		lenD = 0
		b = []byte{byte(unum)}
	} else if unum < 0xFFFF {
		lenD = 1
		b = make([]byte, lenD+1)
		binary.LittleEndian.PutUint16(b, uint16(unum))
	} else if unum < 0xFFFFFF {
		lenD = 2
		b = []byte{
			byte((0xFF0000 & unum) >> 16),
			byte((0x00FF00 & unum) >> 8),
			byte(0x0000FF & unum),
		}
	} else if unum < 0xFFFFFFFF {
		lenD = 3
		b = make([]byte, lenD+1)
		binary.LittleEndian.PutUint32(b, uint32(unum))
	} else if unum < 0xFFFFFFFFFF {
		lenD = 4
		b = []byte{
			byte((0xFF00000000 & unum) >> 32),
			byte((0x00FF000000 & unum) >> 24),
			byte((0x0000FF0000 & unum) >> 16),
			byte((0x000000FF00 & unum) >> 8),
			byte(0x00000000FF & unum),
		}
	} else if unum < 0xFFFFFFFFFFFF {
		lenD = 5
		b = []byte{
			byte((0xFF0000000000 & unum) >> 40),
			byte((0x00FF00000000 & unum) >> 32),
			byte((0x0000FF000000 & unum) >> 24),
			byte((0x000000FF0000 & unum) >> 16),
			byte((0x00000000FF00 & unum) >> 8),
			byte(0x0000000000FF & unum),
		}
	} else if unum < 0xFFFFFFFFFFFFFF {
		lenD = 6
		b = []byte{
			byte((0xFF000000000000 & unum) >> 48),
			byte((0x00FF0000000000 & unum) >> 40),
			byte((0x0000FF00000000 & unum) >> 32),
			byte((0x000000FF000000 & unum) >> 24),
			byte((0x00000000FF0000 & unum) >> 16),
			byte((0x0000000000FF00 & unum) >> 8),
			byte(0x000000000000FF & unum),
		}
	} else {
		lenD = 7
		b = make([]byte, lenD+1)
		binary.LittleEndian.PutUint64(b, uint64(unum))
	}
	return lenD, b
}

func mergeUint64(lenD int, bytes []byte) uint64 {
	if lenD == 1 {
		return uint64(bytes[0])
	} else if lenD == 2 {
		return uint64(binary.LittleEndian.Uint16(bytes))
	} else if lenD == 3 {
		return uint64(bytes[0])<<16 +
			uint64(bytes[1])<<8 +
			uint64(bytes[2])
	} else if lenD == 4 {
		return uint64(binary.LittleEndian.Uint32(bytes))
	} else if lenD == 5 {
		return uint64(bytes[0])<<32 +
			uint64(bytes[1])<<24 +
			uint64(bytes[2])<<16 +
			uint64(bytes[3])<<8 +
			uint64(bytes[4])
	} else if lenD == 6 {
		return uint64(bytes[0])<<40 +
			uint64(bytes[1])<<32 +
			uint64(bytes[2])<<24 +
			uint64(bytes[3])<<16 +
			uint64(bytes[4])<<8 +
			uint64(bytes[5])
	} else if lenD == 7 {
		return uint64(bytes[0])<<48 +
			uint64(bytes[1])<<40 +
			uint64(bytes[2])<<32 +
			uint64(bytes[3])<<24 +
			uint64(bytes[4])<<16 +
			uint64(bytes[5])<<8 +
			uint64(bytes[6])
	} else {
		return binary.LittleEndian.Uint64(bytes)
	}
}
`
