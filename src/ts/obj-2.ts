export interface Obj2Type {
    key: string;
    value: number;
}

export class Obj2 {
    key: string;
    value: number;

    constructor(data: Obj2Type) {
        this.key = data.key;
        this.value = data.value;
    }

    getKey(): string {
        return this.key;
    }
    setKey(v string): void {
        this.key = v
    }
}
