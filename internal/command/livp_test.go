package command

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseContentRangeTotal(t *testing.T) {
	tests := map[string]int64{
		"bytes 0-0/613963":       613963,
		"bytes 0-1023/4273201":   4273201,
		"bytes */4273201":        4273201,
		"":                       -1,
		"bytes 0-0/*":            -1,
		"invalid content range":  -1,
		"bytes 0-0/not-a-number": -1,
	}

	for contentRange, want := range tests {
		if got := parseContentRangeTotal(contentRange); got != want {
			t.Fatalf("content range %q: got %d, want %d", contentRange, got, want)
		}
	}
}

func TestInferHTTPDownloadFileExtension(t *testing.T) {
	tests := map[string]string{
		"image/heif":                   "heic",
		"image/heic; charset=binary":   "heic",
		"image/jpeg":                   "jpg",
		"video/quicktime":              "mov",
		"application/octet-stream":     "",
		"application/x-unknown-format": "",
	}

	for contentType, want := range tests {
		resp := &http.Response{Header: http.Header{"Content-Type": []string{contentType}}}
		if got := inferHTTPDownloadFileExtension(resp); got != want {
			t.Fatalf("content type %q: got %q, want %q", contentType, got, want)
		}
	}
}

func TestInferDownloadFileExtensionFromBytes(t *testing.T) {
	tests := map[string]struct {
		data []byte
		want string
	}{
		"jpeg": {data: []byte{0xff, 0xd8, 0xff, 0xe0}, want: "jpg"},
		"heic": {data: []byte{
			0x00, 0x00, 0x00, 0x18,
			'f', 't', 'y', 'p',
			'h', 'e', 'i', 'c',
		}, want: "heic"},
		"quicktime": {data: []byte{
			0x00, 0x00, 0x00, 0x14,
			'f', 't', 'y', 'p',
			'q', 't', ' ', ' ',
		}, want: "mov"},
		"unknown": {data: []byte("not enough"), want: ""},
	}

	for name, tt := range tests {
		if got := inferDownloadFileExtensionFromBytes(tt.data); got != tt.want {
			t.Fatalf("%s: got %q, want %q", name, got, tt.want)
		}
	}
}

func TestGetHTTPDownloadFileInfoIgnoresErrorContentRange(t *testing.T) {
	const movSize = int64(3128969)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Range", "bytes */298")
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			_, _ = w.Write([]byte("range error"))
		case http.MethodHead:
			w.Header().Set("Content-Type", "video/quicktime")
			w.Header().Set("Content-Length", "3128969")
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	info := getHttpDownloadFileInfo(server.URL)
	if info.FileSize != movSize {
		t.Fatalf("got size %d, want %d", info.FileSize, movSize)
	}
	if info.FileExtension != "mov" {
		t.Fatalf("got extension %q, want mov", info.FileExtension)
	}
}
