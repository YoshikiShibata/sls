// Copyright © 2024 Yoshiki Shibata. All rights reserved.

package main

import (
	"log"
	"net/http"
)

const version = "v0.2.0"

var locker *Locker

func main() {
	log.Printf("Simple Locker Server %s", version)

	locker = NewLocker()

	http.HandleFunc("/", genericHandler)
	http.HandleFunc("/lock", handleLock)
	http.HandleFunc("/unlock", handleUnlock)
	http.HandleFunc("/clear", handleClear)
	log.Printf("Running on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
