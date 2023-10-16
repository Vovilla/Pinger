package main

import (
    // "fmt"
    "sync"
    "time"
    // "reflect"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/go-ping/ping"
)

func NewPinger() *Pinger {
    pinger := new(Pinger)
    pinger.SetCount(5)
    pinger.SetTimeout(5)
    return  pinger
}

type Pinger struct {
    count   int 
    timeout int
}

func (p *Pinger) SetCount(count int) {
    p.count = count
}

func (p *Pinger) SetTimeout(timeout int) {
    p.timeout = timeout
}

func (p *Pinger) Run(c *gin.Context) {

    timeout := time.Duration(p.timeout) * time.Second

    wg := new(sync.WaitGroup)

    hosts := make([]string, 1)
    c.BindJSON(&hosts)

    results := make(map[string]bool)

    for _, addr := range hosts {
        wg.Add(1)
        result := make(chan bool, 1)
        go p.pinger(addr, wg, timeout, result)
        results[addr] = <- result
    }
    wg.Wait()

    c.IndentedJSON(http.StatusOK, results)
}

func (p *Pinger) doPing(addr string) {
    pinger, err := ping.NewPinger(addr)
    if err != nil {
        panic(err)
    }
    pinger.Count = p.count
    err = pinger.Run()
    if err != nil {
        panic(err)
    }
}

func (p *Pinger) pinger(addr string, wg *sync.WaitGroup, t time.Duration, result chan bool) {
	c := make(chan bool)
	go func() {
		defer close(c)
        p.doPing(addr)
	}() 
	select {
        case <-c:
            result <- true
        case <-time.After(t):
            result <- false
	}  
    wg.Done() 
}