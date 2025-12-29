import { test } from 'bun:test';
import { Obj1Arr } from '../dist/ts/obj_Obj1Arr';
import { Obj1 } from '../dist/ts/obj_Obj1';
import { Obj2 } from '../dist/ts/obj_Obj2';

// test('pack and unpack small obj', () => {
//     const ob = new Obj2({
//         key: 'key 1',
//         value: 'value 1',
//     });
//     let start = process.hrtime.bigint();
//     const s1 = JSON.stringify(ob, (key, value) => {
//         return typeof value === 'bigint' ? value.toString() : value;
//     });
//     console.log(
//         `JSON pack:`,
//         process.hrtime.bigint() - start,
//         '-> size:',
//         s1.length,
//     );
//     start = process.hrtime.bigint();
//     JSON.parse(s1);
//     console.log('JSON unpack:', process.hrtime.bigint() - start);
//     start = process.hrtime.bigint();
//     let buf = ob.pack();
//     console.log(
//         `MBin pack:`,
//         process.hrtime.bigint() - start,
//         '-> size:',
//         buf.length,
//         buf,
//     );
//     const ob2 = new Obj2({
//         key: 'key 2',
//         value: 'value 2',
//     });
//     start = process.hrtime.bigint();
//     buf = ob2.pack();
//     console.log(
//         `MBin pack 2:`,
//         process.hrtime.bigint() - start,
//         '-> size:',
//         buf.length,
//         buf,
//     );
//     // start = process.hrtime.bigint();
//     // const o2 = Obj2.unpack(ob);
//     // console.log('MBin unpack:', process.hrtime.bigint() - start);
// });

test('pack and unpack data', () => {
    const items = new Obj1Arr({
        items: [],
    });
    for (let i = 0; i < 10000; i++) {
        items.items.push(
            new Obj1({
                str: `Test string ${i}`,
                strArr: [
                    `Item iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii ${i}-1`,
                    `Item ${i}-2`,
                    `Item ${i}-3`,
                ],
                i32: i,
                i32Arr: [i, -i - 1, -i - 200000],
                i64: BigInt(i) * BigInt(100000),
                i64Arr: [BigInt(i), BigInt(-i - 1), BigInt(-i - 200000)],
                u32: i * 2,
                u32Arr: [i * 2, i * 3, i * 4],
                u64: BigInt(i) * BigInt(200000),
                u64Arr: [
                    BigInt(i) * BigInt(2),
                    BigInt(i) * BigInt(3),
                    BigInt(i) * BigInt(4),
                ],
                f32: i + 0.5,
                f32Arr: [i + 0.1, i + 0.2, i + 0.3],
                f64: i + 0.25,
                f64Arr: [i + 0.01, i + 0.02, i + 0.03],
                bool: i % 2 === 0,
                boolArr: [true, false, true],
                enum1: 'E1',
                enum1Arr: ['E1', 'E2', 'E3'],
                obj2: new Obj2({
                    key: `Obj2 key ${i}`,
                    value: i * 10.1,
                }),
                obj2Arr: [
                    new Obj2({
                        key: `Obj2 key ${i}-1`,
                        value: i * 7.2 + 1,
                    }),
                    new Obj2({
                        key: `Obj2 key ${i}-2`,
                        value: i * 8.54003 + 2,
                    }),
                ],
            }),
        );
    }
    for (let i = 0; i < 10; i++) {
        let start = process.hrtime.bigint();
        const s1 = JSON.stringify(items, (_key, value) => {
            return typeof value === 'bigint' ? value.toString() : value;
        });
        console.debug(
            `JSON pack:`,
            (process.hrtime.bigint() - start) / 1000000n,
            'ms -> size:',
            (s1.length / 8 / 1000000).toFixed(2),
        );
        // start = process.hrtime.bigint();
        // const _s1 = JSON.parse(s1);
        // console.debug(`JSON unpack ${i}:`, process.hrtime.bigint() - start);
        start = process.hrtime.bigint();
        const buf = items.pack();
        console.debug(
            `MBin pack ${i}:`,
            (process.hrtime.bigint() - start) / 1000000n,
            'ms -> size:',
            (buf.length / 8 / 1000000).toFixed(2),
            '\n',
        );
        // start = process.hrtime.bigint();
        // Obj1Arr.unpack(buf);
        // console.debug(`MBin unpack ${i}:`, process.hrtime.bigint() - start);
    }
});
