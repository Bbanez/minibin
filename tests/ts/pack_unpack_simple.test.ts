import { test } from 'bun:test';
import { ObjS } from '../dist/ts/obj_ObjS';

test('pack and unpack data', () => {
    const obj = new ObjS({
        f64: -0.01,
    });
    const s1 = JSON.stringify(obj, (_key, value) => {
        return typeof value === 'bigint' ? value.toString() : value;
    });
    console.debug('Original object JSON:', s1);
    const packed = obj.pack();
    console.debug('Packed data length:', packed);
    const [unpacked, err] = ObjS.unpack(packed);
    if (err) {
        throw err;
    }
    const s2 = JSON.stringify(unpacked, (_key, value) => {
        return typeof value === 'bigint' ? value.toString() : value;
    });
    console.debug('Unpacked object JSON:', s2);
});
