package downloader

import (
	"testing"

	"github.com/tickstep/aliyunpan-api/aliyunpan"
)

func TestSelectFileDownloadUrlByExtensionUsesLivePhotoStreams(t *testing.T) {
	durl := &aliyunpan.GetFileDownloadUrlResult{
		Url: "https://example.com/bad.livp",
		StreamsUrl: &aliyunpan.ShareAlbumFileStreamUrlItem{
			Heic: "https://example.com/live.heic",
			Jpeg: "https://example.com/live.jpg",
			Mov:  "https://example.com/live.mov",
		},
	}

	tests := map[string]string{
		"heic":  "https://example.com/live.heic",
		".jpg":  "https://example.com/live.jpg",
		"jpeg":  "https://example.com/live.jpg",
		".mov":  "https://example.com/live.mov",
		"other": "https://example.com/bad.livp",
	}

	for ext, want := range tests {
		if got := selectFileDownloadUrlByExtension(durl, ext); got != want {
			t.Fatalf("ext %q: got %q, want %q", ext, got, want)
		}
	}
}

func TestSelectFileDownloadUrlByExtensionFallsBackToUrl(t *testing.T) {
	durl := &aliyunpan.GetFileDownloadUrlResult{
		Url: "https://example.com/file.bin",
	}

	if got := selectFileDownloadUrlByExtension(durl, "bin"); got != durl.Url {
		t.Fatalf("got %q, want %q", got, durl.Url)
	}
}

func TestSelectFileDownloadUrlByExtensionFallsBackBetweenStillPhotoStreams(t *testing.T) {
	durl := &aliyunpan.GetFileDownloadUrlResult{
		StreamsUrl: &aliyunpan.ShareAlbumFileStreamUrlItem{
			Jpeg: "https://example.com/still-from-jpeg-field",
		},
	}

	if got := selectFileDownloadUrlByExtension(durl, "heic"); got != durl.StreamsUrl.Jpeg {
		t.Fatalf("got %q, want %q", got, durl.StreamsUrl.Jpeg)
	}
}

func TestSetUseWebApi(t *testing.T) {
	der := NewDownloader(nil, nil, nil, nil, nil)
	der.SetUseWebApi(true)

	if !der.useWebApi {
		t.Fatal("expected downloader to force WebAPI download URL lookup")
	}
}
