package pbservice

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"
	"viewservice"
)

type Clerk struct {
	vs      *viewservice.Clerk
	me      string
	primary string
	// Your declarations here
}

// this may come in handy.
func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(vshost string, me string) *Clerk {
	ck := new(Clerk)
	ck.vs = viewservice.MakeClerk(me, vshost)
	ck.primary = ""
	ck.me = strconv.Itoa(int(nrand()))

	return ck
}

//
// fetch a key's value from the current primary;
// if they key has never been set, return "".
// Get() must keep trying until it either the
// primary replies with the value or the primary
// says the key doesn't exist (has never been Put().
//
func (ck *Clerk) Get(key string) string {
	for ck.primary == "" {
		ck.primary = ck.vs.Primary()
		time.Sleep(viewservice.PingInterval)
	}
	args := &GetArgs{Key: key}
	reply := GetReply{}
	ok := call(ck.primary, "PBServer.Get", args, &reply)
	for reply.Err != OK || !ok {
		time.Sleep(viewservice.PingInterval)
		ck.primary = ck.vs.Primary()
		ok = call(ck.primary, "PBServer.Get", args, &reply)
	}

	return reply.Value
}

//
// send a Put or Append RPC
//
func (ck *Clerk) PutAppend(key string, value string, op string) {

	for ck.primary == "" {
		ck.primary = ck.vs.Primary()
		time.Sleep(viewservice.PingInterval)
	}
	args := &PutAppendArgs{Key: key, Value: value, Operation: op, From: ck.me}
	reply := PutAppendReply{}
	ok := call(ck.primary, "PBServer.PutAppend", args, &reply)
	for reply.Err != OK || !ok {
		time.Sleep(viewservice.PingInterval)
		ck.primary = ck.vs.Primary()
		ok = call(ck.primary, "PBServer.PutAppend", args, &reply)
	}
}

//
// tell the primary to update key's value.
// must keep trying until it succeeds.
//
func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

//
// tell the primary to append to key's value.
// must keep trying until it succeeds.
//
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
