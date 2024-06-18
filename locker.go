// Copyright © 2024 Yoshiki Shibata. All rights reserved.

package main

import "log"

type LockerRequest struct {
	UUID     string
	Path     string
	response chan struct{}
}

type Locker struct {
	lock   chan *LockerRequest
	unlock chan *LockerRequest

	lockedMap    map[string]string
	lockRequests []*LockerRequest
}

func NewLocker() *Locker {
	locker := &Locker{
		lock:      make(chan *LockerRequest, 10),
		unlock:    make(chan *LockerRequest, 10),
		lockedMap: make(map[string]string),
	}
	go locker.monitor()
	return locker
}

func (l *Locker) Lock(request LockRequest) chan struct{} {
	response := make(chan struct{})
	l.lock <- &LockerRequest{
		UUID:     request.UUID,
		Path:     request.Path,
		response: response,
	}

	return response
}

func (l *Locker) Unlock(request UnlockRequest) chan struct{} {
	response := make(chan struct{})
	l.unlock <- &LockerRequest{
		UUID:     request.UUID,
		Path:     request.Path,
		response: response,
	}

	return response
}

func (l *Locker) monitor() {
	for {
		select {
		case lockReq := <-locker.lock:
			if lockReq.UUID == "" || lockReq.Path == "" {
				log.Printf("Invalid Lock Request[%s:%s]: Ignored", lockReq.Path, lockReq.UUID)
				continue
			}

			l.lockRequests = append(l.lockRequests, lockReq)

			l.rescanLockRequests()
		case unlockReq := <-locker.unlock:
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
			close(unlockReq.response)
			log.Printf("[%s:%s] Unlocked", unlockReq.Path, unlockReq.UUID)

			l.rescanLockRequests()
		}
	}
}

func (l *Locker) rescanLockRequests() {
	newLockRequests := make([]*LockerRequest, 0, len(l.lockRequests))
	for _, request := range l.lockRequests {
		_, ok := l.lockedMap[request.Path]
		if !ok {
			l.lockedMap[request.Path] = request.UUID
			close(request.response)
			log.Printf("[%s:%s] Locked", request.Path, request.UUID)
			continue
		}
		newLockRequests = append(newLockRequests, request)
	}
	l.lockRequests = newLockRequests
}
