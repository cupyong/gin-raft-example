package sqllite

import (
	"fmt"
	"github.com/hashicorp/raft"
	"io"
	"log"
	"encoding/json"
)

type FSM struct {
	DataBases *DataBase
}

type logEntryData struct {
	Sql string
	LogId int64
}

// Apply applies a Raft log entry to the key-value store.
func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	e := logEntryData{}

	if err := json.Unmarshal(logEntry.Data, &e); err != nil {
		panic("Failed unmarshaling Raft log entry. This is a bug.")
	}
	if len(e.Sql)>0 {
		fmt.Println(e.LogId)
		ret := f.DataBases.Prepare(e.Sql,e.LogId)
		log.Printf("fms.Apply(), logEntry:%s, ret:%v\n", logEntry.Data, ret)
		return ret
	}
	return nil
}

// Snapshot returns a latest snapshot
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshotsql{dataBase: f.DataBases}, nil
}

// Restore stores the key-value store to a previous state.
func (f *FSM) Restore(serialized io.ReadCloser) error {
	return nil
}
