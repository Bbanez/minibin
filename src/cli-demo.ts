#!/usr/bin/env node

/**
 * Simple CLI demo for minibin
 * This demonstrates practical usage of the library
 */

const { encode, decode, FieldType } = require('./index.js');

// Define a Person schema
const personSchema = {
  id: { type: FieldType.INT32, tag: 1 },
  name: { type: FieldType.STRING, tag: 2 },
  email: { type: FieldType.STRING, tag: 3 },
  age: { type: FieldType.INT32, tag: 4 },
  active: { type: FieldType.BOOL, tag: 5 },
};

// Example person object
const person = {
  id: 1001,
  name: 'John Doe',
  email: 'john.doe@example.com',
  age: 28,
  active: true,
};

console.log('╔════════════════════════════════════════════════════════╗');
console.log('║           Minibin - Binary Encoding Demo              ║');
console.log('╚════════════════════════════════════════════════════════╝');
console.log();

console.log('📝 Original Object:');
console.log(JSON.stringify(person, null, 2));
console.log();

// Encode the person object
const encoded = encode(personSchema, person);

console.log('🔐 Encoded Binary Data:');
console.log('   Length:', encoded.length, 'bytes');
console.log('   Hex:', Array.from(encoded).map(b => (b as number).toString(16).padStart(2, '0')).join(' '));
console.log('   Raw:', Array.from(encoded).join(', '));
console.log();

// Calculate compression ratio
const jsonSize = JSON.stringify(person).length;
const compressionRatio = ((1 - encoded.length / jsonSize) * 100).toFixed(1);
console.log('📊 Size Comparison:');
console.log('   JSON:', jsonSize, 'bytes');
console.log('   Binary:', encoded.length, 'bytes');
console.log('   Compression:', compressionRatio + '%');
console.log();

// Decode the binary data
const decoded = decode(personSchema, encoded);

console.log('🔓 Decoded Object:');
console.log(JSON.stringify(decoded, null, 2));
console.log();

// Verify data integrity
const matches = JSON.stringify(person) === JSON.stringify(decoded);
console.log('✅ Data Integrity:', matches ? 'PASSED' : 'FAILED');
console.log();

console.log('╔════════════════════════════════════════════════════════╗');
console.log('║                    Demo Complete                       ║');
console.log('╚════════════════════════════════════════════════════════╝');
