package server

import (
	"encoding/json"
	"io"
	"log"
	"github.com/hashicorp/raft"
)
type fsmSnapshot struct {
	name map[string]string
}

type FSM struct {
	Cache   *Cache
}

type logEntryData struct {
	Key   string
	Value string
}


type PostRequest struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// Apply applies a Raft log entry to the key-value store.
func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	e := logEntryData{}
	if err := json.Unmarshal(logEntry.Data, &e); err != nil {
		panic("Failed unmarshaling Raft log entry. This is a bug.")
	}
	ret := f.Cache.Set(e.Key, e.Value)
	log.Printf("fms.Apply(), logEntry:%s, ret:%v\n", logEntry.Data, ret)
	return ret
}

// Snapshot returns a latest snapshot
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{cache: f.Cache}, nil
}

// Restore stores the key-value store to a previous state.
func (f *FSM) Restore(serialized io.ReadCloser) error {
	return f.Cache.UnMarshal(serialized)
}
