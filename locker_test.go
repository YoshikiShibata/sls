package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

// slsサーバーを別途起動しておくこと
func TestSimpleLockAndUnlock(t *testing.T) {
	const (
		lockURL   = "http://localhost:8000/lock"
		unlockURL = "http://localhost:8000/unlock"
	)

	readBodyAndShowResponse := func(resp *http.Response) {
		// レスポンスのBodyを読み込む
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("ioutil.ReadAll failed: %v", err)
		}

		// レスポンスの内容を表示
		t.Logf("Response Status: %v", resp.Status)
		t.Logf("Response Body: %s", string(body))
	}

	paths := []string{"login", "shop", "operator"}

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			lockReq := &LockRequest{
				UUID: uuid.NewString(),
				Path: path,
			}
			unlockReq := &UnlockRequest{
				UUID: lockReq.UUID,
				Path: lockReq.Path,
			}

			lockReqBuf, err := json.Marshal(&lockReq)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}
			unlockReqBuf, err := json.Marshal(&unlockReq)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}

			// HTTPリクエストを作成
			resp, err := http.Post(lockURL, "application/json", bytes.NewBuffer(lockReqBuf))
			if err != nil {
				t.Fatalf("http.Post failed: %v", err)
			}
			defer resp.Body.Close()

			readBodyAndShowResponse(resp)

			time.Sleep(400 * time.Millisecond)

			// HTTPリクエストを作成
			resp, err = http.Post(unlockURL, "application/json", bytes.NewBuffer(unlockReqBuf))
			if err != nil {
				t.Fatalf("http.Post failed: %v", err)
			}
			defer resp.Body.Close()

			readBodyAndShowResponse(resp)
		}(paths[i%3])
	}
	wg.Wait()
}
