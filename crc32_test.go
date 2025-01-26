package crc32combine

import (
	"hash/crc32"
	"io"
	mathrand "math/rand"
	"testing"
)

func TestCombine(t *testing.T) {
	poly := map[string]uint32{"IEEE": crc32.IEEE, "Castagnoli": crc32.Castagnoli, "Koopman": crc32.Koopman}

	for name, poly := range poly {
		t.Run(name, func(t *testing.T) {
			table := crc32.MakeTable(poly)
			rand := mathrand.New(mathrand.NewSource(0))
			data1, _ := io.ReadAll(io.LimitReader(rand, 1<<20))
			sum1 := crc32.Checksum(data1, table)

			data2, _ := io.ReadAll(io.LimitReader(rand, 2<<20))
			sum2 := crc32.Checksum(data2, table)

			expectedDigest := crc32.Checksum(append(data1, data2...), table)
			digest := Combine(table, sum1, sum2, len(data2))

			if digest != expectedDigest {
				t.Errorf("expected %x, got %x", expectedDigest, digest)
			}
		})
	}
}

func BenchmarkCombine(b *testing.B) {
	poly := map[string]uint32{"IEEE": crc32.IEEE, "Castagnoli": crc32.Castagnoli, "Koopman": crc32.Koopman}

	for name, poly := range poly {
		b.Run(name, func(b *testing.B) {
			table := crc32.MakeTable(poly)
			rand := mathrand.New(mathrand.NewSource(0))

			sum1 := crc32.New(table)
			_, _ = io.Copy(sum1, io.LimitReader(rand, 1<<20))
			sum2 := crc32.New(table)
			_, _ = io.Copy(sum2, io.LimitReader(rand, 1<<20))

			Combine(table, sum1.Sum32(), sum2.Sum32(), 1<<20)
		})
	}
}
