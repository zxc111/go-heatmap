package heatmap

import (
	"image"
	"io"
	"math"
	"math/rand"
	"testing"

	"github.com/zxc111/go-heatmap/schemes"
)

var testPoints = []DataPoint{}

func init() {
	rand.Seed(4828458249)
	for n := 0; n < 350; n++ {
		testPoints = append(testPoints,
			P(rand.Float64(), rand.Float64()))
	}
}

func assertEpsilon(t *testing.T, field string, expected, got float64) {
	if math.Abs(got-expected) > 0.00001 {
		t.Fatalf("Expected %v for %v, got %v",
			expected, field, got)
	}
}
func TestFindLimits(t *testing.T) {
	l := findLimits(testPoints)

	assertEpsilon(t, "minx", 0.005923822722460793, l.Min.X())
	assertEpsilon(t, "miny", 0.0038807964380815894, l.Min.Y())
	assertEpsilon(t, "maxx", 0.984841695550329, l.Max.X())
	assertEpsilon(t, "maxy", 0.9990553050102642, l.Max.Y())

	assertEpsilon(t, "dx", 0.9789178728278681, l.Dx())
	assertEpsilon(t, "dy", 0.9951745085721826, l.Dy())
}

func TestHeatmapKMLLimits(t *testing.T) {
	tests := []DataPoint{
		P(-200, 0),
		P(200, 0),
		P(0, 100),
		P(0, -100),
	}

	for _, test := range tests {
		lim := findLimits([]DataPoint{test})
		if lim.inRange(-180, 180, -90, 90) {
			t.Errorf("Expected out of range on %v", test)
		}
	}
}

const expHash = uint64(62624876249118208)

func TestMust(t *testing.T) {
	must(nil) // no panic
	panicked := false
	func() {
		defer func() { panicked = recover() != nil }()
		must(io.EOF)
	}()
	if !panicked {
		t.Fatalf("Expected a panic, but didn't get one")
	}
}

func BenchmarkPlacement(b *testing.B) {
	b.StopTimer()
	l := findLimits(testPoints)
	size := image.Rect(0, 0, 4096, 4096)
	dot := mkDot(float64(100))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		bw := image.NewRGBA(size)
		b.StartTimer()

		placePoints(size, l, bw, testPoints, dot)
	}
}

func BenchmarkWarming(b *testing.B) {
	b.StopTimer()

	l := findLimits(testPoints)
	size := image.Rect(0, 0, 4096, 4096)
	dot := mkDot(float64(100))
	colors := schemes.AlphaFire

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		bw := image.NewRGBA(size)
		placePoints(size, l, bw, testPoints, dot)
		out := image.NewRGBA(size)
		b.StartTimer()

		warm(out, bw, 64, colors)
	}
}

func BenchmarkWarmingParallel(b *testing.B) {
	l := findLimits(testPoints)
	size := image.Rect(0, 0, 4096, 4096)
	dot := mkDot(float64(100))
	colors := schemes.AlphaFire

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		bw := image.NewRGBA(size)
		placePoints(size, l, bw, testPoints, dot)
		out := image.NewRGBA(size)

		for pb.Next() {
			warm(out, bw, 64, colors)
		}
	})
}
