raft demo

go run main.go -bootstrap true  

go run main.go -http 6001 -raft 127.0.0.1:7001 -node node2  --join=127.0.0.1:6000 

go run main.go -http 6002 -raft 127.0.0.1:7002 -node node3  --join=127.0.0.1:6000

启动上面三个节点

测试：
 
curl http://localhost:6000/set\?key\=test8\&value\=test1
   
curl http://localhost:cur6001/get\?key\=test8



raft+sqlite demo

go run main2.go -bootstrap true  

go run main2.go -http 6001 -raft 127.0.0.1:7001 -node node2  --join=127.0.0.1:6000 

go run main2.go -http 6002 -raft 127.0.0.1:7002 -node node3  --join=127.0.0.1:6000

启动上面三个节点

测试：
 
启动 go run test.go运行 测试



