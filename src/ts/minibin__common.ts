export interface UnpackageEntry {
    setPropAtPos(pos: number, value: unknown): void;
}

export function packString(buffer: number[], s: string, pos: number): void {
    const data = Array.from(new TextEncoder().encode(s));
    const [lenD, dataLenBytes] = splitUint32(data.length);
    const typLenD = mergeDataTypeAndLenDataLen(1, lenD);
    buffer.push(pos, typLenD, ...dataLenBytes, ...data);
}

function toUint32(num: number): number {
    return num >>> 0;
}

function toInt32(num: number): number {
    return num | 0;
}

function toUint64(num: bigint): bigint {
    return num & 0xffffffffffffffffn;
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
    } else if (unum < 0xffffffffffffff) {
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
