package parser_ts

const Common = `export interface UnpackageEntry {
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
                        const [data, next] = this.unpackInt32(
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
                        const [data, next] = this.unpackInt64(
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
                case 11:
                    {
                        const [data, next] = this.unpackBytes(
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
        let neg = 0;
        if (num < 0) {
            neg = 1;
            num = -num;
        }
        const [lenD, data] = splitUint32(num);
        const typLenD = mergeDataTypeAndLenDataLen(2, lenD);
        buffer.push(pos, typLenD, neg, ...data);
    }
    static unpackInt32(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [number, number] {
        // This is required because lenD == 0 represents 1 byte of data
        lenD++;
        const neg = buffer[atByte];
        atByte++;
        const data = mergeUint32(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        return [neg == 1 ? -data : data, atByte];
    }

    static packInt64(buffer: number[], num: bigint, pos: number): void {
        let neg = 0;
        if (num < 0) {
            neg = 1;
            num = -num;
        }
        const [lenD, data] = splitUint64(num);
        const typLenD = mergeDataTypeAndLenDataLen(3, lenD);
        buffer.push(pos, typLenD, neg, ...data);
    }
    static unpackInt64(
        buffer: number[],
        atByte: number,
        lenD: number,
    ): [bigint, number] {
        // This is required because lenD == 0 represents 1 byte of data
        lenD++;
        const neg = buffer[atByte];
        atByte++;
        const data = mergeUint64(lenD, buffer.slice(atByte, atByte + lenD));
        atByte += lenD;
        return [neg == 1 ? -data : data, atByte];
    }

    static packUint32(buffer: number[], num: number, pos: number): void {
        const [lenD, data] = splitUint32(num >>> 0);
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
        return [data >>> 0, atByte];
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

    static packFloat32(buffer: number[], fnum: number, pos: number, decimals: number): void {
        let num = fnum * decimals
        let neg = 0;
        if (num < 0) {
            neg = 1;
            num = -num;
        }
        const [lenD, data] = splitUint32(num);
        const typLenD = mergeDataTypeAndLenDataLen(6, lenD);
        buffer.push(pos, typLenD, neg, ...data);
    }

    static packFloat64(buffer: number[], fnum: number, pos: number, decimals: number): void {
        let num = BigInt(fnum * decimals)
        let neg = 0;
        if (num < 0) {
            neg = 1;
            num = -num;
        }
        const [lenD, data] = splitUint64(num);
        const typLenD = mergeDataTypeAndLenDataLen(7, lenD);
        buffer.push(pos, typLenD, neg, ...data);
        return;
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
        const [lenD, dataLenBytes] = splitUint32(data.length);
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

    static packBytes(buffer: number[], data: number[], pos: number): void {
        const [lenD, dataLenBytes] = splitUint32(data.length);
        const typLenD = mergeDataTypeAndLenDataLen(11, lenD);
        buffer.push(pos, typLenD, ...dataLenBytes, ...data);
    }
    static unpackBytes(
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
}

export function mergeDataTypeAndLenDataLen(typ: number, lenD: number): number {
    return (lenD & 0xff) + ((typ & 0xff) << 4);
}

export function unmergeDataTypeAndLenDataLen(b: number): [number, number] {
    const lenD = b & 0b00001111;
    const typ = (b & 0b11110000) >> 4;
    return [typ, lenD];
}

export function splitUint32(unum: number): [number, number[]] {
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

export function mergeUint32(lenD: number, bytes: number[]): number {
    if (lenD == 1) {
        return bytes[0];
    } else if (lenD == 2) {
        return (bytes[0] << 8) + bytes[1];
    } else if (lenD == 3) {
        return (bytes[0] << 16) + (bytes[1] << 8) + bytes[2];
    } else {
        return (bytes[0] << 24) + (bytes[1] << 16) + (bytes[2] << 8) + bytes[3];
    }
}

export function splitUint64(unum: bigint): [number, number[]] {
    if (unum < 0xffn) {
        return [0, [Number(unum) & 0xff]];
    } else if (unum < 0xffffn) {
        const num = Number(unum);
        return [1, [(0xff00 & num) >> 8, 0x00ff & num]];
    } else if (unum < 0xffffffn) {
        const num = Number(unum);
        return [
            2,
            [(0xff0000 & num) >> 16, (0x00ff00 & num) >> 8, 0x0000ff & num],
        ];
    } else if (unum < 0xffffffffn) {
        return [
            3,
            [
                Number((0xff000000n & unum) >> 24n),
                Number((0x00ff0000n & unum) >> 16n),
                Number((0x0000ff00n & unum) >> 8n),
                Number(0x000000ffn & unum),
            ],
        ];
    } else if (unum < 0xffffffffffn) {
        return [
            4,
            [
                Number((0xFF00000000n & unum) >> 32n),
                Number((0x00FF000000n & unum) >> 24n),
                Number((0x0000FF0000n & unum) >> 16n),
                Number((0x000000FF00n & unum) >> 8n),
                Number(0x00000000FFn & unum),
            ],
        ];
    } else if (unum < 0xffffffffffffn) {
        return [
            5,
            [
                Number((0xFF0000000000n & unum) >> 40n),
                Number((0x00FF00000000n & unum) >> 32n),
                Number((0x0000FF000000n & unum) >> 24n),
                Number((0x000000FF0000n & unum) >> 16n),
                Number((0x00000000FF00n & unum) >> 8n),
                Number(0x0000000000FFn & unum),
            ],
        ];
    } else if (unum < 0xffffffffffffffn) {
        return [
            6,
            [
                Number((0xFF000000000000n & unum) >> 48n),
                Number((0x00FF0000000000n & unum) >> 40n),
                Number((0x0000FF00000000n & unum) >> 32n),
                Number((0x000000FF000000n & unum) >> 24n),
                Number((0x00000000FF0000n & unum) >> 16n),
                Number((0x0000000000FF00n & unum) >> 8n),
                Number(0x000000000000FFn & unum),
            ],
        ];
    } else {
        return [
            7,
            [
                Number((0xFF00000000000000n & unum) >> 56n),
                Number((0x00FF000000000000n & unum) >> 48n),
                Number((0x0000FF0000000000n & unum) >> 40n),
                Number((0x000000FF00000000n & unum) >> 32n),
                Number((0x00000000FF000000n & unum) >> 24n),
                Number((0x0000000000FF0000n & unum) >> 16n),
                Number((0x000000000000FF00n & unum) >> 8n),
                Number(0x00000000000000FFn & unum),
            ],
        ];
    }
}

export function mergeUint64(lenD: number, bytes: number[]): bigint {
    if (lenD == 1) {
        return BigInt(bytes[0]);
    } else if (lenD == 2) {
        return BigInt((bytes[0] << 8) + bytes[1]);
    } else if (lenD == 3) {
        return BigInt((bytes[0] << 16) + (bytes[1] << 8) + bytes[2]);
    } else if (lenD == 4) {
        return (
            (BigInt(bytes[0]) << 24n) +
            (BigInt(bytes[1]) << 16n) +
            (BigInt(bytes[2]) << 8n) +
            BigInt(bytes[3])
        );
    } else if (lenD == 5) {
        return (
            (BigInt(bytes[0]) << 32n) +
            (BigInt(bytes[1]) << 24n) +
            (BigInt(bytes[2]) << 16n) +
            (BigInt(bytes[3]) << 8n) +
            BigInt(bytes[4])
        );
    } else if (lenD == 6) {
        return (
            (BigInt(bytes[0]) << 40n) +
            (BigInt(bytes[1]) << 32n) +
            (BigInt(bytes[2]) << 24n) +
            (BigInt(bytes[3]) << 16n) +
            (BigInt(bytes[4]) << 8n) +
            BigInt(bytes[5])
        );
    } else if (lenD == 7) {
        return (
            (BigInt(bytes[0]) << 48n) +
            (BigInt(bytes[1]) << 40n) +
            (BigInt(bytes[2]) << 32n) +
            (BigInt(bytes[3]) << 24n) +
            (BigInt(bytes[4]) << 16n) +
            (BigInt(bytes[5]) << 8n) +
            BigInt(bytes[6])
        );
    } else {
        return (
            (BigInt(bytes[0]) << 52n) +
            (BigInt(bytes[1]) << 48n) +
            (BigInt(bytes[2]) << 40n) +
            (BigInt(bytes[3]) << 32n) +
            (BigInt(bytes[4]) << 24n) +
            (BigInt(bytes[5]) << 16n) +
            (BigInt(bytes[6]) << 8n) +
            BigInt(bytes[7])
        );
    }
}
`
