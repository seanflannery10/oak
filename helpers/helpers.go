package helpers

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func GetBody(t *testing.T, rs *http.Response) string {
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return string(body)
}
