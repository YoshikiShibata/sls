package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	lockURL   = "http://localhost:8000/lock"
	unlockURL = "http://localhost:8000/unlock"
	clearURL  = "http://localhost:8000/clear"
)

// slsサーバーを別途起動しておくこと
func TestSimpleLockAndUnlock(t *testing.T) {
	paths := []string{"login", "shop", "operator"}

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			uuid := uuid.NewString()
			lockReq := createLockReq(t, uuid, path)
			unlockReq := createUnlockReq(t, uuid, path)

			resp := sendRequest(t, lockURL, lockReq)
			readBodyAndShowResponse(t, resp)

			time.Sleep(400 * time.Millisecond)

			resp = sendRequest(t, unlockURL, unlockReq)
			readBodyAndShowResponse(t, resp)
		}(paths[i%3])
	}
	wg.Wait()
}

func TestSimpleClearAll(t *testing.T) {
	paths := []string{"login", "shop", "operator"}

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			uuid := uuid.NewString()
			lockReq := createLockReq(t, uuid, path)
			resp := sendRequest(t, lockURL, lockReq)

			readBodyAndShowResponse(t, resp)
		}(paths[i%3])
	}

	time.Sleep(3 * time.Second)
	clearReq := createClearReq(t, true)
	resp := sendRequest(t, clearURL, clearReq)
	readBodyAndShowResponse(t, resp)

	wg.Wait()
}

func readBodyAndShowResponse(
	t *testing.T,
	resp *http.Response,
) {
	// レスポンスのBodyを読み込む
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll failed: %v", err)
	}

	// レスポンスの内容を表示
	t.Logf("Response Status: %v", resp.Status)
	t.Logf("Response Body: %s", string(body))
}

func createLockReq(
	t *testing.T,
	uuid, path string,
) io.Reader {
	t.Helper()

	req := &LockRequest{UUID: uuid, Path: path}
	buf, err := json.Marshal(&req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	return bytes.NewBuffer(buf)
}

func createUnlockReq(
	t *testing.T,
	uuid, path string,
) io.Reader {
	t.Helper()

	req := &UnlockRequest{UUID: uuid, Path: path}
	buf, err := json.Marshal(&req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	return bytes.NewBuffer(buf)
}

func createClearReq(
	t *testing.T,
	clearAll bool,
) io.Reader {
	t.Helper()

	req := &ClearRequest{ClearAll: clearAll}
	buf, err := json.Marshal(&req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	return bytes.NewBuffer(buf)
}

func sendRequest(
	t *testing.T,
	url string,
	body io.Reader,
) *http.Response {
	// HTTPリクエストを作成
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		t.Fatalf("http.Post failed: %v", err)
	}
	t.Cleanup(func() {
		resp.Body.Close()
	})

	return resp
}
