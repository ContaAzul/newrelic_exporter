package newrelic

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

// responseStub fakes a NewRelic API response
func responseStub(t *testing.T, filename string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bts := readTestData(t, filename)
		w.WriteHeader(status)
		w.Write(bts)
	}))
}

// readTestData reads the file named by filename under the testdata directory and returns the contents
func readTestData(t *testing.T, filename string) []byte {
	path := filepath.Join("testdata", filename)
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bts
}
