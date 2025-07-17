package http_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewRequestWithHeaders(method, url string, body []byte, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func DoRequest(t *testing.T, client *http.Client, method, url string, body any, headers map[string]string) ([]byte, int) {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	require.NoError(t, err)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Log("error while closing req body")
		}
	}()

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	fmt.Printf("%s %s Response: %s\n", method, url, string(data))

	return data, resp.StatusCode
}
