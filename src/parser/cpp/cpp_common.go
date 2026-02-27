package parser_cpp

var HClass = `
// -----------------------------------------------------------------
// ---- @name Class
// -----------------------------------------------------------------
class @name {
public:
	@name();
	@name(@constructorArgs);
	@props

	std::string getPropNameAtPos(uint8_t pos);
    std::vector<uint8_t> pack();
	std::string print(int indent = 4);
};
@name unpack@name(const std::vector<uint8_t>& b, std::string level = "@name");`

var CClass = `
// -----------------------------------------------------------------
// ---- @name Class
// -----------------------------------------------------------------
@name::@name() {
@emptyConstructorArgs
}
@name::@name(@constructorArgs) {
@constructorBody
}
std::string @name::getPropNameAtPos(uint8_t pos) {
@posToPropName
	return "__unknown__[" + std::to_string(pos) + "]";
}
std::vector<uint8_t> @name::pack() {
	std::vector<uint8_t> result;
@packProps
	return result;
}
@name unpack@name(const std::vector<uint8_t>& b, std::string level) {
	@name result = @name();
	uint32_t atByte = 0;
	while (atByte < b.size()) {
        uint8_t pos = b[atByte];
        atByte++;
        Tuple<uint8_t, uint8_t> res = _unmergeDataTypeAndLenDataLen(b[atByte]);
		uint8_t lenD = res.a;
		uint8_t typ  = res.b;
        atByte++;
        std::string propName = result.getPropNameAtPos(pos);
        std::string lvl      = level + "." + propName;
@unpackProps
    }
	return result;
}
std::string @name::print(int indent) {
	std::string indentStr(indent, ' ');
@printPrep
	return std::string("@name {\n") + @printStr + "\n" + indentStr + "}";
}`

var HEnum = `
// -----------------------------------------------------------------
// ---- @name Enum
// -----------------------------------------------------------------
enum class @name {
@enumValues
};
std::string @nameToString(@name e);
@name @nameFromString(const std::string& s);
@name newEmpty@name();`

var CEnum = `
// -----------------------------------------------------------------
// ---- @name Enum
// -----------------------------------------------------------------
std::string @nameToString(@name e) {
	switch (e) {
@enumToStringCases
	default:
		return "__unknown__";
	}
}
@name @nameFromString(const std::string& s) {
@stringToEnumCases
	return (@name)0;
}`

var CommonFunctionsH = `#ifndef MINIBIN_H
#define MINIBIN_H

#include <cstdint>
#include <string>
#include <vector>

template <typename T, typename K>
class Tuple {
public:
    Tuple(T a, K b) : a(a), b(b) {}

    T a;
    K b;
};


// ---------------------------------------------------------------------------
// ---- Common functions for packing and unpacking data
// ---------------------------------------------------------------------------
uint8_t _mergeDataTypeAndLenDataLen(uint8_t typ, uint8_t lenD);
Tuple<uint8_t, uint8_t> _unmergeDataTypeAndLenDataLen(uint8_t b);

Tuple<uint8_t, std::vector<uint8_t>> _splitUint32(uint32_t unum);
uint32_t _mergeUint32(uint8_t lenD, const std::vector<uint8_t>& bytes);

Tuple<uint8_t, std::vector<uint8_t>> _splitUint64(uint64_t unum);
uint64_t _mergeUint64(uint8_t lenD, const std::vector<uint8_t>& bytes);

std::vector<uint8_t> _packString(const std::string& str, uint8_t pos);
Tuple<std::string, uint32_t> _unpackString(const std::vector<uint8_t>& b,
                                                uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packInt32(int32_t num, uint8_t pos);
Tuple<int32_t, uint32_t> _unpackInt32(const std::vector<uint8_t>& b,
                                           uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packInt64(int64_t num, uint8_t pos);
Tuple<int64_t, uint32_t> _unpackInt64(const std::vector<uint8_t>& b,
                                           uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packUint32(uint32_t num, uint8_t pos);
Tuple<uint32_t, uint32_t> _unpackUint32(const std::vector<uint8_t>& b,
                                             uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packUint64(uint64_t num, uint8_t pos);
Tuple<uint64_t, uint32_t> _unpackUint64(const std::vector<uint8_t>& b,
                                             uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packFloat32(float num, uint8_t pos, float decimals);

std::vector<uint8_t> _packFloat64(double num, uint8_t pos, double decimals);

std::vector<uint8_t> _packBool(bool num, uint8_t pos);
Tuple<bool, uint32_t> _unpackBool(const std::vector<uint8_t>& b,
                                       uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packObject(const std::vector<uint8_t>& data, uint8_t pos);
Tuple<std::vector<uint8_t>, uint32_t> _unpackObject(
    const std::vector<uint8_t>& b, uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packEnum(const std::string& s, uint8_t pos);
Tuple<std::string, uint32_t> _unpackEnum(const std::vector<uint8_t>& b,
                                              uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packBytes(const std::vector<uint8_t>& data, uint8_t pos);
Tuple<std::vector<uint8_t>, uint32_t> _unpackBytes(
	const std::vector<uint8_t>& b, uint32_t atByte, uint8_t lenD);
// ---------------------------------------------------------------------------`

var CommonFunctionsCPP = `#include "minibin.hpp"

#include <cstdint>
#include <cstdio>
#include <stdexcept>
#include <vector>
#include <string>

// ---------------------------------------------------------------------------
// ---- Common functions for packing and unpacking data
// ---------------------------------------------------------------------------
uint8_t _mergeDataTypeAndLenDataLen(uint8_t typ, uint8_t lenD) {
    return lenD + (typ << 4);
}

Tuple<uint8_t, uint8_t> _unmergeDataTypeAndLenDataLen(uint8_t b) {
    uint8_t lenD = b & 0b00001111;
    uint8_t typ  = (b & 0b11110000) >> 4;
    return Tuple<uint8_t, uint8_t>(lenD, typ);
}

Tuple<uint8_t, std::vector<uint8_t>> _splitUint32(uint32_t unum) {
    uint8_t lenD;
    std::vector<uint8_t> b;
    if (unum < 0xFF) {
        lenD = 0;
        b    = {uint8_t(unum)};
    } else if (unum < 0xFFFF) {
        lenD = 1;
        b    = {uint8_t((0xFF00 & unum) >> 8), uint8_t(0x00FF & unum)};
    } else if (unum < 0xFFFFFF) {
        lenD = 2;
        b    = {
            uint8_t((0xFF0000 & unum) >> 16),
            uint8_t((0x00FF00 & unum) >> 8),
            uint8_t(0x0000FF & unum),
        };
    } else {
        lenD = 3;
        b    = {
            uint8_t((0xFF000000 & unum) >> 24),
            uint8_t((0x00FF0000 & unum) >> 16),
            uint8_t((0x0000FF00 & unum) >> 8),
            uint8_t(0x000000FF & unum),
        };
    }
    return Tuple<uint8_t, std::vector<uint8_t>>(lenD, b);
}

uint32_t _mergeUint32(uint8_t lenD, const std::vector<uint8_t>& bytes) {
    if (lenD == 1) {
        return uint32_t(bytes[0]);
    } else if (lenD == 2) {
        return (uint32_t(bytes[0]) << 8) + uint32_t(bytes[1]);
    } else if (lenD == 3) {
        return (uint32_t(bytes[0]) << 16) + (uint32_t(bytes[1]) << 8) +
               uint32_t(bytes[2]);
    } else {
        return (uint32_t(bytes[0]) << 24) + (uint32_t(bytes[1]) << 16) +
               (uint32_t(bytes[2]) << 8) + uint32_t(bytes[3]);
    }
}

Tuple<uint8_t, std::vector<uint8_t>> _splitUint64(uint64_t unum) {
    uint8_t lenD;
    std::vector<uint8_t> b;
    if (unum < 0xFF) {
        lenD = 0;
        b    = {uint8_t(unum)};
    } else if (unum < 0xFFFF) {
        lenD = 1;
        b    = {
            uint8_t((0xFF00 & unum) >> 8),
            uint8_t(0x00FF & unum),
        };
    } else if (unum < 0xFFFFFF) {
        lenD = 2;
        b    = {
            uint8_t((0xFF0000 & unum) >> 16),
            uint8_t((0x00FF00 & unum) >> 8),
            uint8_t(0x0000FF & unum),
        };
    } else if (unum < 0xFFFFFFFF) {
        lenD = 3;
        b    = {
            uint8_t((0xFF000000 & unum) >> 24),
            uint8_t((0x00FF0000 & unum) >> 16),
            uint8_t((0x0000FF00 & unum) >> 8),
            uint8_t(0x000000FF & unum),
        };
    } else if (unum < 0xFFFFFFFFFF) {
        lenD = 4;
        b    = {
            uint8_t((0xFF00000000 & unum) >> 32),
            uint8_t((0x00FF000000 & unum) >> 24),
            uint8_t((0x0000FF0000 & unum) >> 16),
            uint8_t((0x000000FF00 & unum) >> 8),
            uint8_t(0x00000000FF & unum),
        };
    } else if (unum < 0xFFFFFFFFFFFF) {
        lenD = 5;
        b    = {
            uint8_t((0xFF0000000000 & unum) >> 40),
            uint8_t((0x00FF00000000 & unum) >> 32),
            uint8_t((0x0000FF000000 & unum) >> 24),
            uint8_t((0x000000FF0000 & unum) >> 16),
            uint8_t((0x00000000FF00 & unum) >> 8),
            uint8_t(0x0000000000FF & unum),
        };
    } else if (unum < 0xFFFFFFFFFFFFFF) {
        lenD = 6;
        b    = {
            uint8_t((0xFF000000000000 & unum) >> 48),
            uint8_t((0x00FF0000000000 & unum) >> 40),
            uint8_t((0x0000FF00000000 & unum) >> 32),
            uint8_t((0x000000FF000000 & unum) >> 24),
            uint8_t((0x00000000FF0000 & unum) >> 16),
            uint8_t((0x0000000000FF00 & unum) >> 8),
            uint8_t(0x000000000000FF & unum),
        };
    } else {
        lenD = 7;
        b    = {
            uint8_t((0xFF00000000000000 & unum) >> 56),
            uint8_t((0x00FF000000000000 & unum) >> 48),
            uint8_t((0x0000FF0000000000 & unum) >> 40),
            uint8_t((0x000000FF00000000 & unum) >> 32),
            uint8_t((0x00000000FF000000 & unum) >> 24),
            uint8_t((0x0000000000FF0000 & unum) >> 16),
            uint8_t((0x000000000000FF00 & unum) >> 8),
            uint8_t(0x00000000000000FF & unum),
        };
    }
    return Tuple<uint8_t, std::vector<uint8_t>>(lenD, b);
}

uint64_t _mergeUint64(uint8_t lenD, const std::vector<uint8_t>& bytes) {
    if (lenD == 1) {
        return uint64_t(bytes[0]);
    } else if (lenD == 2) {
        return (uint64_t(bytes[0]) << 8) + uint64_t(bytes[1]);
    } else if (lenD == 3) {
        return (uint64_t(bytes[0]) << 16) + (uint64_t(bytes[1]) << 8) +
               uint64_t(bytes[2]);
    } else if (lenD == 4) {
        return (uint64_t(bytes[0]) << 24) + (uint64_t(bytes[1]) << 16) +
               (uint64_t(bytes[2]) << 8) + uint64_t(bytes[3]);
    } else if (lenD == 5) {
        return (uint64_t(bytes[0]) << 32) + (uint64_t(bytes[1]) << 24) +
               (uint64_t(bytes[2]) << 16) + (uint64_t(bytes[3]) << 8) +
               uint64_t(bytes[4]);
    } else if (lenD == 6) {
        return (uint64_t(bytes[0]) << 40) + (uint64_t(bytes[1]) << 32) +
               (uint64_t(bytes[2]) << 24) + (uint64_t(bytes[3]) << 16) +
               (uint64_t(bytes[4]) << 8) + uint64_t(bytes[5]);
    } else if (lenD == 7) {
        return (uint64_t(bytes[0]) << 48) + (uint64_t(bytes[1]) << 40) +
               (uint64_t(bytes[2]) << 32) + (uint64_t(bytes[3]) << 24) +
               (uint64_t(bytes[4]) << 16) + (uint64_t(bytes[5]) << 8) +
               uint64_t(bytes[6]);
    } else {
        return (uint64_t(bytes[0]) << 56) + (uint64_t(bytes[1]) << 48) +
               (uint64_t(bytes[2]) << 40) + (uint64_t(bytes[3]) << 32) +
               (uint64_t(bytes[4]) << 24) + (uint64_t(bytes[5]) << 16) +
               (uint64_t(bytes[6]) << 8) + uint64_t(bytes[7]);
    }
}

std::vector<uint8_t> _packString(const std::string& s, uint8_t pos) {
    std::vector<uint8_t> data(s.begin(), s.end());
     Tuple<uint8_t, std::vector<uint8_t>> res  = _splitUint32(uint32_t(data.size()));
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(1, res.a);
    std::vector<uint8_t> result = {pos, typLenD};
    result.insert(result.end(), res.b.begin(), res.b.end());
    result.insert(result.end(), data.begin(), data.end());
    return result;
}
Tuple<std::string, uint32_t> _unpackString(const std::vector<uint8_t>& b,
                                                uint32_t atByte, uint8_t lenD) {
    lenD++;
    std::vector<uint8_t> data(b.begin() + atByte, b.begin() + atByte + lenD);
    uint32_t dataLen = _mergeUint32(lenD, data);
    atByte += lenD;
    std::vector<uint8_t> dataBytes = std::vector<uint8_t>(
        b.begin() + atByte, b.begin() + atByte + dataLen);
    std::string str(dataBytes.begin(), dataBytes.end());
    atByte += dataLen;
	return Tuple<std::string, uint32_t>(str, atByte);
}

std::vector<uint8_t> _packInt32(int32_t num, uint8_t pos) {
    uint8_t neg = 0;
    if (num < 0) {
        neg = 1;
        num = -num;
    }
    Tuple<uint8_t, std::vector<uint8_t>> res = _splitUint32(uint32_t(num));
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(2, res.a);
    std::vector<uint8_t> result = {pos, typLenD, neg};
    result.insert(result.end(), res.b.begin(),res.b.end());
    return result;
}
Tuple<int32_t, uint32_t> _unpackInt32(const std::vector<uint8_t>& b,
                                           uint32_t atByte, uint8_t lenD) {
    lenD++;
    uint8_t neg = b.at(atByte);
    atByte++;
    std::vector<uint8_t> dataBytes(b.begin() + atByte,
                                   b.begin() + atByte + lenD);
    uint32_t data = _mergeUint32(lenD, dataBytes);
    atByte += lenD;
    if (neg) {
        return Tuple<int32_t, uint32_t>(-int32_t(data), atByte);
    }
    return Tuple<int32_t, uint32_t>(int32_t(data), atByte);
}

std::vector<uint8_t> _packInt64(int64_t num, uint8_t pos) {
    uint8_t neg = 0;
    if (num < 0) {
        neg = 1;
        num = -num;
    }
    Tuple<uint8_t, std::vector<uint8_t>> res = _splitUint64(uint64_t(num));
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(3, res.a);
    std::vector<uint8_t> result = {pos, typLenD, neg};
    result.insert(result.end(),res.b.begin(),res.b.end());
    return result;
}
Tuple<int64_t, uint32_t> _unpackInt64(const std::vector<uint8_t>& b,
                                           uint32_t atByte, uint8_t lenD) {
    lenD++;
    uint8_t neg = b.at(atByte);
    atByte++;
    std::vector<uint8_t> dataBytes(b.begin() + atByte,
                                   b.begin() + atByte + lenD);
    uint64_t data = _mergeUint64(lenD, dataBytes);
    atByte += lenD;
    if (neg) {
        return Tuple<int64_t, uint32_t>(-int64_t(data), atByte);
    }
    return Tuple<int64_t, uint32_t>(int64_t(data), atByte);
}

std::vector<uint8_t> _packUint32(uint32_t num, uint8_t pos) {
    Tuple<uint8_t, std::vector<uint8_t>> res = _splitUint32(num);
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(4, res.a);
    std::vector<uint8_t> result = {pos, typLenD};
    result.insert(result.end(), res.b.begin(), res.b.end());
    return result;
}
Tuple<uint32_t, uint32_t> _unpackUint32(const std::vector<uint8_t>& b,
                                             uint32_t atByte, uint8_t lenD) {
    lenD++;
    std::vector<uint8_t> dataBytes(b.begin() + atByte,
                                   b.begin() + atByte + lenD);
    uint32_t data = _mergeUint32(lenD, dataBytes);
    atByte += lenD;
    return Tuple<uint32_t, uint32_t>(data, atByte);
}

std::vector<uint8_t> _packUint64(uint64_t num, uint8_t pos) {
    Tuple<uint8_t, std::vector<uint8_t>> res = _splitUint64(num);
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(5,res.a);
    std::vector<uint8_t> result = {pos, typLenD};
    result.insert(result.end(),res.b.begin(), res.b.end());
    return result;
}
Tuple<uint64_t, uint32_t> _unpackUint64(const std::vector<uint8_t>& b,
                                             uint32_t atByte, uint8_t lenD) {
    lenD++;
    std::vector<uint8_t> dataBytes(b.begin() + atByte,
                                   b.begin() + atByte + lenD);
    uint64_t data = _mergeUint64(lenD, dataBytes);
    atByte += lenD;
    return Tuple<uint64_t, uint32_t>(data, atByte);
}

std::vector<uint8_t> _packFloat32(float fnum, uint8_t pos, float decimals) {
    int32_t num = int32_t(fnum * decimals);
    uint8_t neg = 0;
    if (num < 0) {
        neg = 1;
        num = -num;
    }
    Tuple<uint8_t, std::vector<uint8_t>> res = _splitUint32(num);
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(6, res.a);
    std::vector<uint8_t> result = {pos, typLenD, neg};
    result.insert(result.end(), res.b.begin(), res.b.end());
    return result;
}

std::vector<uint8_t> _packFloat64(double fnum, uint8_t pos, double decimals) {
    int64_t num = fnum * decimals;
    uint8_t neg = 0;
    if (num < 0) {
        neg = 1;
        num = -num;
    }
    Tuple<uint8_t, std::vector<uint8_t>> res           = _splitUint64(num);
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(7, res.a);
    std::vector<uint8_t> result = {pos, typLenD, neg};
    result.insert(result.end(),res.b.begin(),res.b.end());
    return result;
}

std::vector<uint8_t> _packBool(bool num, uint8_t pos) {
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(8, 0);
    std::vector<uint8_t> result = {pos, typLenD};
    if (num) {
        result.push_back(1);
    } else {
        result.push_back(0);
    }
    return result;
}
Tuple<bool, uint32_t> _unpackBool(const std::vector<uint8_t>& b,
                                       uint32_t atByte, uint8_t lenD) {
    uint8_t dataBytes = b.at(atByte);
    atByte++;
    bool data;
    if (dataBytes == 0) {
        data = false;
    } else {
        data = true;
    }
    return Tuple<bool, uint32_t>(data, atByte);
}

std::vector<uint8_t> _packObject(const std::vector<uint8_t>& data,
                                 uint8_t pos) {
    Tuple<uint8_t, std::vector<uint8_t>> res   = _splitUint32(data.size());
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(9, res.a);
    std::vector<uint8_t> result = {pos, typLenD};
    result.insert(result.end(), res.b.begin(), res.b.end());
    result.insert(result.end(), data.begin(), data.end());
    return result;
}
Tuple<std::vector<uint8_t>, uint32_t> _unpackObject(
    const std::vector<uint8_t>& b, uint32_t atByte, uint8_t lenD) {
    lenD++;
    uint8_t dataLen = _mergeUint32(
        lenD,
        std::vector<uint8_t>(b.begin() + atByte, b.begin() + atByte + lenD));
    atByte += lenD;
    std::vector<uint8_t> dataBytes = std::vector<uint8_t>(
        b.begin() + atByte, b.begin() + atByte + dataLen);
    atByte += dataLen;
    return Tuple<std::vector<uint8_t>, uint32_t>(dataBytes, atByte);
}

std::vector<uint8_t> _packEnum(const std::string& s, uint8_t pos) {
    std::vector<uint8_t> data(s.begin(), s.end());
    Tuple<uint8_t, std::vector<uint8_t>> res = _splitUint32(data.size());
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(10, res.a);
    std::vector<uint8_t> result = {pos, typLenD};
    result.insert(result.end(), res.b.begin(), res.b.end());
    result.insert(result.end(), data.begin(), data.end());
    return result;
}
Tuple<std::string, uint32_t> _unpackEnum(const std::vector<uint8_t>& b,
                                              uint32_t atByte, uint8_t lenD) {
    lenD++;
    std::vector<uint8_t> data(b.begin() + atByte, b.begin() + atByte + lenD);
    uint32_t dataLen = _mergeUint32(lenD, data);
    atByte += lenD;
    std::vector<uint8_t> dataBytes = std::vector<uint8_t>(
        b.begin() + atByte, b.begin() + atByte + dataLen);
    std::string str(dataBytes.begin(), dataBytes.end());
    atByte += dataLen;
    return Tuple<std::string, uint32_t>(str, atByte);
}

std::vector<uint8_t> _packBytes(const std::vector<uint8_t>& data, uint8_t pos) {
    Tuple<uint8_t, std::vector<uint8_t>> res   = _splitUint32(data.size());
    uint8_t typLenD             = _mergeDataTypeAndLenDataLen(11, res.a);
    std::vector<uint8_t> result = {pos, typLenD};
    result.insert(result.end(), res.b.begin(), res.b.end());
    result.insert(result.end(), data.begin(), data.end());
    return result;
}
Tuple<std::vector<uint8_t>, uint32_t> _unpackBytes(
    const std::vector<uint8_t>& b, uint32_t atByte, uint8_t lenD) {
    lenD++;
    uint8_t dataLen = _mergeUint32(
        lenD,
        std::vector<uint8_t>(b.begin() + atByte, b.begin() + atByte + lenD));
    atByte += lenD;
    std::vector<uint8_t> dataBytes = std::vector<uint8_t>(
        b.begin() + atByte, b.begin() + atByte + dataLen);
    atByte += dataLen;
    return Tuple<std::vector<uint8_t>, uint32_t>(dataBytes, atByte);
}
// ---------------------------------------------------------------------------`
