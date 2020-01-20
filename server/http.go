package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/hashicorp/raft"
	"encoding/json"
	"time"
	"net/http"
	"io/ioutil"
	"errors"
	"fmt"
)

type HttpService struct {
	Options *Options
	Cache   *Cache
	raft    *raft.Raft
	fsm     *FSM
	LeaderNotifyCh chan bool
}

func NewHttpService(h *HttpService) *HttpService {
	return h
}

func (h *HttpService) Get(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		log.Println("doGet() error, get nil key")
		c.String(200, "null")
		return
	}
	log.Println("Key == ", )
	log.Println("Value == ", h.Cache.Get(key))
	c.String(200, h.Cache.Get(key))
}

func (h *HttpService) Set(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	if key == "" || value == "" {
		c.String(200, "null")
		return
	}

	event := logEntryData{Key: key, Value: value}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		c.String(200, "null")
		return
	}

	applyFuture := h.raft.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		c.String(200, "null")
		return
	}
	c.String(200, "ok")
	return
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

func JoinRaftCluster(opts *Options) error {
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
