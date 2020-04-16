package microservice

import (
	"github.com/go-redis/redis/v7"
	"github.com/qubard/claack-go/lib/util"
	"log"
	"time"
)

type ActiveRace struct {
	RacerIds  []string // A list of all the racer ids (usernames) in the race
	Capacity  int      // The maximum # of racers in the race
	Started   bool
	StartedAt time.Time
}

// We will need a cleanup routine that runs as well
// periodically that cleans up finished races from the table
// which have not been updated

// The bottleneck here is not the Redis GET/SET or pubsub time
// So this single process approach should work for now
// It is technically synchronized and each racer has access
// to every other racer.

// The only thing I am concerned about is how relaying will work
// Could we store each user id as user id -> race neighbors
// so that we emit to each of those neighbors when we
// receive UpdateRace?

// JoinQueue Algorithm:
// Note that this is all thread-safe because the poolchan is syncrhonized for only queue messages
// We will also need an unqueue message that removes users from the pool
// But we will need a lock to interact with the unpooled array

// When the user enters the queue, add them to the local unpooled array

// In a second goroutine, combine unpooled users into their desired capacity periodically if it is possible
// of races using some sort of pairing/matching algorithm
// Once a race has been generated, send all users the text and race id
// As well as map them to each other using the id, so participants[raceId] -> string[] maps to all the
// current participants of the race. So when we receive that race id (which is meant to be private)
// we can just emit that packet to the people in the race
// Also map user id -> race id for security purposes just so that we make sure they actually belong to that race
// OR just store participants[raceId] -> set string[] of user ids so we can lookup if they are in the set
// before emitting given the raceId
// Use SISMEMBER to do this

// raceId is mapped to the current local state of the race so this should be fine

// before relaying a packet to them (cause we are inherently trusting the user)
// This also avoids N set() ops for each of the users

// If we assume 10 gets/user/s we get support for 14k users per
// server before Redis throttles us (max 140k GETs/sec) with
// SISMEMBER

// There's no race condition with the emits to users in a race
// Because even if they technically receive a race packet that they
// aren't in anymore we can just drop the packet if it doesn't match
// their current race id from local state

const (
	MAX_POOL_SIZE = 100000 // Allocate enough pool space for 100,000 users
)

type LinkNode struct {
	Next  *LinkNode
	Prev  *LinkNode
	Value string
}

type LinkedList struct {
	Dummy *LinkNode
	Tail  *LinkNode
	m     map[string]*LinkNode // Assume unique keys
}

func (list *LinkedList) Push(val string) {
	tmp := LinkNode{}
	tmp.Prev = list.Tail
	tmp.Value = val
	list.Tail.Next = &tmp
	list.Tail = &tmp
	list.m[val] = &tmp
}

func (list *LinkedList) Remove(val string) {
	curr := list.m[val]
	if curr.Next != nil {
		curr.Next.Prev = curr.Prev
	}
	curr.Prev.Next = curr.Next.Next

}

type RacePool struct {
	Id       string // The identifier of the pool in Redis
	Redis    *redis.Client
	Unpooled []string
	// Inbound join pool messages
	EnqueueChan <-chan *redis.Message
	// Inbound remove from pool messages
	DequeueChan <-chan *redis.Message
}

func (pool *RacePool) PoolRacers(ticker *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-ticker.C:
			// Run the pooling algorithm described above
			// This involves generating a race for users once they are pooled
			// And removing them from the `Unpooled` list
			// For now we will pool 2 racers together..
			// Ideally we have separate pools for each game type (100v100, 1v1, 1v2, 1v3)
			log.Println("Tick")
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (pool *RacePool) Run() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go pool.PoolRacers(ticker, quit)
	defer close(quit)

	for {
		select {
		// `id` wants to join the pool
		case message := <-pool.EnqueueChan:
			id := message.Payload
			// We will assume this data is sane and is the user's id
			pool.Unpooled = append(pool.Unpooled, id)
			log.Println("pool", id)
		case message := <-pool.DequeueChan:
			id := message.Payload
			log.Println("unpool", id)
			// TODO: convert this to a linked list + a map so we can
			// O(1) removal and O(1) adding
		}
	}
}

func CreateRacePool(client *redis.Client, id string, enqChanId string, deqChanId string) *RacePool {
	return &RacePool{
		Id:          id,
		Redis:       client,
		EnqueueChan: util.CreateSubChannel(client, enqChanId),
		DequeueChan: util.CreateSubChannel(client, deqChanId),
		Unpooled:    make([]string, MAX_POOL_SIZE),
	}
}
