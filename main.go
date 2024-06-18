// Copyright © 2024 Yoshiki Shibata. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var locker *Locker

func main() {
	log.Printf("Simple Locker Server v0.0")

	locker = NewLocker()

	http.HandleFunc("/", genericHandler)
	http.HandleFunc("/lock", handleLock)
	http.HandleFunc("/unlock", handleUnlock)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func genericHandler(w http.ResponseWriter, r *http.Request) {
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

	<-locker.Lock(lockRequest)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("LOCKED"))
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

	<-locker.Unlock(unlockRequest)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("UNLOCKED"))
}
