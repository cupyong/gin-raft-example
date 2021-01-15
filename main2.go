package main

import (
	"gin-raft-example/common"
	"gin-raft-example/server"
	"gin-raft-example/sqllite"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	opts := &common.Options{}
	opts = common.NewOptions()
	database := sqllite.NewDataBase(opts.DataDir + "/foo.db")

	r := gin.Default()
	node := &sqllite.HttpService{
		Options:   opts,
		DataBases: database,
	}
	node.NewRaftNode(opts, database)
	httpService := sqllite.NewHttpServiceSql(node)
	//r.GET("/get", httpService.s)
	//r.GET("/set", httpService.Set)
	r.GET("/join", httpService.Join)
	r.GET("/sql", httpService.Sql)

	if opts.JoinAddress != "" {
		server.JoinRaftCluster(opts)
	}
	go func() {
		r.Run(":" + opts.HttpAddress) // listen and serve on 0.0.0.0:8080
	}()
	for {
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
