import { Obj1Arr } from '../dist/ts/obj_Obj1Arr.ts';
import { Obj1 } from '../dist/ts/obj_Obj1.ts';
import { Enum1 } from '../dist/ts/enum_Enum1.ts';
import { Obj2 } from '../dist/ts/obj_Obj2.ts';

function main() {
    const items = new Obj1Arr({
        items: [],
    });
    const arr = new Uint8Array();
    for (let i = 0; i < 1; i++) {
        items.items.push(
            new Obj1({
                str: `Test string ${i}`,
                strArr: [`Item ${i}-1`, `Item ${i}-2`, `Item ${i}-3`],
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
                enum1: Enum1.e1,
                enum1Arr: [Enum1.e1, Enum1.e2, Enum1.e3],
                obj2: new Obj2({
                    key: `Obj2 key ${i}`,
                    value: i * 10,
                }),
                obj2Arr: [
                    new Obj2({
                        key: `Obj2 key ${i}-1`,
                        value: i * 10 + 1,
                    }),
                    new Obj2({
                        key: `Obj2 key ${i}-2`,
                        value: i * 10 + 2,
                    }),
                ],
            }),
        );
    }
    let start = Date.now();
    const s1 = JSON.stringify(items, (key, value) => {
        return typeof value === 'bigint' ? value.toString() : value;
    });
    console.log(`JSON pack:`, Date.now() - start);
    start = Date.now();
    const buf = items.pack();
    console.log(`MBin pack:`, Date.now() - start);
}
main();
