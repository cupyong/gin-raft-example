package common

import (
	"flag"
)

type Options struct {
	DataDir        string // data directory
	HttpAddress    string // http server address
	RaftTCPAddress string // construct Raft Address
	Bootstrap      bool   // start as master or not
	JoinAddress    string // peer address to join
}

func NewOptions() *Options {
	opts := &Options{}
	var httpAddress = flag.String("http", "6000", "Http address")
	var raftTCPAddress = flag.String("raft", "127.0.0.1:7000", "raft tcp address")
	var node = flag.String("node", "node1", "raft node name")
	var bootstrap = flag.Bool("bootstrap", false, "start as raft cluster")
	var joinAddress = flag.String("join", "", "join address for raft cluster")
	flag.Parse()

	opts.DataDir = "./" + *node
	opts.HttpAddress = *httpAddress
	opts.Bootstrap = *bootstrap
	opts.RaftTCPAddress = *raftTCPAddress
	opts.JoinAddress = *joinAddress
	return opts
}
