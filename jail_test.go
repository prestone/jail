package jail

//testing

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var r = New("test.db")

func TestNew(t *testing.T) {
	a := assert.New(t)
	_ = a

	r.Put(1, time.Second*2, "ban")
	r.Put(2, time.Second*5, "ban")
	r.Put(3, time.Second*10000, "ban")
	r.Put(4, time.Second*10000, "ban")
	a.Equal(4, r.count)

	a.True(r.Check(1))
	a.True(r.Check(2))
	a.True(r.Check(3))

	time.Sleep(time.Second * 2)

	a.False(r.Check(1))
	a.True(r.Check(2))
	a.True(r.Check(3))

	r.Delete(3)
	a.False(r.Check(3))

	a.Equal("ban", r.Reason(2))

}

func BenchmarkPut(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.Put(rand.Int(), time.Second, "ban")
	}
}

func BenchmarkCheck(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.Check(rand.Int())
	}
}

func BenchmarkGetFreeParallel(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Put(rand.Int(), time.Second, "ban")
		}
	})
}

func TestD(t *testing.T) {
	os.Remove("test.db")
}
