package main

import (
	"crypto/sha512"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"time"
)

type doHashOp struct {
	password string
	delayS   int
	resp     chan int
}

type getHashOp struct {
	index int
	resp  chan string
}

type getStatsOp struct {
	pattern string
	resp    chan *statsObj
}

type shutdownOp struct {
	should bool
	resp   chan bool
}

type statsObj struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
}

type setValueOp struct {
	index     int
	hashValue string
}

type logResponseOp struct {
	pattern string
	time    int
}

type pendingHashOp struct {
	resp chan int
}

type Server interface {
	getStats(pattern string) *statsObj
	getPendingHashCt() int
	doHash(password string, delayS int) int
	getHash(index int) string
	logResponse(pattern string, time int)
	isShutdown() bool
	shutdown()
}

type ServerState struct {
	getHashC       chan *getHashOp
	doHashC        chan *doHashOp
	logResponseC   chan *logResponseOp
	getStatsC      chan *getStatsOp
	shutdownOpC    chan *shutdownOp
	pendingHashesC chan *pendingHashOp
}

func (state *ServerState) getPendingHashCt() int {
	data := &pendingHashOp{resp: make(chan int)}
	state.pendingHashesC <- data
	return <-data.resp
}

func (state *ServerState) getHash(index int) string {
	data := &getHashOp{
		index: index,
		resp:  make(chan string)}
	state.getHashC <- data
	return <-data.resp
}

func (state *ServerState) shutdown() {
	state.shutdownOpC <- &shutdownOp{should: true}
}

func (state *ServerState) isShutdown() bool {
	data := &shutdownOp{should: false, resp: make(chan bool)}
	state.shutdownOpC <- data
	return <-data.resp
}

func (state *ServerState) logResponse(pattern string, time int) {
	state.logResponseC <- &logResponseOp{pattern: pattern, time: time}
}

func (state *ServerState) getStats(pattern string) *statsObj {
	data := &getStatsOp{pattern: pattern, resp: make(chan *statsObj)}
	state.getStatsC <- data
	return <-data.resp
}

func (state *ServerState) doHash(password string, delayS int) int {
	data := &doHashOp{password: password, resp: make(chan int), delayS: delayS}
	state.doHashC <- data
	return <-data.resp
}

func computeHash(password string, index int, setVal chan *setValueOp, delayS int) {
	time.Sleep(time.Duration(delayS) * time.Second)
	passwordByteValue := []byte(password)
	hasher := sha512.New()
	hasher.Write(passwordByteValue)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	log.Printf("computed hash #%s value of %s", strconv.Itoa(index), sha)
	setVal <- &setValueOp{index: index, hashValue: sha}
}

func makeState() Server {
	doHashC := make(chan *doHashOp)
	getHashC := make(chan *getHashOp)
	logResponseC := make(chan *logResponseOp)
	getStatsC := make(chan *getStatsOp)
	shutdownOpC := make(chan *shutdownOp)
	pendingHashesC := make(chan *pendingHashOp)
	go func() {
		var hashes []string
		timings := make(map[string]int)
		counts := make(map[string]int)
		pending := make(map[int]bool)
		setValueC := make(chan *setValueOp)
		isShuttingDown := false
		for {
			select {
			case req := <-pendingHashesC:
				req.resp <- len(pending)
			case req := <-shutdownOpC:
				if !req.should {
					req.resp <- isShuttingDown
					breaks
				}
				isShuttingDown = true
				if len(pending) == 0 {
					log.Print("shutting down gracefully")
					os.Exit(0)
				}
			case hashReq := <-doHashC:
				hashes = append(hashes, "")
				index := len(hashes)
				pending[index] = true
				hashReq.resp <- index
				go computeHash(hashReq.password, index, setValueC, hashReq.delayS)
			case setReq := <-setValueC:
				hashes[setReq.index-1] = setReq.hashValue
				delete(pending, setReq.index)
				if len(pending) == 0 && isShuttingDown {
					log.Print("shutting down gracefully")
					os.Exit(0)
				}

			case getHashReq := <-getHashC:
				if getHashReq.index-1 >= len(hashes) || getHashReq.index-1 < 0 {
					getHashReq.resp <- "not found"
				} else {
					getHashReq.resp <- hashes[getHashReq.index-1]
				}
			case timing := <-logResponseC:
				timings[timing.pattern] += timing.time
				counts[timing.pattern]++
			case statsReq := <-getStatsC:
				time := timings[statsReq.pattern]
				count := counts[statsReq.pattern]
				if count == 0 {
					statsReq.resp <- &statsObj{Total: 0, Average: 0}
				} else {
					statsReq.resp <- &statsObj{Total: count, Average: float64(time) / float64(count)}
				}
			}
		}
	}()

	return &ServerState{doHashC: doHashC, getHashC: getHashC, logResponseC: logResponseC, getStatsC: getStatsC, shutdownOpC: shutdownOpC, pendingHashesC: pendingHashesC}
}
