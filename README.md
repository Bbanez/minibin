# minibin

A minimal and primitive version of Protobuf concept without dependencies.

## Features

- 🚀 Zero dependencies
- 📦 Lightweight binary serialization
- 🔧 Simple schema-based encoding/decoding
- 🎯 Support for common data types (integers, floats, strings, bytes, booleans)
- 🔄 Compatible with Protobuf wire format concepts

## Installation

```bash
npm install minibin
```

## Usage

### Basic Example

```typescript
import { encode, decode, FieldType, Schema } from 'minibin';

// Define a schema
const userSchema: Schema = {
  id: { type: FieldType.INT32, tag: 1 },
  name: { type: FieldType.STRING, tag: 2 },
  active: { type: FieldType.BOOL, tag: 3 },
};

// Create an object
const user = {
  id: 42,
  name: 'Alice',
  active: true,
};

// Encode to binary
const encoded = encode(userSchema, user);
console.log(encoded); // Uint8Array

// Decode from binary
const decoded = decode(userSchema, encoded);
console.log(decoded); // { id: 42, name: 'Alice', active: true }
```

### Supported Field Types

- `FieldType.INT32` - 32-bit signed integer
- `FieldType.INT64` - 64-bit signed integer
- `FieldType.UINT32` - 32-bit unsigned integer
- `FieldType.UINT64` - 64-bit unsigned integer
- `FieldType.BOOL` - Boolean value
- `FieldType.FLOAT` - 32-bit floating point
- `FieldType.DOUBLE` - 64-bit floating point
- `FieldType.STRING` - UTF-8 string
- `FieldType.BYTES` - Byte array (Uint8Array)

### Schema Definition

A schema is a plain object where each key represents a field name, and the value defines the field type and tag number:

```typescript
const schema: Schema = {
  fieldName: { type: FieldType.STRING, tag: 1 },
  // tag numbers should be unique and start from 1
};
```

### Complex Example

```typescript
import { encode, decode, FieldType, Schema } from 'minibin';

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
  timestamp: Date.now(),
  sender: 'user@example.com',
  content: 'Hello, World!',
  priority: 1.5,
  encrypted: false,
  attachment: new Uint8Array([1, 2, 3, 4]),
};

const encoded = encode(messageSchema, message);
const decoded = decode(messageSchema, encoded);
```

### Partial Messages

You can encode objects with only some fields defined. Undefined or null fields are automatically skipped:

```typescript
const user = {
  id: 42,
  name: 'Bob',
  // active field is not set
};

const encoded = encode(userSchema, user);
const decoded = decode(userSchema, encoded);
// decoded.active will be undefined
```

### Forward Compatibility

The decoder automatically skips unknown fields, allowing for schema evolution:

```typescript
const oldSchema: Schema = {
  id: { type: FieldType.INT32, tag: 1 },
  name: { type: FieldType.STRING, tag: 2 },
};

const newSchema: Schema = {
  id: { type: FieldType.INT32, tag: 1 },
  name: { type: FieldType.STRING, tag: 2 },
  email: { type: FieldType.STRING, tag: 3 }, // new field
};

// Data encoded with newSchema can be decoded with oldSchema
// The email field will be silently ignored
```

## API Reference

### `encode(schema: Schema, obj: any): Uint8Array`

Serializes an object according to the provided schema into a binary format.

**Parameters:**
- `schema`: The schema definition
- `obj`: The object to encode

**Returns:** `Uint8Array` containing the binary data

### `decode(schema: Schema, buffer: Uint8Array): any`

Deserializes binary data according to the provided schema.

**Parameters:**
- `schema`: The schema definition
- `buffer`: The binary data to decode

**Returns:** The decoded object

### Low-Level Functions

The library also exports low-level encoding/decoding functions:

- `encodeVarint(value: number): Uint8Array`
- `decodeVarint(buffer: Uint8Array, offset: number): { value: number; length: number }`
- `encodeFloat(value: number): Uint8Array`
- `decodeFloat(buffer: Uint8Array, offset: number): number`
- `encodeDouble(value: number): Uint8Array`
- `decodeDouble(buffer: Uint8Array, offset: number): number`
- `encodeString(value: string): Uint8Array`
- `decodeString(buffer: Uint8Array, offset: number): { value: string; length: number }`
- `encodeBytes(value: Uint8Array): Uint8Array`
- `decodeBytes(buffer: Uint8Array, offset: number): { value: Uint8Array; length: number }`

## How It Works

minibin uses a wire format similar to Protocol Buffers:

1. **Varint Encoding**: Variable-length encoding for integers
2. **Wire Types**: Different encoding methods for different data types
   - Varint (0): int32, int64, uint32, uint64, bool
   - Fixed64 (1): double
   - Length-delimited (2): string, bytes
   - Fixed32 (5): float
3. **Tag-Value Pairs**: Each field is encoded as a key (tag + wire type) followed by the value

## Limitations

- No support for nested messages (yet)
- No support for repeated fields (arrays)
- No support for maps
- 32-bit integer limitations for very large numbers
- No built-in schema validation

## License

MIT