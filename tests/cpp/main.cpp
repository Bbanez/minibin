#include <cstdint>
#include <cstdio>
#include <string>

#include "minibin.hpp"
int main() {
    printf("Hello, World!\n");
    // ObjS objS = ObjS(3000000000001);
    // printf("%s\n", objS.print().c_str());
    // std::vector<uint8_t> packedS = objS.pack();
    // printf("Packed data size: %zu bytes\n", packedS.size());
    // for (size_t i = 0; i < packedS.size(); i++) {
    //     printf("%02X ", packedS[i]);
    // }
    // printf("\n");
    // ObjS unpackedObjS = unpackObjS(packedS);
    // printf("%s\n", unpackedObjS.print().c_str());
    // printf("value1: %ld == value2: %ld = %d\n", objS.value,
    //        unpackedObjS.value, objS.value == unpackedObjS.value);

    Obj1 obj1 = Obj1(
        "Hello, World!", std::vector<std::string>{"Hello", "World"},
        new int32_t(12345), std::vector<int32_t>{1, 2, 3}, 123456789012345,
        std::vector<int64_t>{123456789012345, 123456789012346}, 1234567890,
        std::vector<uint32_t>{1234567890, 1234567891}, 12345678901234567890ULL,
        std::vector<uint64_t>{12345678901234567890ULL, 12345678901234567891ULL},
        3.14f, std::vector<float>{3.14f, 2.71f}, 3.141592653589793,
        std::vector<double>{3.141592653589793, 2.718281828459045}, true,
        std::vector<bool>{true, false, true}, Obj2("Nested Object", 42),
        std::vector<Obj2>{Obj2("Nested Object 1", 42),
                          Obj2("Nested Object 2", 43)},
        Enum1::E1, std::vector<Enum1>{Enum1::E2, Enum1::E3});
    printf("%s\n", obj1.print().c_str());
    std::vector<uint8_t> packed = obj1.pack();
    printf("Packed data size: %zu bytes\n", packed.size());
    Obj1 unpackedObj1 = unpackObj1(packed);
    printf("%s\n", unpackedObj1.print().c_str());
    return 0;
}
