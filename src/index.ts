/**
 * Wire types used in the binary encoding
 */
export enum WireType {
  VARINT = 0,        // int32, int64, uint32, uint64, bool
  FIXED64 = 1,       // fixed64, double
  LENGTH_DELIMITED = 2, // string, bytes, embedded messages
  FIXED32 = 5,       // fixed32, float
}

/**
 * Field types supported by minibin
 */
export enum FieldType {
  INT32 = 'int32',
  INT64 = 'int64',
  UINT32 = 'uint32',
  UINT64 = 'uint64',
  BOOL = 'bool',
  FLOAT = 'float',
  DOUBLE = 'double',
  STRING = 'string',
  BYTES = 'bytes',
}

/**
 * Field definition in a schema
 */
export interface FieldDef {
  type: FieldType;
  tag: number;
}

/**
 * Schema definition for a message
 */
export interface Schema {
  [fieldName: string]: FieldDef;
}

/**
 * Encodes a varint (variable-length integer)
 */
export function encodeVarint(value: number): Uint8Array {
  const bytes: number[] = [];
  let n = value >>> 0; // Convert to unsigned 32-bit
  
  while (n >= 0x80) {
    bytes.push((n & 0x7f) | 0x80);
    n >>>= 7;
  }
  bytes.push(n);
  
  return new Uint8Array(bytes);
}

/**
 * Decodes a varint from a buffer
 * Note: Limited to 32-bit integers due to JavaScript number precision
 */
export function decodeVarint(buffer: Uint8Array, offset: number): { value: number; length: number } {
  let value = 0;
  let shift = 0;
  let length = 0;
  const maxBytes = 10; // Maximum varint length for 64-bit values
  
  while (offset + length < buffer.length && length < maxBytes) {
    const byte = buffer[offset + length];
    length++;
    
    // For values beyond 32 bits, JavaScript number will lose precision
    value |= (byte & 0x7f) << shift;
    
    if ((byte & 0x80) === 0) {
      break;
    }
    
    shift += 7;
  }
  
  return { value, length };
}

/**
 * Encodes a field key (tag + wire type)
 */
export function encodeKey(tag: number, wireType: WireType): Uint8Array {
  return encodeVarint((tag << 3) | wireType);
}

/**
 * Decodes a field key
 */
export function decodeKey(value: number): { tag: number; wireType: WireType } {
  return {
    tag: value >>> 3,
    wireType: value & 0x7,
  };
}

/**
 * Encodes a 32-bit float
 */
export function encodeFloat(value: number): Uint8Array {
  const buffer = new ArrayBuffer(4);
  new DataView(buffer).setFloat32(0, value, true); // little-endian
  return new Uint8Array(buffer);
}

/**
 * Decodes a 32-bit float
 */
export function decodeFloat(buffer: Uint8Array, offset: number): number {
  const view = new DataView(buffer.buffer, buffer.byteOffset + offset, 4);
  return view.getFloat32(0, true); // little-endian
}

/**
 * Encodes a 64-bit double
 */
export function encodeDouble(value: number): Uint8Array {
  const buffer = new ArrayBuffer(8);
  new DataView(buffer).setFloat64(0, value, true); // little-endian
  return new Uint8Array(buffer);
}

/**
 * Decodes a 64-bit double
 */
export function decodeDouble(buffer: Uint8Array, offset: number): number {
  const view = new DataView(buffer.buffer, buffer.byteOffset + offset, 8);
  return view.getFloat64(0, true); // little-endian
}

/**
 * Encodes a string
 */
export function encodeString(value: string): Uint8Array {
  const encoder = new TextEncoder();
  const bytes = encoder.encode(value);
  const length = encodeVarint(bytes.length);
  
  const result = new Uint8Array(length.length + bytes.length);
  result.set(length, 0);
  result.set(bytes, length.length);
  
  return result;
}

/**
 * Decodes a string
 */
export function decodeString(buffer: Uint8Array, offset: number): { value: string; length: number } {
  const { value: strLength, length: headerLength } = decodeVarint(buffer, offset);
  const decoder = new TextDecoder();
  const value = decoder.decode(buffer.slice(offset + headerLength, offset + headerLength + strLength));
  
  return { value, length: headerLength + strLength };
}

/**
 * Encodes bytes
 */
export function encodeBytes(value: Uint8Array): Uint8Array {
  const length = encodeVarint(value.length);
  
  const result = new Uint8Array(length.length + value.length);
  result.set(length, 0);
  result.set(value, length.length);
  
  return result;
}

/**
 * Decodes bytes
 */
export function decodeBytes(buffer: Uint8Array, offset: number): { value: Uint8Array; length: number } {
  const { value: bytesLength, length: headerLength } = decodeVarint(buffer, offset);
  const value = buffer.slice(offset + headerLength, offset + headerLength + bytesLength);
  
  return { value, length: headerLength + bytesLength };
}

/**
 * Serializes an object according to a schema
 */
export function encode(schema: Schema, obj: any): Uint8Array {
  const parts: Uint8Array[] = [];
  
  for (const [fieldName, fieldDef] of Object.entries(schema)) {
    const value = obj[fieldName];
    
    if (value === undefined || value === null) {
      continue;
    }
    
    let wireType: WireType;
    let encodedValue: Uint8Array;
    
    switch (fieldDef.type) {
      case FieldType.INT32:
      case FieldType.INT64:
      case FieldType.UINT32:
      case FieldType.UINT64:
        wireType = WireType.VARINT;
        encodedValue = encodeVarint(value);
        break;
        
      case FieldType.BOOL:
        wireType = WireType.VARINT;
        encodedValue = encodeVarint(value ? 1 : 0);
        break;
        
      case FieldType.FLOAT:
        wireType = WireType.FIXED32;
        encodedValue = encodeFloat(value);
        break;
        
      case FieldType.DOUBLE:
        wireType = WireType.FIXED64;
        encodedValue = encodeDouble(value);
        break;
        
      case FieldType.STRING:
        wireType = WireType.LENGTH_DELIMITED;
        encodedValue = encodeString(value);
        break;
        
      case FieldType.BYTES:
        wireType = WireType.LENGTH_DELIMITED;
        encodedValue = encodeBytes(value);
        break;
        
      default:
        throw new Error(`Unsupported field type: ${fieldDef.type}`);
    }
    
    const key = encodeKey(fieldDef.tag, wireType);
    parts.push(key);
    parts.push(encodedValue);
  }
  
  // Combine all parts
  const totalLength = parts.reduce((sum, part) => sum + part.length, 0);
  const result = new Uint8Array(totalLength);
  let offset = 0;
  
  for (const part of parts) {
    result.set(part, offset);
    offset += part.length;
  }
  
  return result;
}

/**
 * Deserializes binary data according to a schema
 */
export function decode(schema: Schema, buffer: Uint8Array): any {
  const result: any = {};
  
  // Create reverse mapping from tag to field name and type
  const tagToField: Map<number, { name: string; type: FieldType }> = new Map();
  for (const [fieldName, fieldDef] of Object.entries(schema)) {
    tagToField.set(fieldDef.tag, { name: fieldName, type: fieldDef.type });
  }
  
  let offset = 0;
  
  while (offset < buffer.length) {
    // Read the key
    const { value: keyValue, length: keyLength } = decodeVarint(buffer, offset);
    offset += keyLength;
    
    const { tag, wireType } = decodeKey(keyValue);
    const field = tagToField.get(tag);
    
    if (!field) {
      // Unknown field, skip it
      offset = skipField(buffer, offset, wireType);
      continue;
    }
    
    let value: any;
    
    switch (field.type) {
      case FieldType.INT32:
      case FieldType.INT64:
      case FieldType.UINT32:
      case FieldType.UINT64: {
        const decoded = decodeVarint(buffer, offset);
        value = decoded.value;
        offset += decoded.length;
        break;
      }
      
      case FieldType.BOOL: {
        const decoded = decodeVarint(buffer, offset);
        value = decoded.value !== 0;
        offset += decoded.length;
        break;
      }
      
      case FieldType.FLOAT:
        value = decodeFloat(buffer, offset);
        offset += 4;
        break;
        
      case FieldType.DOUBLE:
        value = decodeDouble(buffer, offset);
        offset += 8;
        break;
        
      case FieldType.STRING: {
        const decoded = decodeString(buffer, offset);
        value = decoded.value;
        offset += decoded.length;
        break;
      }
      
      case FieldType.BYTES: {
        const decoded = decodeBytes(buffer, offset);
        value = decoded.value;
        offset += decoded.length;
        break;
      }
      
      default:
        throw new Error(`Unsupported field type: ${field.type}`);
    }
    
    result[field.name] = value;
  }
  
  return result;
}

/**
 * Skips an unknown field
 */
function skipField(buffer: Uint8Array, offset: number, wireType: WireType): number {
  switch (wireType) {
    case WireType.VARINT: {
      const { length } = decodeVarint(buffer, offset);
      return offset + length;
    }
    
    case WireType.FIXED64:
      return offset + 8;
      
    case WireType.LENGTH_DELIMITED: {
      const { value: length, length: headerLength } = decodeVarint(buffer, offset);
      return offset + headerLength + length;
    }
    
    case WireType.FIXED32:
      return offset + 4;
      
    default:
      throw new Error(`Unknown wire type: ${wireType}`);
  }
}
