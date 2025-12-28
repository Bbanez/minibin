package parser_ts

const Common = `
export interface UnpackageEntry {
    setPropAtPos(pos: number, value: unknown): void;
}

export class Minibin {
    static unpack<T extends UnpackageEntry>(
        o: T,
        bytes: number[],
    ): Error | null {
        let atByte = 0;
        while (atByte < bytes.length) {
            const pos = bytes[atByte];
            atByte += 1;
            const [typ, lenD] = unmergeDataTypeAndLenDataLen(bytes[atByte]);
            atByte += 1;
            switch (typ) {
                case 0: {
                    atByte += 1;
                    continue;
                }
                case 1:
                    {
                        const [data, next] = this.unpackString(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 2:
                    {
                        const [data, next] = this.unpackInt32(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 3:
                    {
                        const [data, next] = this.unpackInt64(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 4:
                    {
                        const [data, next] = this.unpackUint32(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 5:
                    {
                        const [data, next] = this.unpackUint64(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 6:
                    {
                        const [data, next] = this.unpackFloat32(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 7:
                    {
                        const [data, next] = this.unpackFloat64(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 8:
                    {
                        const [data, next] = this.unpackBool(bytes, atByte);
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 9:
                    {
                        const [data, next] = this.unpackObject(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                case 10:
                    {
                        const [data, next] = this.unpackEnum(
                            bytes,
                            atByte,
                            lenD,
                        );
                        atByte = next;
                        o.setPropAtPos(pos, data);
                    }
                    break;
                default: {
                    return new Error('Unknown data type: ' + typ);
                }
            }
        }
        return null;
    }

    static packString(buffer: number[], s: string, pos: number): void {
        const data = Array.from(new TextEncoder().encode(s));
        const [lenD, dataLenBytes] = splitUint32(data.length);
        const typLenD = mergeDataTypeAndLenDataLen(1, lenD);
        buffer.push(pos, typLenD, ...dataLenBytes, ...data);
    }
    static unpackString(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [string, number] {
        lenD++;
        const dataLen = mergeUint32(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        const dataBytes = buffer.slice(atByte, atByte + dataLen);
        return [
            new TextDecoder().decode(new Uint8Array(dataBytes)),
            atByte + dataLen,
        ];
    }

    static packInt32(buffer: number[], num: number, pos: number): void {
        const [lenD, data] = splitUint32(uint32(num));
        const typLenD = mergeDataTypeAndLenDataLen(2, lenD);
        buffer.push(pos, typLenD, ...data);
    }
    static unpackInt32(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [number, number] {
        // This is required because lenD == 0 represents 1 byte of data
        lenD++;
        const data = mergeUint32(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        return [int32(data), atByte];
    }

    static packInt64(buffer: number[], num: bigint, pos: number): void {
        const [lenD, data] = splitUint64(uint64(num));
        const typLenD = mergeDataTypeAndLenDataLen(3, lenD);
        buffer.push(pos, typLenD, ...data);
    }
    static unpackInt64(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [bigint, number] {
        // This is required because lenD == 0 represents 1 byte of data
        lenD++;
        const data = mergeUint64(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        return [int64(data), atByte];
    }

    static packUint32(buffer: number[], num: number, pos: number): void {
        const [lenD, data] = splitUint32(num);
        const typLenD = mergeDataTypeAndLenDataLen(4, lenD);
        buffer.push(pos, typLenD, ...data);
    }
    static unpackUint32(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [number, number] {
        // This is required because lenD == 0 represents 1 byte of data
        lenD++;
        const data = mergeUint32(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        return [data, atByte];
    }

    static packUint64(buffer: number[], num: bigint, pos: number): void {
        const [lenD, data] = splitUint64(num);
        const typLenD = mergeDataTypeAndLenDataLen(5, lenD);
        buffer.push(pos, typLenD, ...data);
    }
    static unpackUint64(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [bigint, number] {
        // This is required because lenD == 0 represents 1 byte of data
        lenD++;
        const data = mergeUint64(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        return [data, atByte];
    }

    static packFloat32(buffer: number[], num: number, pos: number): void {
        const typLenD = mergeDataTypeAndLenDataLen(6, 3);
        const arr = new ArrayBuffer(4);
        new DataView(arr).setFloat32(0, num, true);
        const data = Array.from(new Uint8Array(arr));
        buffer.push(pos, typLenD, ...data);
    }
    static unpackFloat32(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [number, number] {
        lenD++;
        const dataBytes = buffer.slice(atByte, atByte + lenD);
        atByte += lenD;
        const arr = new Uint8Array(dataBytes);
        const data = new DataView(arr.buffer).getFloat32(0, true); // true for little-endian
        return [parseFloat(data.toFixed(5)), atByte];
    }

    static packFloat64(buffer: number[], num: number, pos: number): void {
        const typLenD = mergeDataTypeAndLenDataLen(7, 7);
        const arr = new ArrayBuffer(8);
        new DataView(arr).setFloat64(0, num, true);
        const data = Array.from(new Uint8Array(arr));
        buffer.push(pos, typLenD, ...data);
    }
    static unpackFloat64(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [number, number] {
        lenD++;
        const dataBytes = buffer.slice(atByte, atByte + lenD);
        atByte += lenD;
        const arr = new Uint8Array(dataBytes);
        const view = new DataView(arr.buffer);
        const data = view.getFloat64(0, true); // true for little-endian
        return [data, atByte];
    }

    static packBool(buffer: number[], num: boolean, pos: number): void {
        const typLenD = mergeDataTypeAndLenDataLen(8, 0);
        buffer.push(pos, typLenD);
        if (num) {
            buffer.push(1);
        } else {
            buffer.push(0);
        }
    }
    static unpackBool(buffer: number[], atByte: number): [boolean, number] {
        const dataBytes = buffer[atByte];
        atByte += 1;
        let data: boolean;
        if (dataBytes > 0) {
            data = true;
        } else {
            data = false;
        }
        return [data, atByte];
    }

    static packObject(buffer: number[], data: number[], pos: number): void {
        const [lenD, dataLenBytes] = splitUint32(uint32(data.length));
        const typLenD = mergeDataTypeAndLenDataLen(9, lenD);
        buffer.push(pos, typLenD, ...dataLenBytes, ...data);
    }
    static unpackObject(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [number[], number] {
        lenD++;
        const dataLen = mergeUint32(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        const dataBytes = buffer.slice(atByte, atByte + dataLen);
        return [dataBytes, atByte + dataLen];
    }

    static packEnum(buffer: number[], s: string, pos: number): void {
        const data = Array.from(new TextEncoder().encode(s));
        const [lenD, dataLenBytes] = splitUint32(data.length);
        const typLenD = mergeDataTypeAndLenDataLen(10, lenD);
        buffer.push(pos, typLenD, ...dataLenBytes, ...data);
    }
    static unpackEnum(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [string, number] {
        lenD++;
        const dataLen = mergeUint32(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        const dataBytes = buffer.slice(atByte, atByte + dataLen);
        return [
            new TextDecoder().decode(new Uint8Array(dataBytes)),
            atByte + dataLen,
        ];
    }
}

function uint32(num: number): number {
    return num >>> 0;
}
function int32(num: number): number {
    return num | 0;
}

function uint64(num: bigint): bigint {
    return num & 0xffffffffffffffffn;
}
function int64(num: bigint): bigint {
    const maxInt64 = BigInt('0x7FFFFFFFFFFFFFFF');
    if (num > maxInt64) {
        return num - BigInt('0x10000000000000000');
    }
    return num;
}

function mergeDataTypeAndLenDataLen(typ: number, lenD: number): number {
    return (lenD & 0xff) + ((typ & 0xff) << 4);
}

function unmergeDataTypeAndLenDataLen(b: number): [number, number] {
    const lenD = b & 0b00001111;
    const typ = (b & 0b11110000) >> 4;
    return [typ, lenD];
}

function splitUint32(unum: number): [number, number[]] {
    let lenD: number;
    let b: number[];
    if (unum < 0xff) {
        lenD = 0;
        b = [unum];
    } else if (unum < 0xffff) {
        lenD = 1;
        b = [(0xff00 & unum) >> 8, 0x00ff & unum];
    } else if (unum < 0xffffff) {
        lenD = 2;
        b = [(0xff0000 & unum) >> 16, (0x00ff00 & unum) >> 8, 0x0000ff & unum];
    } else {
        lenD = 3;
        b = [
            (0xff000000 & unum) >> 24,
            (0x00ff0000 & unum) >> 16,
            (0x0000ff00 & unum) >> 8,
            0x000000ff & unum,
        ];
    }
    return [lenD, b];
}

function mergeUint32(lenD: number, bytes: number[]): number {
    if (lenD == 1) {
        return bytes[0];
    } else if (lenD == 2) {
        return bytes[1] << (8 + bytes[0]);
    } else if (lenD == 3) {
        return (bytes[0] << (16 + bytes[1])) << (8 + bytes[2]);
    } else {
        return (
            ((bytes[0] << (24 + bytes[1])) << (16 + bytes[2])) << (8 + bytes[3])
        );
    }
}

function splitUint64(unum: bigint): [number, number[]] {
    if (unum < 0xff) {
        return [0, [Number(unum) & 0xff]];
    } else if (unum < 0xffff) {
        const num = Number(unum);
        return [1, [(0xff00 & num) >> 8, 0x00ff & num]];
    } else if (unum < 0xffffff) {
        const num = Number(unum);
        return [
            2,
            [(0xff0000 & num) >> 16, (0x00ff00 & num) >> 8, 0x0000ff & num],
        ];
    } else if (unum < 0xffffffff) {
        const num = Number(unum);
        return [
            3,
            [
                (0xff000000 & num) >> 24,
                (0x00ff0000 & num) >> 16,
                (0x0000ff00 & num) >> 8,
                0x000000ff & num,
            ],
        ];
    } else if (unum < 0xffffffffff) {
        return [
            4,
            [
                Number((unum >> BigInt(32)) & 0xffn),
                Number((unum >> BigInt(24)) & 0xffn),
                Number((unum >> BigInt(16)) & 0xffn),
                Number((unum >> BigInt(8)) & 0xffn),
                Number(unum & 0xffn),
            ],
        ];
    } else if (unum < 0xffffffffffff) {
        return [
            5,
            [
                Number((unum >> BigInt(40)) & 0xffn),
                Number((unum >> BigInt(32)) & 0xffn),
                Number((unum >> BigInt(24)) & 0xffn),
                Number((unum >> BigInt(16)) & 0xffn),
                Number((unum >> BigInt(8)) & 0xffn),
                Number(unum & 0xffn),
            ],
        ];
    } else if (unum < 0xffffffffffffffn) {
        return [
            6,
            [
                Number((unum >> BigInt(48)) & 0xffn),
                Number((unum >> BigInt(40)) & 0xffn),
                Number((unum >> BigInt(32)) & 0xffn),
                Number((unum >> BigInt(24)) & 0xffn),
                Number((unum >> BigInt(16)) & 0xffn),
                Number((unum >> BigInt(8)) & 0xffn),
                Number(unum & 0xffn),
            ],
        ];
    } else {
        return [
            7,
            [
                Number((unum >> BigInt(52)) & 0xffn),
                Number((unum >> BigInt(48)) & 0xffn),
                Number((unum >> BigInt(40)) & 0xffn),
                Number((unum >> BigInt(32)) & 0xffn),
                Number((unum >> BigInt(24)) & 0xffn),
                Number((unum >> BigInt(16)) & 0xffn),
                Number((unum >> BigInt(8)) & 0xffn),
                Number(unum & 0xffn),
            ],
        ];
    }
}

function mergeUint64(lenD: number, bytes: number[]): bigint {
    if (lenD == 1) {
        return BigInt(bytes[0]);
    } else if (lenD == 2) {
        return BigInt(bytes[1] << (8 + bytes[0]));
    } else if (lenD == 3) {
        return BigInt((bytes[0] << 16) + (bytes[1] << 8) + bytes[2]);
    } else if (lenD == 4) {
        return BigInt(
            (bytes[0] << 24) + (bytes[0] << 16) + (bytes[1] << 8) + bytes[2],
        );
    } else if (lenD == 5) {
        return (
            (BigInt(bytes[0]) << BigInt(32)) +
            (BigInt(bytes[1]) << BigInt(24)) +
            (BigInt(bytes[2]) << BigInt(16)) +
            (BigInt(bytes[3]) << BigInt(8)) +
            BigInt(bytes[4])
        );
    } else if (lenD == 6) {
        return (
            (BigInt(bytes[0]) << BigInt(40)) +
            (BigInt(bytes[1]) << BigInt(32)) +
            (BigInt(bytes[2]) << BigInt(24)) +
            (BigInt(bytes[3]) << BigInt(16)) +
            (BigInt(bytes[4]) << BigInt(8)) +
            BigInt(bytes[5])
        );
    } else if (lenD == 7) {
        return (
            (BigInt(bytes[0]) << BigInt(48)) +
            (BigInt(bytes[1]) << BigInt(40)) +
            (BigInt(bytes[2]) << BigInt(32)) +
            (BigInt(bytes[3]) << BigInt(24)) +
            (BigInt(bytes[4]) << BigInt(16)) +
            (BigInt(bytes[5]) << BigInt(8)) +
            BigInt(bytes[6])
        );
    } else {
        return (
            (BigInt(bytes[0]) << BigInt(52)) +
            (BigInt(bytes[1]) << BigInt(48)) +
            (BigInt(bytes[2]) << BigInt(40)) +
            (BigInt(bytes[3]) << BigInt(32)) +
            (BigInt(bytes[4]) << BigInt(24)) +
            (BigInt(bytes[5]) << BigInt(16)) +
            (BigInt(bytes[6]) << BigInt(8)) +
            BigInt(bytes[7])
        );
    }
}
`
