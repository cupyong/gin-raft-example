package sqllite

import (
	"github.com/hashicorp/raft"
)

type snapshotsql struct {
	dataBase *DataBase
}

// Persist saves the FSM snapshot out to the given sink.
func (s *snapshotsql) Persist(sink raft.SnapshotSink) error {
	//snapshotBytes, err := s.dataBase.Marshal()
	//if err != nil {
	//	sink.Cancel()
	//	return err
	//}
	//
	//if _, err := sink.Write(snapshotBytes); err != nil {
	//	sink.Cancel()
	//	return err
	//}
	//
	//if err := sink.Close(); err != nil {
	//	sink.Cancel()
	//	return err
	//}
	return nil
}

func (f *snapshotsql) Release() {}
