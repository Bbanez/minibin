# MiniBin

Binary data format for transferring data between clients and services. It
was inspired by ProtoBuf and the reason why developed this data format is
because ProtoBuf is a big library with a lot of features. I usually need
data encoding and decoding with type checking while transferring data and
I need it to be able to generate types, encoders and decoders for multiple
languages.

## How it works

### Binary bytes

| Value | Type   | Pos | Data len bites | Data bits |
| ----- | ------ | --- | -------------- | --------- |
| 0     | null   | 8   | 0              | 0         |
| 1     | string | 8   | 32             | N         |
| 2     | i32    | 8   | 0              | 32        |
| 3     | i64    | 8   | 0              | 64        |
| 4     | u32    | 8   | 0              | 32        |
| 5     | u64    | 8   | 0              | 64        |
| 6     | f32    | 8   | 0              | 32        |
| 7     | f64    | 8   | 0              | 64        |
| 8     | bool   | 8   | 0              | 8         |
| 9     | object | 8   | 32             | N         |
| 10    | enum   | 8   | 8              | N         |
| 11    | bytes  | 8   | 8              | N         |

### String encoding

  1 byte  000000  00  1-4 bytes  N bytes
  |----|  |----| |--| |-------|  |-----| 
  |       |      |    |          | String data bytes
  |       |      |    | Data bytes len
  |       |      | Len of data bytes len
  |       | Data type
  | Property position


Lets say that we have a data type like this:

```json
{
  "name": "User",
  "props": [
    {
      "name": "id",
      "pos": 0,
      "typ": "string",
      "required": true
    },
    {
      "name": "createdAt",
      "pos": 1,
      "typ": "u64",
      "required": true
    },
    {
      "name": "updatedAt",
      "pos": 2,
      "typ": "u64",
      "required": true
    },
    {
      "name": "name",
      "pos": 3,
      "typ": "string"
    },
    {
      "name": "email",
      "pos": 4,
      "typ": "string",
      "required": true
    }
  ]
}
```

```json
{
  "id": "1",
  "createdAt": 1,
  "updatedAt": 2,
  "name": "John",
  "email": "test@test.com"
}
```

Encoded message will look like this:

```txt
00000001 | 0000 0000 0000 0000 0000 0000 0000 0001 | 0000 0000 0000 0000 0011 0001 >
00000101 | 0000 0000 0000 0000 0000 0000 0000 0001 >
00000101 | 0000 0000 0000 0000 0000 0000 0000 0010 >
00000001 | 0000 0000 0000 0000 0000 0000 0000 0100 | 0000 0000 0000 0000 0100 1010 _ 0000 0000 0000 0000 0110 1111 _ 0000 0000 0000 0000 0110 1000 _ 0000 0000 0000 0000 0110 1110
.....
```
