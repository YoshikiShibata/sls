// Copyright © 2024 Yoshiki Shibata. All rights reserved.

package main

import (
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
	log.Printf("Running on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
