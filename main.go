package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    pinger := NewPinger()

    router := gin.Default()
    router.POST("/pinger", pinger.Run)
    router.Run("localhost:8080")

}



