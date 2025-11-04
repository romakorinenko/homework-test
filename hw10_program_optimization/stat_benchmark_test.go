package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

// BenchmarkGetDomainStat
// goos: darwin
// goarch: arm64
// pkg: github.com/romakorinenko/homework-test/hw10_program_optimization
// cpu: Apple M3 Pro
// BenchmarkGetDomainStat
// BenchmarkGetDomainStat-11    	 3134178	       374.4 ns/op
// PASS
// .
func BenchmarkGetDomainStat(b *testing.B) {
	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer func() {
		_ = r.Close()
	}()

	data, _ := r.File[0].Open()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
}
