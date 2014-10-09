package ratecounter

import (
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	c := NewCounter()

	check := func(expected int64) {
		val := c.Value()
		if val != expected {
			t.Error("Expected ", val, " to equal ", expected)
		}
	}

	check(0)
	c.Incr(1)
	check(1)
	c.Incr(9)
	check(10)

	// Concurrent usage
	wg := &sync.WaitGroup{}
	wg.Add(3)
	for i := 1; i <= 3; i++ {
		go func(val int64) {
			c.Incr(val)
			wg.Done()
		}(int64(i))
	}
	wg.Wait()
	check(16)
}

func TestRateCounter(t *testing.T) {
	interval := 500 * time.Millisecond
	r := NewRateCounter(interval)

	check := func(expected int64) {
		val := r.Rate()
		if val != expected {
			t.Error("Expected ", val, " to equal ", expected)
		}
	}

	check(0)
	r.Incr(1)
	check(1)
	r.Incr(2)
	check(3)
	time.Sleep(2 * interval)
	check(0)
}

func BenchmarkRateCounter(b *testing.B) {
	interval := 0 * time.Millisecond
	r := NewRateCounter(interval)

	for i := 0; i < b.N; i++ {
		r.Incr(1)
		r.Rate()
	}
}

func BenchmarkRateCounter_ScheduleDecrement(b *testing.B) {
	interval := 0 * time.Millisecond
	r := NewRateCounter(interval)

	for i := 0; i < b.N; i++ {
		r.scheduleDecrement(-1)
	}
}

func Benchmark_TimeNowAndAdd(b *testing.B) {
	a := []time.Time{}
	for i := 0; i < b.N; i++ {
		a = append(a, time.Now().Add(1*time.Second))
	}
	fmt.Fprintln(ioutil.Discard, a)
}
