package main
/*
import (
	"net/http"
	"path/filepath"
	"runtime"
	"testing"
)

type mockResponseWriter struct {
	header      http.Header
	captureData []byte
}

func (mrw mockResponseWriter) Header() http.Header {
	return mrw.header
}

func (mrw mockResponseWriter) Write(data []byte) (int, error) {
	mrw.captureData = data
	return len(data), nil
}

func (mrw mockResponseWriter) WriteHeader(code int) {
	//no op
}

func Test(t *testing.T) {



	rr := NewRestRouter()

	rr.LoadLibrary()

	rr.LoadLibrary(mrw, nil)


}
*/