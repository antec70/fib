package internal

import (
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var s Server

func TestFib(t *testing.T) {
	testsTable := []struct {
		name          string
		input         uint32
		expectedValue uint32
	}{
		{
			name:          "Ok",
			input:         10,
			expectedValue: 55,
		},
		{
			name:          "Ok",
			input:         13,
			expectedValue: 233,
		},
	}
	for _, testCase := range testsTable {
		t.Logf("Running test case %s", testCase.name)
		t.Run(testCase.name, func(t *testing.T) {
			res := Fib(testCase.input)
			assert.Equal(t, testCase.expectedValue, res)
		})
	}

}
func BenchmarkSlowFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SlowFib(50)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(50)
	}
}

func TestServer_CalcFib(t *testing.T) {
	db, _ := redismock.NewClientMock()
	s.rdb = db
	testsTable := []struct {
		name               string
		x                  uint32
		y                  uint32
		expectedCountItems int
	}{
		{
			name:               "ok",
			x:                  0,
			y:                  16,
			expectedCountItems: 16,
		},
	}
	for _, testCase := range testsTable {
		t.Logf("Running test case %s", testCase.name)
		t.Run(testCase.name, func(t *testing.T) {
			res := s.CalcFib(testCase.x, testCase.y)
			assert.Equal(t, testCase.expectedCountItems, len(res))
		})
	}
}

func TestHandler_handleIndex(t *testing.T) {
	t.Parallel()
	db, _ := redismock.NewClientMock()
	s.rdb = db
	testsTable := []struct {
		name         string
		url          string
		x            string
		y            string
		expectedCode int
	}{
		{
			name:         "Ok",
			url:          "http://localhost:8080/fib",
			x:            "0",
			y:            "16",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Bad request for X value",
			url:          "http://localhost:8080/fib",
			x:            "-1",
			y:            "16",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Empty y value",
			url:          "http://localhost:8080/fib",
			x:            "1",
			y:            "",
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, testCase := range testsTable {
		t.Logf("Running test case %s", testCase.name)
		t.Run(testCase.name, func(t *testing.T) {

			req, err := http.NewRequest("GET", testCase.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			q := req.URL.Query()
			q.Add("x", testCase.x)
			q.Add("y", testCase.y)
			req.URL.RawQuery = q.Encode()
			response := executeRequest(req)
			assert.Equal(t, testCase.expectedCode, response.Code)
		})
	}
}
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	http.HandlerFunc(s.handleIndex()).ServeHTTP(rr, req)
	return rr
}
