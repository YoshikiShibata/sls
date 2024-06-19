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
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
