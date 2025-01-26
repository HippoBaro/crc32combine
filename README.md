# CRC32 Combine

A lightweight, zero-dependency module that exports a single function: `Combine`.

`Combine` merges two CRC32 checksums as if their corresponding data streams had been processed sequentially. Given `A`, the checksum of an initial byte stream, and `B`, the checksum of a second byte stream of size `length`, this function computes `AB`, the CRC32 that would have been obtained if both streams had been concatenated and processed as a single, continuous input.

For a deeper understanding of why this works, see the definitive explanation by Mark Adler (yes, as in `adler32`): https://stackoverflow.com/a/23126768

> Adapted from https://github.com/madler/zlib; modified so it plays well with Go's `crc32` package.

## Benchmarks

```
goos: darwin
goarch: arm64
pkg: github.com/HippoBaro/crc32combine
cpu: Apple M1 Max
BenchmarkCombine
BenchmarkCombine/IEEE
BenchmarkCombine/IEEE-10         	1000000000	         0.002435 ns/op
BenchmarkCombine/Castagnoli
BenchmarkCombine/Castagnoli-10   	1000000000	         0.002134 ns/op
BenchmarkCombine/Koopman
BenchmarkCombine/Koopman-10      	1000000000	         0.006929 ns/op
PASS
```