package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type StressTest struct {
	Url              string
	Requests         int
	Concurrency      int
	MetricsCollected MetricsCollected
}

type MetricsCollected struct {
	TotalRequests        int64
	TotalTimeNanoSeconds int64
	TotalRequest2xx      int64
	TotalRequest4xx      int64
	TotalRequest5xx      int64
	TotalWithErrors      int64
}

func NewStressTest(url string, requests int, concurrency int) *StressTest {
	return &StressTest{
		Url:              url,
		Requests:         requests,
		Concurrency:      concurrency,
		MetricsCollected: MetricsCollected{},
	}
}

func (st *StressTest) ExecuteTests() {
	fmt.Println("\n Running tests...")
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(st.Concurrency)
	go func() {
		for i := 0; i < st.Concurrency; i++ {
			go func() {
				defer waitGroup.Done()
				for j := 0; j < st.Requests; j++ {
					start := time.Now()
					atomic.AddInt64(&st.MetricsCollected.TotalRequests, 1)

					resp, err := http.Get(st.Url)
					if err != nil {
						atomic.AddInt64(&st.MetricsCollected.TotalWithErrors, 1)
						return
					}

					strStatusCode := strconv.Itoa(resp.StatusCode)
					switch {
					case strings.HasPrefix(strStatusCode, "2"):
						atomic.AddInt64(&st.MetricsCollected.TotalRequest2xx, 1)
					case strings.HasPrefix(strStatusCode, "4"):
						atomic.AddInt64(&st.MetricsCollected.TotalRequest4xx, 1)
					case strings.HasPrefix(strStatusCode, "5"):
						atomic.AddInt64(&st.MetricsCollected.TotalRequest5xx, 1)
					}
					elapsed := time.Since(start).Nanoseconds()
					atomic.AddInt64(&st.MetricsCollected.TotalTimeNanoSeconds, elapsed)
				}
			}()
		}
	}()
	waitGroup.Wait()
}

func (st *StressTest) PrintMetrics() {
	fmt.Println("-----")
	fmt.Printf("URL Tested %s \n", st.Url)
	fmt.Printf("Number of Requests %d \n", st.Requests)
	fmt.Printf("Using Concurrency of %d \n", st.Concurrency)
	fmt.Println("-----")
	fmt.Println("RESULTS:")
	fmt.Printf("TotalRequests %d \n", st.MetricsCollected.TotalRequests)
	fmt.Printf("TotalTimeNanoSeconds %d \n", st.MetricsCollected.TotalTimeNanoSeconds)
	fmt.Printf("TotalRequest2xx %d \n", st.MetricsCollected.TotalRequest2xx)
	fmt.Printf("TotalRequest4xx %d \n", st.MetricsCollected.TotalRequest4xx)
	fmt.Printf("TotalRequest5xx %d \n", st.MetricsCollected.TotalRequest5xx)
	fmt.Printf("TotalWithErrors %d \n", st.MetricsCollected.TotalWithErrors)
}
