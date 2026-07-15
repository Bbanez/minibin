import { expect, test } from 'bun:test';
import { AllTypesDocument } from '../dist/ts/obj_AllTypesDocument';
import { AllTypesDocumentBatch } from '../dist/ts/obj_AllTypesDocumentBatch';
import { Attachment } from '../dist/ts/obj_Attachment';
import { AuditEntry } from '../dist/ts/obj_AuditEntry';
import { Coordinate } from '../dist/ts/obj_Coordinate';

function document(): AllTypesDocument {
    return new AllTypesDocument({
        id: 'document-001',
        title: 'A fully populated document',
        aliases: ['primary', 'backup'],
        optionalAliases: ['legacy'],
        minimumI32: -2000000000,
        optionalI32: 2000000000,
        i32Values: [-1, 0, 1],
        optionalI32Values: [-42, 42],
        minimumI64: -9000000000000n,
        optionalI64: 9000000000000n,
        i64Values: [-9000000000n, 0n, 9000000000n],
        optionalI64Values: [-7n, 7n],
        maximumU32: 4294967295,
        optionalU32: 42,
        u32Values: [0, 1, 4294967295],
        optionalU32Values: [9],
        maximumU64: 9000000000000n,
        optionalU64: 99n,
        u64Values: [0n, 1n, 9000000000000n],
        optionalU64Values: [11n],
        price: -123.45,
        optionalRatio: 0.12345,
        f32Values: [-2.5, 0, 3.125],
        optionalF32Values: [-1.2, 4.5],
        measurement: -987654.1234567,
        optionalMeasurement: 0.1256,
        f64Values: [-1.234567, 0, 9.876543],
        optionalF64Values: [12.34],
        enabled: true,
        optionalEnabled: false,
        flags: [true, false, true],
        optionalFlags: [false, true],
        state: 'in_review',
        priorities: ['low', 'critical'],
        payload: [0, 1, 127, 128, 255],
        optionalPayload: [9, 8, 7],
        chunks: [[], [1], [2, 3, 4]],
        optionalChunks: [[5, 6]],
        origin: new Coordinate({
            latitude: 51.5,
            longitude: -0.1,
            altitudeMeters: 35.25,
        }),
        optionalDestination: new Coordinate({
            latitude: -33.8,
            longitude: 151.2,
            altitudeMeters: undefined,
        }),
        attachments: [
            new Attachment({
                fileName: 'report.bin',
                content: [1, 2, 3],
                checksum: [4, 5],
                labels: ['binary', 'report'],
                sizes: [3n, 1024n],
            }),
        ],
        optionalAuditTrail: [
            new AuditEntry({
                sequence: 1n,
                actorId: -42n,
                message: 'published',
                previousState: 'draft',
                nextState: 'published',
                changedFields: ['state', 'title'],
                location: new Coordinate({
                    latitude: 40.7,
                    longitude: -74,
                    altitudeMeters: undefined,
                }),
            }),
        ],
    });
}

test('round trips every supported type and nested shape', () => {
    const source = document();
    const [unpacked, error] = AllTypesDocument.unpack(source.pack());
    expect(error).toBeNull();
    expect(unpacked).toEqual(source);
});

test('round trips a batch with empty required arrays', () => {
    const source = new AllTypesDocumentBatch({
        batchId: 'batch-001',
        documents: [document()],
        failedDocuments: [],
    });
    const [unpacked, error] = AllTypesDocumentBatch.unpack(source.pack());
    expect(error).toBeNull();
    expect(unpacked).toEqual(source);
});

test('getters, setters, copy, and post-mutation packing preserve values', () => {
    const source = document();
    source.setId('mutated-id');
    source.setTitle('mutated title');
    source.setAliases(['mutated', 'aliases']);
    source.setOptionalI32(-99);
    source.setState('archived');
    source.setPayload([9, 8, 7]);
    source.setOrigin(
        new Coordinate({ latitude: 3, longitude: 4, altitudeMeters: 5.5 }),
    );

    expect(source.getId()).toBe('mutated-id');
    expect(source.getTitle()).toBe('mutated title');
    expect(source.getAliases()).toEqual(['mutated', 'aliases']);
    expect(source.getOptionalI32()).toBe(-99);
    expect(source.getState()).toBe('archived');
    expect(source.getPayload()).toEqual([9, 8, 7]);
    expect(source.getOrigin().getLatitude()).toBe(3);

    const copy = source.copy();
    copy.setId('copy-id');
    copy.setAliases(['copy-only']);
    copy.getOrigin().setLatitude(99);
    expect(source.getId()).toBe('mutated-id');
    expect(source.getAliases()).toEqual(['mutated', 'aliases']);
    expect(source.getOrigin().getLatitude()).toBe(3);

    const [unpacked, error] = AllTypesDocument.unpack(source.pack());
    expect(error).toBeNull();
    expect(unpacked?.getId()).toBe('mutated-id');
    expect(unpacked?.getState()).toBe('archived');
    expect(unpacked?.getPayload()).toEqual([9, 8, 7]);
    expect(unpacked?.getOrigin().getLatitude()).toBe(3);
});
