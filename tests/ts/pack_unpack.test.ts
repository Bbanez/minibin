import { expect, test } from 'bun:test';
import { Obj1 } from '../dist/ts/obj_Obj1';
import { Obj2 } from '../dist/ts/obj_Obj2';

test('pack and unpack data', () => {
    const obj = new Obj1({
        str: `Test string`,
        strArr: [`Item 1`, `Item 2`, `Item 3`],
        i32: 1,
        i32Arr: [2, -3 - 1, -200000],
        i64: BigInt(-100000101),
        i64Arr: [BigInt(1), BigInt(-1), BigInt(-200000)],
        u32: 2,
        u32Arr: [2, 3, 4],
        u64: BigInt(200000),
        u64Arr: [BigInt(2), BigInt(3), BigInt(4)],
        f32: 0.5,
        f32Arr: [0.1, 0.2, 0.3],
        f64: 1928340.25,
        f64Arr: [0.01, 0.02, 0.03],
        bool: true,
        boolArr: [true, false, true],
        enum1: 'E1',
        enum1Arr: ['E1', 'E2', 'E3'],
        obj2: new Obj2({
            key: `Obj2 key`,
            value: 10.10001,
        }),
        obj2Arr: [
            new Obj2({
                key: `Obj2 key 1`,
                value: 8.2,
            }),
            new Obj2({
                key: `Obj2 key 2`,
                value: 8.54003,
            }),
        ],
    });
    let start = process.hrtime.bigint();
    const s1 = JSON.stringify(obj, (_key, value) => {
        return typeof value === 'bigint' ? value.toString() : value;
    });
    console.debug(
        `JSON pack:`,
        (process.hrtime.bigint() - start) / 1000n,
        'us -> size:',
        (s1.length / 8).toFixed(2),
        'B',
    );
    start = process.hrtime.bigint();
    const buf = obj.pack();
    console.debug(
        `MBin pack:`,
        (process.hrtime.bigint() - start) / 1000n,
        'us -> size:',
        (buf.length / 8).toFixed(2),
        'B\n',
    );
    start = process.hrtime.bigint();
    const [s2, err] = Obj1.unpack(buf);
    console.debug(
        `MBin unpack:`,
        (process.hrtime.bigint() - start) / 1000n,
        'us',
    );
    if (err) {
        throw err;
    }
    const s3 = JSON.stringify(s2, (_key, value) => {
        return typeof value === 'bigint' ? value.toString() : value;
    });
    expect(s1).toEqual(s3);
});
