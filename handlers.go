// Copyright © 2024 Yoshiki Shibata. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func genericHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleLock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var lockRequest LockRequest
	if err := json.Unmarshal(body, &lockRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Unmarshal request failed: %v", err))
		return
	}

	statusCode := <-locker.Lock(lockRequest)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusOK {
		w.Write([]byte("LOCKED"))
	}
}

func handleUnlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var unlockRequest UnlockRequest
	if err := json.Unmarshal(body, &unlockRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Unmarshal request failed: %v", err))
		return
	}

	statusCode := <-locker.Unlock(unlockRequest)

	w.WriteHeader(statusCode)
	if statusCode == http.StatusOK {
		w.Write([]byte("UNLOCKED"))
	}
}

func handleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var clearRequest ClearRequest
	if err := json.Unmarshal(body, &clearRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Unmarshal request failed: %v", err))
		return
	}

	if clearRequest.ClearAll {
		statusCode := <-locker.ClearAll()
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			w.Write([]byte("All Requests Cleared"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("NO Requests Cleared: clear_all is false"))
}
