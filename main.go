package main

import (
	"github.com/gin-gonic/gin"
	"gin-raft-example/server"
	"log"
)

func main() {
	opts := &server.Options{}
	opts = server.NewOptions()
	cache := server.NewCache()

	r := gin.Default()
	node := &server.HttpService{
		Options: opts,
		Cache:   cache,
	}
	node.NewRaftNode(opts, cache)
	httpService := server.NewHttpService(node)
	r.GET("/get", httpService.Get)
	r.GET("/set", httpService.Set)
	r.GET("/join", httpService.Join)

	if opts.JoinAddress != ""{
         server.JoinRaftCluster(opts)
	}
	go func() {
		r.Run(":"+opts.HttpAddress) // listen and serve on 0.0.0.0:8080
	}()
    for{
		select {
		case leader := <-httpService.LeaderNotifyCh:
			if leader {
				log.Println("become leader, enable write api")
			} else {
				log.Println("become follower, close write api")
			}
		}
	}
}
