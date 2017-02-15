package diskv

import (
	"shardmaster"
)

//
// Sharded key/value server.
// Lots of replica groups, each running op-at-a-time paxos.
// Shardmaster decides which group serves each shard.
// Shardmaster may change shard assignment from time to time.
//
// You will have to modify these definitions.
//

const (
	OK            = "OK"
	ErrNoKey      = "ErrNoKey"
	ErrWrongGroup = "ErrWrongGroup"
)

type Err string

type PutAppendArgs struct {
	Key       string
	Value     string
	Op        string // "Put" or "Append"
	ID        string
	Seq       int
	ConfigNum int
	Shard     int
}

type PutAppendReply struct {
	Err Err
}

type GetArgs struct {
	Key       string
	ID        string
	Seq       int
	ConfigNum int
	Shard     int
}

type GetReply struct {
	Err   Err
	Value string
}

type UpdateArgs struct {
	Database     map[string]string
	MaxClientSeq map[string]string
	Shard        int
	ConfigNum    int
	Seq          int
	ID           string
}

type UpdateReply struct {
	Err string
}

type RecoverArgs struct {
}

type RecoverReply struct {
	Database     map[int]map[string]string
	MaxClientSeq map[int]map[string]string
	Config       shardmaster.Config
	LastApply    int
	IsReceived   map[int]bool
}
