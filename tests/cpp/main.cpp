#include <cassert>
#include <cstdint>
#include <string>
#include <vector>

#include "../dist/cpp/minibin.hpp"

void testRoundTrip() {
    Coordinate origin;
    origin.latitude = 0;
    origin.longitude = 0;

    Attachment attachment;
    attachment.fileName = "report.bin";
    attachment.content = {1, 2, 3};
    attachment.sizes = {3, 1024};

    AllTypesDocument document;
    document.id = "document-001";
    document.aliases = {"primary", "backup"};
    document.minimumI32 = INT32_MIN;
    document.i32Values = {-1, 0, 1};
    document.minimumI64 = INT64_MIN;
    document.i64Values = {-9000000000LL, 0, 9000000000LL};
    document.maximumU32 = UINT32_MAX;
    document.u32Values = {0, 1, UINT32_MAX};
    document.maximumU64 = UINT64_MAX;
    document.u64Values = {0, 1, UINT64_MAX};
    document.price = -123.45f;
    document.f32Values = {-2.5f, 0.0f, 3.125f};
    document.measurement = -987654.123456789;
    document.f64Values = {-1.234567, 0.0, 9.876543};
    document.enabled = true;
    document.flags = {true, false, true};
    document.state = Lifecycle::in_review;
    document.priorities = {Priority::low, Priority::critical};
    document.payload = {0, 1, 127, 128, 255};
    document.chunks = {{1}, {2, 3, 4}};
    document.origin = origin;
    document.attachments = {attachment};

    AllTypesDocument unpacked = unpackAllTypesDocument(document.pack());
    assert(unpacked.id == document.id);
    assert(unpacked.minimumI32 == document.minimumI32);
    assert(unpacked.minimumI64 == document.minimumI64);
    assert(unpacked.maximumU32 == document.maximumU32);
    assert(unpacked.maximumU64 == document.maximumU64);
    assert(unpacked.i32Values == document.i32Values);
    assert(unpacked.i64Values == document.i64Values);
    assert(unpacked.u32Values == document.u32Values);
    assert(unpacked.u64Values == document.u64Values);
    assert(unpacked.payload == document.payload);
    assert(unpacked.chunks == document.chunks);
    assert(unpacked.state == document.state);
    assert(unpacked.priorities == document.priorities);
    assert(unpacked.attachments.size() == 1);
    assert(unpacked.attachments[0].content == attachment.content);

}

void testMutationAndPropertyLookup() {
    AllTypesDocument document;
    document.id = "document-001";
    document.aliases = {"primary"};
    document.i32Values = {1};
    document.i64Values = {2};
    document.u32Values = {3};
    document.u64Values = {4};
    document.f32Values = {1.5f};
    document.f64Values = {2.5};
    document.flags = {true};
    document.priorities = {Priority::low};
    document.chunks = {{1}};
    Attachment attachment;
    attachment.fileName = "before.bin";
    attachment.content = {1};
    attachment.sizes = {1};
    document.attachments = {attachment};

    document.id = "mutated-document";
    document.enabled = false;
    document.flags = {false, true};
    document.state = Lifecycle::archived;
    document.payload = {9, 8, 7};
    document.origin.latitude = 7;
    document.attachments[0].content = {4, 5, 6};

    assert(document.getPropNameAtPos(0) == "id");
    assert(document.getPropNameAtPos(34) == "payload");
    assert(document.getPropNameAtPos(255) == "__unknown__[255]");

    AllTypesDocument mutated = unpackAllTypesDocument(document.pack());
    assert(mutated.id == "mutated-document");
    assert(!mutated.enabled);
    assert(mutated.flags == std::vector<bool>({false, true}));
    assert(mutated.state == Lifecycle::archived);
    assert(mutated.payload == std::vector<uint8_t>({9, 8, 7}));
    assert(mutated.origin.latitude == 7);
    assert(mutated.attachments[0].content == std::vector<uint8_t>({4, 5, 6}));

}

int main(int argc, char* argv[]) {
    const std::string test = argc > 1 ? argv[1] : "all";
    if (test == "round-trip" || test == "all") testRoundTrip();
    if (test == "mutation" || test == "all") testMutationAndPropertyLookup();
    return 0;
}
