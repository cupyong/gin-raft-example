package server

import (
	"github.com/hashicorp/raft"
	"time"
	"net"
	"os"
	"github.com/hashicorp/raft-boltdb"
	"path/filepath"
)

func (h *HttpService)NewRaftNode(opts *Options, cache *Cache) error {
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(opts.RaftTCPAddress)
	raftConfig.SnapshotInterval = 20 * time.Second
	raftConfig.SnapshotThreshold = 2
	leaderNotifyCh := make(chan bool, 1)
	raftConfig.NotifyCh = leaderNotifyCh
	h.LeaderNotifyCh= leaderNotifyCh
	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", opts.RaftTCPAddress)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(addr.String(), addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return  err
	}
	fsm := &FSM{
		Cache: cache,
	}
	snapshotStore, err := raft.NewFileSnapshotStore(opts.DataDir, 1, os.Stderr)
	if err != nil {
		return  err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(opts.DataDir, "raft-log.bolt"))
	if err != nil {
		return err
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(opts.DataDir, "raft-stable.bolt"))
	if err != nil {
		return err
	}
	ra, err := raft.NewRaft(raftConfig, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return err
	}
	h.raft = ra
	h.fsm = fsm
    if opts.Bootstrap{
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raftConfig.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}
	return nil
}
