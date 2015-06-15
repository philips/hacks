package restish

import (
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc/benchmark/stats"
)

func caller(hc testClient) {
	DoCall(hc, 1, 1)
}

func run(b *testing.B, maxConcurrentCalls int) {
	s := stats.AddStats(b, 38)
	b.StopTimer()
	target, stopper := StartServer("localhost:0")
	defer stopper()
	hc := testClient{target, &http.Client{}}

	// Warm up connection.
	for i := 0; i < 10; i++ {
		caller(hc)
	}
	ch := make(chan int, maxConcurrentCalls*4)
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	wg.Add(maxConcurrentCalls)

	// Distribute the b.N calls over maxConcurrentCalls workers.
	for i := 0; i < maxConcurrentCalls; i++ {
		go func() {
			for _ = range ch {
				start := time.Now()
				caller(hc)
				elapse := time.Since(start)
				mu.Lock()
				s.Add(elapse)
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
	b.StopTimer()
	close(ch)
	wg.Wait()
}

func BenchmarkClientc8(b *testing.B) {
	run(b, 8)
}

func BenchmarkClientc64(b *testing.B) {
	run(b, 64)
}

func BenchmarkClient512(b *testing.B) {
	run(b, 512)
}

func TestMain(m *testing.M) {
	os.Exit(stats.RunTestMain(m))
}
