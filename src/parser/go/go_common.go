package parser_go

const Common string = `package minibin

import (
	"fmt"
	"strings"
	"strconv"
)

const maxUint32 = ^uint32(0)

func uint32Length(length int) uint32 {
	if uint64(length) > uint64(maxUint32) {
		panic("minibin payload exceeds uint32 length limit")
	}
	return uint32(length)
}

func int32SliceToStringSlice(items []int32) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.Itoa(int(v))
    }
    return strs
}
func int32SliceToStringSliceRef(items []*int32) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.Itoa(int(*v))
    }
    return strs
}

func int64SliceToStringSlice(items []int64) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatInt(v, 10)
    }
    return strs
}
func int64SliceToStringSliceRef(items []*int64) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatInt(*v, 10)
    }
    return strs
}

func uint32SliceToStringSlice(items []uint32) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatUint(uint64(v), 10)
    }
    return strs
}
func uint32SliceToStringSliceRef(items []*uint32) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatUint(uint64(*v), 10)
    }
    return strs
}

func uint64SliceToStringSlice(items []uint64) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatUint(v, 10)
    }
    return strs
}
func uint64SliceToStringSliceRef(items []*uint64) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatUint(*v, 10)
    }
    return strs
}

func float32SliceToStringSlice(items []float32) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatFloat(float64(v), 'f', -1, 32)
    }
    return strs
}
func float32SliceToStringSliceRef(items []*float32) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatFloat(float64(*v), 'f', -1, 32)
    }
    return strs
}

func float64SliceToStringSlice(items []float64) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatFloat(v, 'f', -1, 64)
    }
    return strs
}
func float64SliceToStringSliceRef(items []*float64) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatFloat(*v, 'f', -1, 64)
    }
    return strs
}

func boolSliceToStringSlice(items []bool) []string {
    strs := make([]string, len(items))
    for i, v := range items {
	    if v {
	        strs[i] = "true"
	    } else {
	        strs[i] = "false"
	    }
    }
    return strs
}
func boolSliceToStringSliceRef(items []*bool) []string {
    strs := make([]string, len(items))
    for i, v := range items {
	    if *v {
	        strs[i] = "true"
	    } else {
	        strs[i] = "false"
	    }
    }
    return strs
}

type EnumGen interface {
	ToStr() string
}
func enumSliceToStringSlice[E EnumGen](items []E) []string {
	strs := make([]string, len(items))
	for i, v := range items {
		strs[i] = v.ToStr()
	}
	return strs
}

func enumSliceToJSONStringSlice[E EnumGen](items []E) []string {
	strs := make([]string, len(items))
	for i, v := range items {
		strs[i] = jsonQuote(v.ToStr())
	}
	return strs
}

type ObjGen interface {
	ToJson() string
}
func objSliceToStringSlice[O ObjGen](items []O) []string {
	strs := make([]string, len(items))
	for i, v := range items {
		strs[i] = v.ToJson()
	}
	return strs
}

func bytesToStringSlice(items []byte) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = strconv.FormatUint(uint64(v), 10)
    }
    return strs
}

func cloneBytes(items []byte) []byte {
	if items == nil {
		return nil
	}
	output := make([]byte, len(items))
	copy(output, items)
	return output
}
func bytesToStringSliceRef(items *[]byte) []string {
	if items == nil {
		return []string{}
	}
    strs := make([]string, len(*items))
    for i, v := range *items {
        strs[i] = strconv.FormatUint(uint64(v), 10)
    }
    return strs
}

func bytesSliceToStringSlice(items [][]byte) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = "[" + strings.Join(bytesToStringSlice(v), ",") + "]"
    }
    return strs
}
func bytesSliceToStringSliceRef(items []*[]byte) []string {
    strs := make([]string, len(items))
    for i, v := range items {
        strs[i] = "[" + strings.Join(bytesToStringSliceRef(v), ",") + "]"
    }
    return strs
}

func joinStringPointers(ptrs []*string, sep string) string {
    strs := make([]string, len(ptrs))
    for i, p := range ptrs {
        if p != nil {
            strs[i] = *p
        } else {
            strs[i] = "" // or handle nil pointers as needed
        }
    }
    return strings.Join(strs, sep)
}

func quoteStrings(items []string) []string {
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = jsonQuote(item)
	}
	return quoted
}

func quoteStringPointers(items []*string) []string {
	quoted := make([]string, len(items))
	for i, item := range items {
		if item == nil {
			quoted[i] = "null"
		} else {
			quoted[i] = jsonQuote(*item)
		}
	}
	return quoted
}

func jsonQuote(value string) string {
	for i := 0; i < len(value); i++ {
		if value[i] < 0x20 || value[i] == '"' || value[i] == '\\' {
			var output strings.Builder
			output.Grow(len(value) + 2)
			output.WriteByte('"')
			for j := 0; j < len(value); j++ {
				switch value[j] {
				case '"', '\\':
					output.WriteByte('\\')
					output.WriteByte(value[j])
				case '\b':
					output.WriteString("\\b")
				case '\f':
					output.WriteString("\\f")
				case '\n':
					output.WriteString("\\n")
				case '\r':
					output.WriteString("\\r")
				case '\t':
					output.WriteString("\\t")
				default:
					if value[j] < 0x20 {
						output.WriteString("\\u00")
						output.WriteByte("0123456789abcdef"[value[j]>>4])
						output.WriteByte("0123456789abcdef"[value[j]&0x0f])
					} else {
						output.WriteByte(value[j])
					}
				}
			}
			output.WriteByte('"')
			return output.String()
		}
	}
	return "\"" + value + "\""
}

func Compress(data []byte) ([]byte, error) {
	return data, nil
}

func Decompress(compressed []byte) ([]byte, error) {
	return compressed, nil
}

type UnpackabeEntry interface {
	SetPropAtPos(pos int, value any, level string) error
	GetPropNameAtPos(pos int) string
}

func Unpack[T UnpackabeEntry](o T, b []byte, level string) error {
	return unpack(o, b, level, 0)
}

const maxUnpackDepth = 64

func unpack[T UnpackabeEntry](o T, b []byte, level string, depth int) error {
	if depth > maxUnpackDepth {
		return fmt.Errorf("maximum nesting depth exceeded at %s", level)
	}
	bytes, err := Decompress(b)
	if err != nil {
		return err
	}
	atByte := 0
	for atByte < len(bytes) {
		if len(bytes)-atByte < 2 {
			return fmt.Errorf("truncated entry header at byte %d", atByte)
		}
		pos := int(bytes[atByte])
		atByte += 1
		typ, lenD := unmergeDataTypeAndLenDataLen(bytes[atByte])
		atByte += 1
		propName := o.GetPropNameAtPos(pos)
		lvl := level + "." + propName
		length := lenD + 1
		remaining := len(bytes) - atByte
		switch typ {
		case 0:
			if atByte >= len(bytes) {
				return fmt.Errorf("truncated data for %s", lvl)
			}
			atByte++
			continue
		case 1:
			if length > 4 || length > remaining {
				return fmt.Errorf("invalid or truncated length for %s", lvl)
			}
			dataLen := mergeUint32(length, bytes[atByte:atByte+length])
			if uint64(dataLen) > uint64(remaining-length) {
				return fmt.Errorf("truncated data for %s", lvl)
			}
			data, next := UnpackString(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 2:
			if length > 4 || length+1 > remaining {
				return fmt.Errorf("invalid or truncated data for %s", lvl)
			}
			data, next := UnpackInt32(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 3:
			if length > 8 || length+1 > remaining {
				return fmt.Errorf("invalid or truncated data for %s", lvl)
			}
			data, next := UnpackInt64(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 4:
			if length > 4 || length > remaining {
				return fmt.Errorf("invalid or truncated data for %s", lvl)
			}
			data, next := UnpackUint32(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 5:
			if length > 8 || length > remaining {
				return fmt.Errorf("invalid or truncated data for %s", lvl)
			}
			data, next := UnpackUint64(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 6:
			if length > 4 || length+1 > remaining {
				return fmt.Errorf("invalid or truncated data for %s", lvl)
			}
			data, next := UnpackInt32(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 7:
			if length > 8 || length+1 > remaining {
				return fmt.Errorf("invalid or truncated data for %s", lvl)
			}
			data, next := UnpackInt64(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 8:
			if lenD != 0 || atByte >= len(bytes) {
				return fmt.Errorf("invalid data length for %s", lvl)
			}
			data, next := UnpackBool(bytes, atByte)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 9:
			if length > 4 || length > remaining {
				return fmt.Errorf("invalid or truncated length for %s", lvl)
			}
			dataLen := mergeUint32(length, bytes[atByte:atByte+length])
			if uint64(dataLen) > uint64(remaining-length) {
				return fmt.Errorf("truncated data for %s", lvl)
			}
			data, next := UnpackObject(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 10:
			if length > 4 || length > remaining {
				return fmt.Errorf("invalid or truncated length for %s", lvl)
			}
			dataLen := mergeUint32(length, bytes[atByte:atByte+length])
			if uint64(dataLen) > uint64(remaining-length) {
				return fmt.Errorf("truncated data for %s", lvl)
			}
			data, next := UnpackString(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		case 11:
			if length > 4 || length > remaining {
				return fmt.Errorf("invalid or truncated length for %s", lvl)
			}
			dataLen := mergeUint32(length, bytes[atByte:atByte+length])
			if uint64(dataLen) > uint64(remaining-length) {
				return fmt.Errorf("truncated data for %s", lvl)
			}
			data, next := UnpackBytes(bytes, atByte, lenD)
			atByte = next
			if err := o.SetPropAtPos(pos, data, lvl); err != nil { return err }
		default:
			return fmt.Errorf(
				"Invalid datatype %s",
				lvl,
			)
		}
	}
	return nil
}

func PackString(s string, pos int) []byte {
	data := []byte(s)
	lenD, dataLenBytes := SplitUint32(uint32Length(len(data)))
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
	var neg byte = 0
	if num < 0 {
		neg = 1
		num = -num
	}
	lenD, data := SplitUint32(uint32(num))
	typLenD := mergeDataTypeAndLenDataLen(2, byte(lenD))
	result := []byte{byte(pos), typLenD, neg}
	result = append(result, data...)
	return result
}
func UnpackInt32(b []byte, atByte int, lenD int) (int32, int) {
	// This is required because lenD == 0 represents 1 byte of data
	lenD++
	neg := b[atByte]
	atByte++
	data := mergeUint32(lenD, b[atByte:atByte+lenD])
	atByte += lenD
	if neg == 1 {
		return -int32(data), atByte
	}
	return int32(data), atByte
}

func PackInt64(num int64, pos int) []byte {
	var neg byte = 0
	if num < 0 {
		neg = 1
		num = -num
	}
	lenD, data := SplitUint64(uint64(num))
	typLenD := mergeDataTypeAndLenDataLen(3, byte(lenD))
	result := []byte{byte(pos), typLenD, neg}
	result = append(result, data...)
	return result
}
func UnpackInt64(b []byte, atByte int, lenD int) (int64, int) {
	// This is required because lenD == 0 represents 1 byte of data
	lenD++
	neg := b[atByte]
	atByte++
	data := mergeUint64(lenD, b[atByte:atByte+lenD])
	atByte += lenD
	if neg == 1 {
		return -int64(data), atByte
	}
	return int64(data), atByte
}

func PackUint32(num uint32, pos int) []byte {
	lenD, data := SplitUint32(num)
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
	lenD, data := SplitUint64(num)
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

func PackFloat32(fnum float32, pos int, decimals float32) []byte {
	num := int32(fnum * decimals)
	var neg byte = 0
	if num < 0 {
		neg = 1
		num = -num
	}
	lenD, data := SplitUint32(uint32(num))
	typLenD := mergeDataTypeAndLenDataLen(6, byte(lenD))
	result := []byte{byte(pos), typLenD, neg}
	result = append(result, data...)
	return result
}

func PackFloat64(fnum float64, pos int, decimals float32) []byte {
	num := int64(fnum * float64(decimals))
	var neg byte = 0
	if num < 0 {
		neg = 1
		num = -num
	}
	lenD, data := SplitUint64(uint64(num))
	typLenD := mergeDataTypeAndLenDataLen(7, byte(lenD))
	result := []byte{byte(pos), typLenD, neg}
	result = append(result, data...)
	return result
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
	lenD, dataLenBytes := SplitUint32(uint32Length(len(data)))
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
	lenD, dataLenBytes := SplitUint32(uint32Length(len(data)))
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
	lenD, dataLenBytes := SplitUint32(uint32Length(len(s)))
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
	dataBytes := cloneBytes(b[atByte : atByte+dataLen])
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

func SplitUint32(unum uint32) (int, []byte) {
	var lenD int
	var b []byte
	if unum < 0xFF {
		lenD = 0
		b = []byte{byte(unum)}
	} else if unum < 0xFFFF {
		lenD = 1
		b = []byte{
			byte((0xFF00 & unum) >> 8),
			byte(0x00FF & unum),
		}
	} else if unum < 0xFFFFFF {
		lenD = 2
		b = []byte{
			byte((0xFF0000 & unum) >> 16),
			byte((0x00FF00 & unum) >> 8),
			byte(0x0000FF & unum),
		}
	} else {
		lenD = 3
		b = []byte{
			byte((0xFF000000 & unum) >> 24),
			byte((0x00FF0000 & unum) >> 16),
			byte((0x0000FF00 & unum) >> 8),
			byte(0x000000FF & unum),
		}
	}
	return lenD, b
}

func mergeUint32(lenD int, bytes []byte) uint32 {
	if lenD == 1 {
		return uint32(bytes[0])
	} else if lenD == 2 {
		return uint32(bytes[0])<<8 +
			uint32(bytes[1])
	} else if lenD == 3 {
		return uint32(bytes[0])<<16 +
			uint32(bytes[1])<<8 +
			uint32(bytes[2])
	} else {
		return uint32(bytes[0])<<24 +
			uint32(bytes[1])<<16 +
			uint32(bytes[2])<<8 +
			uint32(bytes[3])
	}
}

func SplitUint64(unum uint64) (int, []byte) {
	var lenD int
	var b []byte
	if unum < 0xFF {
		lenD = 0
		b = []byte{byte(unum)}
	} else if unum < 0xFFFF {
		lenD = 1
		b = []byte{
			byte((0xFF00 & unum) >> 8),
			byte(0x00FF & unum),
		}
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
		b = []byte{
			byte((0xFF000000 & unum) >> 24),
			byte((0x00FF0000 & unum) >> 16),
			byte((0x0000FF00 & unum) >> 8),
			byte(0x000000FF & unum),
		}
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
		b = []byte{
			byte((0xFF00000000000000 & unum) >> 56),
			byte((0x00FF000000000000 & unum) >> 48),
			byte((0x0000FF0000000000 & unum) >> 40),
			byte((0x000000FF00000000 & unum) >> 32),
			byte((0x00000000FF000000 & unum) >> 24),
			byte((0x0000000000FF0000 & unum) >> 16),
			byte((0x000000000000FF00 & unum) >> 8),
			byte(0x00000000000000FF & unum),
		}
	}
	return lenD, b
}

func mergeUint64(lenD int, bytes []byte) uint64 {
	if lenD == 1 {
		return uint64(bytes[0])
	} else if lenD == 2 {
		return uint64(bytes[0])<<8 +
			uint64(bytes[1])
	} else if lenD == 3 {
		return uint64(bytes[0])<<16 +
			uint64(bytes[1])<<8 +
			uint64(bytes[2])
	} else if lenD == 4 {
		return uint64(bytes[0])<<24 +
			uint64(bytes[1])<<16 +
			uint64(bytes[2])<<8 +
			uint64(bytes[3])
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
		return uint64(bytes[0])<<56 +
			uint64(bytes[1])<<48 +
			uint64(bytes[2])<<40 +
			uint64(bytes[3])<<32 +
			uint64(bytes[4])<<24 +
			uint64(bytes[5])<<16 +
			uint64(bytes[6])<<8 +
			uint64(bytes[7])
	}
}
`
