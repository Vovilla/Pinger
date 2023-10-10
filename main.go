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

func runPinger(c *gin.Context) {

    PING_TIMEOUT := 5
    ping_timeout := time.Duration(PING_TIMEOUT) * time.Second

    wg := new(sync.WaitGroup)

    hosts := make([]string, 1)
    c.BindJSON(&hosts)

    results := make(map[string]bool)

    for _, addr := range hosts {
        wg.Add(1)
        result := make(chan bool, 1)
        go pinger(addr, wg, ping_timeout, result)
        results[addr] = <- result
    }
    wg.Wait()

    c.IndentedJSON(http.StatusOK, results)
}

func doPing(addr string) {
    pinger, err := ping.NewPinger(addr)
    if err != nil {
        panic(err)
    }
    pinger.Count = 5
    err = pinger.Run()
    if err != nil {
        panic(err)
    }
}

func pinger(addr string, wg *sync.WaitGroup, t time.Duration, result chan bool) {
	c := make(chan bool)
	go func() {
		defer close(c)
        doPing(addr)
	}() 
	select {
        case <-c:
            result <- true
        case <-time.After(t):
            result <- false
	}  
    wg.Done() 
}

func main() {

    router := gin.Default()
    router.POST("/pinger", runPinger)
    router.Run("localhost:8080")

}



