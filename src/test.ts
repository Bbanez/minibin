import { strict as assert } from 'assert';
import { test } from 'node:test';
import {
  encode,
  decode,
  encodeVarint,
  decodeVarint,
  encodeFloat,
  decodeFloat,
  encodeDouble,
  decodeDouble,
  encodeString,
  decodeString,
  encodeBytes,
  decodeBytes,
  encodeKey,
  decodeKey,
  FieldType,
  WireType,
  Schema,
} from './index';

test('encodeVarint - small numbers', () => {
  const result = encodeVarint(0);
  assert.deepEqual(result, new Uint8Array([0]));
  
  const result2 = encodeVarint(1);
  assert.deepEqual(result2, new Uint8Array([1]));
  
  const result3 = encodeVarint(127);
  assert.deepEqual(result3, new Uint8Array([127]));
});

test('encodeVarint - larger numbers', () => {
  const result = encodeVarint(300);
  assert.deepEqual(result, new Uint8Array([0xac, 0x02]));
});

test('decodeVarint - small numbers', () => {
  const result = decodeVarint(new Uint8Array([0]), 0);
  assert.deepEqual(result, { value: 0, length: 1 });
  
  const result2 = decodeVarint(new Uint8Array([127]), 0);
  assert.deepEqual(result2, { value: 127, length: 1 });
});

test('decodeVarint - larger numbers', () => {
  const result = decodeVarint(new Uint8Array([0xac, 0x02]), 0);
  assert.deepEqual(result, { value: 300, length: 2 });
});

test('encodeKey and decodeKey', () => {
  const key = encodeKey(1, WireType.VARINT);
  const keyValue = decodeVarint(key, 0).value;
  const decoded = decodeKey(keyValue);
  assert.deepEqual(decoded, { tag: 1, wireType: WireType.VARINT });
});

test('encodeFloat and decodeFloat', () => {
  const encoded = encodeFloat(3.14);
  assert.equal(encoded.length, 4);
  
  const decoded = decodeFloat(encoded, 0);
  assert.ok(Math.abs(decoded - 3.14) < 0.001);
});

test('encodeDouble and decodeDouble', () => {
  const encoded = encodeDouble(3.14159265359);
  assert.equal(encoded.length, 8);
  
  const decoded = decodeDouble(encoded, 0);
  assert.ok(Math.abs(decoded - 3.14159265359) < 0.0000001);
});

test('encodeString and decodeString', () => {
  const encoded = encodeString('hello');
  const decoded = decodeString(encoded, 0);
  assert.equal(decoded.value, 'hello');
  assert.equal(decoded.length, encoded.length);
});

test('encodeString and decodeString - empty string', () => {
  const encoded = encodeString('');
  const decoded = decodeString(encoded, 0);
  assert.equal(decoded.value, '');
});

test('encodeString and decodeString - unicode', () => {
  const encoded = encodeString('こんにちは');
  const decoded = decodeString(encoded, 0);
  assert.equal(decoded.value, 'こんにちは');
});

test('encodeBytes and decodeBytes', () => {
  const data = new Uint8Array([1, 2, 3, 4, 5]);
  const encoded = encodeBytes(data);
  const decoded = decodeBytes(encoded, 0);
  assert.deepEqual(decoded.value, data);
  assert.equal(decoded.length, encoded.length);
});

test('encode and decode - simple message with int', () => {
  const schema: Schema = {
    id: { type: FieldType.INT32, tag: 1 },
  };
  
  const obj = { id: 42 };
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.deepEqual(decoded, obj);
});

test('encode and decode - message with string', () => {
  const schema: Schema = {
    name: { type: FieldType.STRING, tag: 1 },
  };
  
  const obj = { name: 'Alice' };
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.deepEqual(decoded, obj);
});

test('encode and decode - message with bool', () => {
  const schema: Schema = {
    active: { type: FieldType.BOOL, tag: 1 },
  };
  
  const obj = { active: true };
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.deepEqual(decoded, obj);
});

test('encode and decode - message with float', () => {
  const schema: Schema = {
    value: { type: FieldType.FLOAT, tag: 1 },
  };
  
  const obj = { value: 3.14 };
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.ok(Math.abs(decoded.value - obj.value) < 0.001);
});

test('encode and decode - message with double', () => {
  const schema: Schema = {
    value: { type: FieldType.DOUBLE, tag: 1 },
  };
  
  const obj = { value: 3.14159265359 };
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.ok(Math.abs(decoded.value - obj.value) < 0.0000001);
});

test('encode and decode - message with bytes', () => {
  const schema: Schema = {
    data: { type: FieldType.BYTES, tag: 1 },
  };
  
  const obj = { data: new Uint8Array([1, 2, 3, 4, 5]) };
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.deepEqual(decoded.data, obj.data);
});

test('encode and decode - complex message', () => {
  const schema: Schema = {
    id: { type: FieldType.INT32, tag: 1 },
    name: { type: FieldType.STRING, tag: 2 },
    active: { type: FieldType.BOOL, tag: 3 },
    score: { type: FieldType.FLOAT, tag: 4 },
  };
  
  const obj = {
    id: 123,
    name: 'Bob',
    active: false,
    score: 98.6,
  };
  
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.equal(decoded.id, obj.id);
  assert.equal(decoded.name, obj.name);
  assert.equal(decoded.active, obj.active);
  assert.ok(Math.abs(decoded.score - obj.score) < 0.001);
});

test('encode and decode - partial message', () => {
  const schema: Schema = {
    id: { type: FieldType.INT32, tag: 1 },
    name: { type: FieldType.STRING, tag: 2 },
    active: { type: FieldType.BOOL, tag: 3 },
  };
  
  const obj = {
    id: 123,
    name: 'Charlie',
  };
  
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.equal(decoded.id, obj.id);
  assert.equal(decoded.name, obj.name);
  assert.equal(decoded.active, undefined);
});

test('encode and decode - skip unknown fields', () => {
  const encodeSchema: Schema = {
    id: { type: FieldType.INT32, tag: 1 },
    name: { type: FieldType.STRING, tag: 2 },
    extra: { type: FieldType.INT32, tag: 3 },
  };
  
  const decodeSchema: Schema = {
    id: { type: FieldType.INT32, tag: 1 },
    name: { type: FieldType.STRING, tag: 2 },
  };
  
  const obj = {
    id: 456,
    name: 'David',
    extra: 999,
  };
  
  const encoded = encode(encodeSchema, obj);
  const decoded = decode(decodeSchema, encoded);
  
  assert.equal(decoded.id, obj.id);
  assert.equal(decoded.name, obj.name);
  assert.equal(decoded.extra, undefined);
});

test('encode - ignores null and undefined values', () => {
  const schema: Schema = {
    id: { type: FieldType.INT32, tag: 1 },
    name: { type: FieldType.STRING, tag: 2 },
  };
  
  const obj = {
    id: 789,
    name: null,
  };
  
  const encoded = encode(schema, obj);
  const decoded = decode(schema, encoded);
  
  assert.equal(decoded.id, obj.id);
  assert.equal(decoded.name, undefined);
});

console.log('All tests passed!');
