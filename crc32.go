// Package crc32combine provides a lightweight, zero-dependency implementation
// for merging two CRC32 checksums without needing to recompute the CRC from
// scratch.
//
// // Ported from zlib (https://github.com/madler/zlib) and modified so it
// plays well with Go's `crc32` package.
package crc32combine

import "hash/crc32"

// Precomputed matrices for GF(2) arithmetic
var (
	initOdd = gf2Mat{
		0, 0, 0, 0, 1 << 0, 1 << 1, 1 << 2, 1 << 3,
		1 << 4, 1 << 5, 1 << 6, 1 << 7, 1 << 8, 1 << 9, 1 << 10, 1 << 11,
		1 << 12, 1 << 13, 1 << 14, 1 << 15, 1 << 16, 1 << 17, 1 << 18, 1 << 19,
		1 << 20, 1 << 21, 1 << 22, 1 << 23, 1 << 24, 1 << 25, 1 << 26, 1 << 27,
	}
	initEven = gf2Mat{
		0, 0, 1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5,
		1 << 6, 1 << 7, 1 << 8, 1 << 9, 1 << 10, 1 << 11, 1 << 12, 1 << 13,
		1 << 14, 1 << 15, 1 << 16, 1 << 17, 1 << 18, 1 << 19, 1 << 20, 1 << 21,
		1 << 22, 1 << 23, 1 << 24, 1 << 25, 1 << 26, 1 << 27, 1 << 28, 1 << 29,
	}
)

type gf2Mat [32]uint32

// gf2Mul performs matrix-vector multiplication in GF(2)
func (m *gf2Mat) mul(vec uint32) (sum uint32) {
	for i := 0; vec != 0; i, vec = i+4, vec>>4 {
		sum ^= (m[i] * (vec & 1)) ^
			(m[i+1] * ((vec >> 1) & 1)) ^
			(m[i+2] * ((vec >> 2) & 1)) ^
			(m[i+3] * ((vec >> 3) & 1))
	}
	return sum
}

// gf2Sqrt computes the square of each element in the given GF(2) matrix
func (m *gf2Mat) sqrt(lhs *gf2Mat) {
	for n, val := range lhs {
		m[n] = lhs.mul(val)
	}
}

// Combine merges two CRC32 checksums as if their corresponding data streams
// had been processed sequentially. Given `crc1`, the checksum of an initial
// byte stream, and `crc2`, the checksum of a second byte stream of size `length`,
// this function computes the CRC that would have been obtained if both streams
// had been concatenated and processed as a single, continuous input.
//
// For a deeper understanding of why this works, see the definitive explanation by
// Mark Adler (as in adler32, yes): https://stackoverflow.com/a/23126768
func Combine(poly *crc32.Table, crc1, crc2 uint32, length int) uint32 {
	var (
		even = initEven
		odd  = initOdd
	)

	// Guard against invalid values
	if length <= 0 {
		return crc1
	}

	// Set polynomial seeds from table
	odd[0], odd[1], odd[2], odd[3] = poly[1<<4], poly[1<<5], poly[1<<6], poly[1<<7]
	even[0], even[1] = poly[1<<6], poly[1<<7]

	// Adjusts `crc1` to account for `length` bytes of zero-padding. The resulting
	// CRC value is equivalent to computing the CRC over `length` zero bytes
	// followed by the original data that produced `crc1`.
	for pOdd, pEven := &odd, &even; length > 0; pOdd, pEven = pEven, pOdd {
		if pEven.sqrt(pOdd); length&1 != 0 {
			crc1 = pEven.mul(crc1)
		}

		// Process next bit
		length >>= 1
	}

	// Finally, XOR the two CRCs to combine them.
	return crc1 ^ crc2
}
