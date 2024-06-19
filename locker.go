// Copyright © 2024 Yoshiki Shibata. All rights reserved.

package main

import (
	"log"
	"net/http"
)

type LockerRequest struct {
	UUID     string
	Path     string
	response chan int
}

type Locker struct {
	lock   chan *LockerRequest
	unlock chan *LockerRequest
	clear  chan *LockerRequest

	lockedMap    map[string]string
	lockRequests []*LockerRequest
}

func NewLocker() *Locker {
	locker := &Locker{
		lock:      make(chan *LockerRequest, 10),
		unlock:    make(chan *LockerRequest, 10),
		clear:     make(chan *LockerRequest),
		lockedMap: make(map[string]string),
	}
	go locker.monitor()
	return locker
}

func (l *Locker) Lock(request LockRequest) chan int {
	response := make(chan int)
	l.lock <- &LockerRequest{
		UUID:     request.UUID,
		Path:     request.Path,
		response: response,
	}

	return response
}

func (l *Locker) Unlock(request UnlockRequest) chan int {
	response := make(chan int)
	l.unlock <- &LockerRequest{
		UUID:     request.UUID,
		Path:     request.Path,
		response: response,
	}

	return response
}

func (l *Locker) ClearAll() chan int {
	response := make(chan int)
	l.clear <- &LockerRequest{
		response: response,
	}

	return response
}

func (l *Locker) monitor() {
	for {
		select {
		case lockReq := <-l.lock:
			if lockReq.UUID == "" || lockReq.Path == "" {
				log.Printf("Invalid Lock Request[%s:%s]: Ignored", lockReq.Path, lockReq.UUID)
				continue
			}

			l.lockRequests = append(l.lockRequests, lockReq)

			l.rescanLockRequests()
		case unlockReq := <-l.unlock:
			if unlockReq.UUID == "" || unlockReq.Path == "" {
				log.Printf("Invalid Unlock Request[%s:%s]: Ignored", unlockReq.Path, unlockReq.UUID)
				continue
			}

			lockedUUID, ok := l.lockedMap[unlockReq.Path]
			if !ok {
				log.Printf("Invalid Unlock Request[%s:%s]: Ignored", unlockReq.Path, unlockReq.UUID)
				continue
			}
			if lockedUUID != unlockReq.UUID {
				log.Printf("Invalid Unlock Request[%s:%s]: Ignored", unlockReq.Path, unlockReq.UUID)
				continue
			}
			delete(l.lockedMap, unlockReq.Path)
			unlockReq.response <- http.StatusOK
			close(unlockReq.response)
			log.Printf("[%s:%s] Unlocked", unlockReq.Path, unlockReq.UUID)

			l.rescanLockRequests()
		case clearReq := <-l.clear:
			for _, request := range l.lockRequests {
				request.response <- http.StatusGone
				close(request.response)
			}
			l.lockRequests = nil
			clearReq.response <- http.StatusOK
			close(clearReq.response)
			clear(l.lockedMap)
		}
	}
}

func (l *Locker) rescanLockRequests() {
	newLockRequests := make([]*LockerRequest, 0, len(l.lockRequests))
	for _, request := range l.lockRequests {
		_, ok := l.lockedMap[request.Path]
		if !ok {
			l.lockedMap[request.Path] = request.UUID
			request.response <- http.StatusOK
			close(request.response)
			log.Printf("[%s:%s] Locked", request.Path, request.UUID)
			continue
		}
		newLockRequests = append(newLockRequests, request)
	}
	l.lockRequests = newLockRequests
	log.Printf("rescanLockRequests: len(l.lockRequests): %d",
		len(l.lockRequests))
}
