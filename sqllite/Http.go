package sqllite

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin-raft-example/common"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type HttpService struct {
	Options        *common.Options
	DataBases      *DataBase
	raft           *raft.Raft
	fsm            *FSM
	LeaderNotifyCh chan bool
}

func NewHttpServiceSql(h *HttpService) *HttpService {
	return h
}

func (h *HttpService) Sql(c *gin.Context) {
	sql := c.Query("sql")
	if len(sql)<1{
		return
	}
	sql = strings.Trim(strings.ToLower(sql), " ")
	if strings.HasPrefix(sql, "update") || strings.HasPrefix(sql, "delete") || strings.HasPrefix(sql, "create")  || strings.HasPrefix(sql, "insert") {
        logId :=h.DataBases.AddId()
		event := logEntryData{Sql: sql,LogId: logId}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			c.String(200, err.Error())
			return
		}

		applyFuture := h.raft.Apply(eventBytes, 5*time.Second)
		if err := applyFuture.Error(); err != nil {
			fmt.Println(err)
			c.String(200, err.Error())
			return
		}
		c.String(200, "ok")
		return
	} else {
		table := h.DataBases.Query(sql)
		jsonStr, err := json.Marshal(table)
		fmt.Println(err)
		c.String(200, string(jsonStr))
	}
}

func (h *HttpService) Join(c *gin.Context) {
	peerAddress := c.Query("peerAddress")
	if peerAddress == "" {
		log.Println("invalid PeerAddress")
		c.String(200, "null")
		return
	}
	addPeerFuture := h.raft.AddVoter(raft.ServerID(peerAddress), raft.ServerAddress(peerAddress), 0, 0)
	if err := addPeerFuture.Error(); err != nil {
		log.Printf("Error joining peer to raft, peeraddress:%s, err:%v, code:%d", peerAddress, err, http.StatusInternalServerError)
		c.String(200, "null")
		return
	}
	c.String(200, "ok")
	return
}

func JoinRaftClustersql(opts *common.Options) error {
	url := fmt.Sprintf("http://%s/join?peerAddress=%s", opts.JoinAddress, opts.RaftTCPAddress)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "ok" {
		return errors.New(fmt.Sprintf("Error joining cluster: %s", body))
	}

	return nil
}
