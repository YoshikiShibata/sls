// Copyright © 2024 Yoshiki Shibata. All rights reserved.

package main

type LockRequest struct {
	UUID string `json:"uuid"`
	Path string `json:"path"`
}

type UnlockRequest struct {
	UUID string `json:"uuid"`
	Path string `json:"path"`
}
