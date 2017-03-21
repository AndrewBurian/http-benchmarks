package main

import (
	"fmt"
	"github.com/andrewburian/http-benchmarks/routes"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sort"
	"time"
)

const (
	MIN_BENCH_TIME = 500
)

type TestResult struct {
	Name         string
	BaseMemSize  uint64
	GithubSize   uint64
	GPlusSize    uint64
	GithubRoutes TimeResult
	GPlusRoutes  TimeResult
	SingleRoute  TimeResult
}

type TimeResult struct {
	min, max, median, average, p95, p99 time.Duration
}

func RunTests(routers []Router) {

	results := make([]*TestResult, 0, len(routers))
	for _, router := range routers {
		result := TestRouter(router)
		if result == nil {
			fmt.Println(router.Name, "panicked")
			continue
		}
		results = append(results, result)
		fmt.Println(result.Name)
		fmt.Println("Base memory:", result.BaseMemSize, "B")
		fmt.Println("Github memory:", result.GithubSize, "B")
		fmt.Println("Google+ memory:", result.GPlusSize, "B")
		fmt.Println("Single route time:", result.SingleRoute)
	}

	for _, _ = range results {
	}
}

func TestRouter(router Router) *TestResult {

	// safe from panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var result TestResult

	result.Name = router.Name
	result.BaseMemSize = TestBaseSize(router)
	result.GithubSize = TestGithubSize(router)
	result.GPlusSize = TestGPlusSize(router)
	result.SingleRoute = BenchSingleRoute(router)

	return &result
}

func TestBaseSize(router Router) uint64 {

	var m runtime.MemStats

	// before
	runtime.GC()
	runtime.ReadMemStats(&m)
	before := m.HeapAlloc

	r := router.SetupRoutes([]routes.Route{}, "")

	// after
	runtime.GC()
	runtime.ReadMemStats(&m)
	after := m.HeapAlloc

	// do something with the router so we are allowed to declare it
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	return after - before
}

func TestGithubSize(router Router) uint64 {

	var m runtime.MemStats

	// before
	runtime.GC()
	runtime.ReadMemStats(&m)
	before := m.HeapAlloc

	r := router.SetupRoutes(routes.Github, "")

	// after
	runtime.GC()
	runtime.ReadMemStats(&m)
	after := m.HeapAlloc

	// do something with the router so we are allowed to declare it
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	return after - before
}

func TestGPlusSize(router Router) uint64 {

	var m runtime.MemStats

	// before
	runtime.GC()
	runtime.ReadMemStats(&m)
	before := m.HeapAlloc

	r := router.SetupRoutes(routes.GPlus, "")

	// after
	runtime.GC()
	runtime.ReadMemStats(&m)
	after := m.HeapAlloc

	// do something with the router so we are allowed to declare it
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	return after - before
}

func BenchSingleRoute(router Router) TimeResult {

	r := router.SetupRoutes([]routes.Route{{http.MethodGet, "/"}}, "")
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// todo all logic below here is garbage. Don't code tired kids
	var result TimeResult
	var iterations uint64 = 5

	for result.min < MIN_BENCH_TIME {

		times := make([]time.Duration, 0, iterations)
		var runningCount time.Duration

		for i := uint64(0); i < iterations; i++ {
			start := time.Now()
			r.ServeHTTP(res, req)
			end := time.Now()

			diff := end.Sub(start)
			times = append(times, diff)

			runningCount += diff
		}

		sort.Slice(times, func(i, j int) bool {
			return times[i].Nanoseconds() < times[j].Nanoseconds()
		})

		result.min = times[0]
		result.max = times[iterations-1]
		result.median = times[iterations/2]
		result.average = time.Duration(uint64(runningCount) / iterations)
		result.p95 = times[iterations*95/100]
		result.p99 = times[iterations*99/100]

		iterations *= 2
	}

	return result
}
