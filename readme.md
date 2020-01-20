raft demo

go run *.go -bootstrap true  

go run *.go -http 6001 -raft 127.0.0.1:7001 -node node2  --join=127.0.0.1:6000 

go run *.go -http 6002 -raft 127.0.0.1:7002 -node node3  --join=127.0.0.1:6000

启动上面三个节点

测试：
 
 curl http://localhost:6000/set\?key\=test8\&value\=test1
   
curl http://localhost:6001/get\?key\=test8