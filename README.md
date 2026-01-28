# MiniBin

Binary data format for transferring data between clients and services or storing it locally. It was inspired by ProtoBuf and the reason why developed this data format is because ProtoBuf is a big library with a lot of features. I usually need data encoding and decoding with type checking while transferring data and I need it to be able to generate types, encoders and decoders for multiple languages.

## Supported languages

- [x] - Go
- [x] - TypeScript
- [ ] - Rust (coming soon)
- [ ] - C/C++

## How it works

After downloading the program you can use it to convert JSON schema to language specific object that is used to encode and decode object data. In addition to that, you can use generated objects in your code to access data directly.

- `minibin -l go` - Running this command will convert all JSON schema files in `minibin-schemas` to Go objects and place them in `src/minibin`.
- `minibin -l ts` - Will do the same as above command just for TypeScript.

### Encoding/Decoding

| Prop pos | Value | Type   | Data len bites | Data bits |
| -------- | ----- | ------ | -------------- | --------- |
| 8        | 0     | null   | 0              | 0         |
| 8        | 1     | string | 8-16-24-32     | N         |
| 8        | 2     | i32    | 0              | 32        |
| 8        | 3     | i64    | 0              | 64        |
| 8        | 4     | u32    | 0              | 32        |
| 8        | 5     | u64    | 0              | 64        |
| 8        | 6     | f32    | 0              | 32        |
| 8        | 7     | f64    | 0              | 64        |
| 8        | 8     | bool   | 0              | 8         |
| 8        | 9     | object | 32             | N         |
| 8        | 10    | enum   | 8              | N         |
| 8        | 11    | bytes  | 8              | N         |

```txt
 0xFF   0xFF   0xFF-0xFFFFFFFF   0xN
|----| |----| |---------------| |---|
|      |      |                 \ Data bytes
|      |      \ 1-4 bytes for data length (optional, might not exist)
|      \ combination of property type and data length
\ property position
```
