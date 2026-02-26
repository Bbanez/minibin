package parser_cpp

var HClass = `class @name {
public:
	@name(@constructorArgs);
	@props

	std::string getPropNameAtPos(uint8_t pos);
    std::vector<uint8_t> pack();
};
@name newEmpty@name();
@name unpack@name(const std::vector<uint8_t>* b, std::string* level);`

var CommonFunctionsH = `
uint8_t _mergeDataTypeAndLenDataLen(uint8_t typ, uint8_t lenD);
std::tuple<uint8_t, uint8_t> _unmergeDataTypeAndLenDataLen(uint8_t b);

std::tuple<uint8_t, std::vector<uint8_t>> _splitUint32(uint32_t unum);
uint32_t _mergeUint32(uint8_t lenD, const std::vector<uint8_t>& bytes);

std::tuple<uint8_t, std::vector<uint8_t>> _splitUint64(uint64_t unum);
uint64_t _mergeUint64(uint8_t lenD, const std::vector<uint8_t>& bytes);

std::vector<uint8_t> _packString(const std::string* str, uint8_t pos);
std::tuple<std::string, uint32_t> _unpackString(const std::vector<uint8_t>* b,
                                                uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packInt32(int32_t num, uint8_t pos);
std::tuple<int32_t, uint32_t> _unpackInt32(const std::vector<uint8_t>* b,
                                           uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packInt64(int64_t num, uint8_t pos);
std::tuple<int64_t, uint32_t> _unpackInt64(const std::vector<uint8_t>* b,
                                           uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packUint32(uint32_t num, uint8_t pos);
std::tuple<uint32_t, uint32_t> _unpackUint32(const std::vector<uint8_t>* b,
                                             uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packUint64(uint64_t num, uint8_t pos);
std::tuple<uint64_t, uint32_t> _unpackUint64(const std::vector<uint8_t>* b,
                                             uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packFloat32(float num, uint8_t pos, float decimals);

std::vector<uint8_t> _packFloat64(double num, uint8_t pos);

std::vector<uint8_t> _packBool(bool num, uint8_t pos);
std::tuple<bool, uint32_t> _unpackBool(const std::vector<uint8_t>* b,
                                       uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packObject(const std::vector<uint8_t>* data, uint8_t pos);
std::tuple<std::vector<uint8_t>, uint32_t> _unpackObject(
    const std::vector<uint8_t>* b, uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packEnum(const std::string* s, uint8_t pos);
std::tuple<std::string, uint32_t> _unpackEnum(const std::vector<uint8_t>* b,
                                              uint32_t atByte, uint8_t lenD);

std::vector<uint8_t> _packBytes(const std::vector<uint8_t>* data, uint8_t pos);
std::tuple<std::vector<uint8_t>, uint32_t> _unpackBytes(
    const std::vector<uint8_t>* b, uint32_t atByte, uint8_t lenD);
	`
