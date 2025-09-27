package gribtest

import (
	"os"
	"testing"

	"github.com/cgray1234/grib/griblib"
)

func BenchmarkReadMessages(b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, err := os.Open("../integrationtestdata/template5_3.grib2")
		if err != nil {
			b.Fatalf("Could not open test-file %v", err)
		}
		griblib.ReadMessages(f)
	}
}
