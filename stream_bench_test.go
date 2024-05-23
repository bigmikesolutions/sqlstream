package sqlstream_test

import (
	"fmt"
	"sqlstream"
	"testing"
)

func BenchmarkStream_ReadRows(b *testing.B) {
	stream := newBenchStream()

	for i := 0; i < b.N; i++ {
		<-stream
	}

}

func BenchmarkStream_ReadRows_Debug(b *testing.B) {
	b.Skipf("used only for checking memory allocations per iteration")

	b.StopTimer()
	testLoops := 5
	stream := newBenchStream()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		entry := <-stream
		b.StopTimer()

		fmt.Printf("%d: %+v, memory: %p, %p\n",
			i, entry.Value,
			&entry, &entry.Value,
		)

		if i == testLoops {
			break
		}
	}

}

func newBenchStream() sqlstream.ReadStream[testObject] {
	id := 0
	obj := &testObject{
		ID:     id,
		Field1: "",
	}
	gen := newTRowsGenerator[testObject](func() (*testObject, error) {
		id += 1
		obj.ID = id
		return obj, nil
	})

	return sqlstream.Read[testObject](gen)
}
