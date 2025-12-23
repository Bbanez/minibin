import { encode, decode, FieldType, Schema } from './index';

// Example 1: Simple user object
console.log('=== Example 1: Simple User ===');
const userSchema: Schema = {
  id: { type: FieldType.INT32, tag: 1 },
  name: { type: FieldType.STRING, tag: 2 },
  active: { type: FieldType.BOOL, tag: 3 },
};

const user = {
  id: 42,
  name: 'Alice',
  active: true,
};

const encodedUser = encode(userSchema, user);
console.log('Original:', user);
console.log('Encoded bytes:', encodedUser);
console.log('Encoded length:', encodedUser.length, 'bytes');

const decodedUser = decode(userSchema, encodedUser);
console.log('Decoded:', decodedUser);
console.log('Match:', JSON.stringify(user) === JSON.stringify(decodedUser));
console.log();

// Example 2: Complex message
console.log('=== Example 2: Complex Message ===');
const messageSchema: Schema = {
  id: { type: FieldType.UINT32, tag: 1 },
  timestamp: { type: FieldType.INT64, tag: 2 },
  sender: { type: FieldType.STRING, tag: 3 },
  content: { type: FieldType.STRING, tag: 4 },
  priority: { type: FieldType.FLOAT, tag: 5 },
  encrypted: { type: FieldType.BOOL, tag: 6 },
  attachment: { type: FieldType.BYTES, tag: 7 },
};

const message = {
  id: 12345,
  timestamp: 1703347200000,
  sender: 'user@example.com',
  content: 'Hello, World!',
  priority: 1.5,
  encrypted: false,
  attachment: new Uint8Array([0x48, 0x65, 0x6c, 0x6c, 0x6f]), // "Hello" in ASCII
};

const encodedMessage = encode(messageSchema, message);
console.log('Original:', message);
console.log('Encoded length:', encodedMessage.length, 'bytes');

const decodedMessage = decode(messageSchema, encodedMessage);
console.log('Decoded:', decodedMessage);
console.log();

// Example 3: Partial message (forward compatibility)
console.log('=== Example 3: Partial Message ===');
const partialUser = {
  id: 99,
  name: 'Bob',
  // active field is omitted
};

const encodedPartial = encode(userSchema, partialUser);
console.log('Original (partial):', partialUser);
console.log('Encoded length:', encodedPartial.length, 'bytes');

const decodedPartial = decode(userSchema, encodedPartial);
console.log('Decoded:', decodedPartial);
console.log('Missing field is undefined:', decodedPartial.active === undefined);
console.log();

// Example 4: Schema evolution (backward compatibility)
console.log('=== Example 4: Schema Evolution ===');
const oldSchema: Schema = {
  id: { type: FieldType.INT32, tag: 1 },
  name: { type: FieldType.STRING, tag: 2 },
};

const newSchema: Schema = {
  id: { type: FieldType.INT32, tag: 1 },
  name: { type: FieldType.STRING, tag: 2 },
  email: { type: FieldType.STRING, tag: 3 },
  age: { type: FieldType.INT32, tag: 4 },
};

const newObject = {
  id: 123,
  name: 'Charlie',
  email: 'charlie@example.com',
  age: 30,
};

// Encode with new schema
const encodedNew = encode(newSchema, newObject);
console.log('Encoded with new schema:', newObject);

// Decode with old schema (unknown fields are skipped)
const decodedOld = decode(oldSchema, encodedNew);
console.log('Decoded with old schema:', decodedOld);
console.log('Email and age fields ignored:', decodedOld.email === undefined && decodedOld.age === undefined);
console.log();

// Example 5: Different data types
console.log('=== Example 5: Different Data Types ===');
const dataTypesSchema: Schema = {
  int32Field: { type: FieldType.INT32, tag: 1 },
  floatField: { type: FieldType.FLOAT, tag: 2 },
  doubleField: { type: FieldType.DOUBLE, tag: 3 },
  boolField: { type: FieldType.BOOL, tag: 4 },
  stringField: { type: FieldType.STRING, tag: 5 },
  bytesField: { type: FieldType.BYTES, tag: 6 },
};

const dataTypesObj = {
  int32Field: 2147483647,
  floatField: 3.14,
  doubleField: 3.141592653589793,
  boolField: true,
  stringField: 'minibin 🚀',
  bytesField: new Uint8Array([255, 128, 64, 32, 16, 8, 4, 2, 1]),
};

const encodedTypes = encode(dataTypesSchema, dataTypesObj);
console.log('Original:', dataTypesObj);
console.log('Encoded length:', encodedTypes.length, 'bytes');

const decodedTypes = decode(dataTypesSchema, encodedTypes);
console.log('Decoded:', decodedTypes);
console.log();

console.log('=== All examples completed successfully! ===');
