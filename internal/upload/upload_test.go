package upload

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, err := w.Write([]byte("mock body"))
		assert.NoError(t, err)

		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Error("Header Content-Type does not contain multipart/form-data")
		}

		if r.RequestURI != "/upload" {
			t.Errorf("Expected uploaded to /upload, got %s", r.RequestURI)
		}
	}))

	defer ts.Close()

	data := []byte(strings.Repeat("na", 512))
	f, err := os.Create("testuploaddata.file")
	if err != nil {
		t.Errorf("Failed to create test datafile for uploading. Reason %s", err)
	}
	_, err = f.Write(data)
	assert.NoError(t, err)
	f.Close()
	defer os.Remove("testuploaddata.file")
	metadata := map[string]string{
		"task_url": "https://www.example.com/12345678",
	}
	body, err := Upload(ts.URL+"/upload", "testuploaddata.file", "", metadata)
	if err != nil {
		t.Error("ERROR from Upload:", err)
	}
	if string(body) != "mock body" {
		t.Error("Retrieved body is not expected")
	}
}
