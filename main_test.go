package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
	"zartekAssignment/variables"
	"zartekAssignment/visitors"
)

func TestVisitors(t *testing.T) {
	v := visitors.NewVisitors()

	// test serving a request with a new IP address
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:12345"
	rec := httptest.NewRecorder()

	v.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// test serving a request with the same IP address within the allowed limit
	for i := 0; i < variables.MaxRequests+1; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "1.2.3.4:12345"
		rec := httptest.NewRecorder()

		v.ServeHTTP(rec, req)

		resp := rec.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}
	}

	// test serving a request with the same IP address that exceeds the allowed limit
	req = httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:12345"
	rec = httptest.NewRecorder()

	v.ServeHTTP(rec, req)

	resp = rec.Result()
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected status code %d, but got %d", http.StatusTooManyRequests, resp.StatusCode)
	}

	// test serving a request after the allowed duration has elapsed
	time.Sleep(variables.Duration)

	req = httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:12345"
	rec = httptest.NewRecorder()

	v.ServeHTTP(rec, req)

	resp = rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// test serving a request with a new IP address while the server is running
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "5.6.7.8:54321"
		rec := httptest.NewRecorder()

		v.ServeHTTP(rec, req)

		resp := rec.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}
	}()

	// wait for the goroutine to finish before exiting the test
	wg.Wait()
}
