package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"docker_test/comon"
	dbM "docker_test/database/mgo"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Mongo db
	dbM.MgoDB.Init()

	go startserver()
	closeServer()

	// Keep alive
	for {
		time.Sleep(1 * time.Second)
	}

}

func startserver() {

	routersInit := InitRouter()
	listenAddr := fmt.Sprintf("0.0.0.0:%d", 3000)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           listenAddr,
		Handler:        routersInit,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}
	fmt.Println("InitHttpServer", listenAddr)
	server.SetKeepAlivesEnabled(false)
	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("Server err: %v", err)
	}
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/receiveMessage", comon.LineReceive)
	r.POST("/sendMessage", comon.LineSend)
	r.GET("/getMessages", comon.GetLineMessages)
	return r
}

func closeServer() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c

		dbM.MgoDB.Close()

		time.Sleep(1 * time.Second)
		fmt.Println("Close Database")
		os.Exit(0)
	}()
}
