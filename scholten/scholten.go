package scholten

import (
	"errors"
	"TerminationDetection/util"
	"log"
	"sync"
	"time"
)

// Scholten is the struct that hold all information that is used by this instance
// of scholten.
type Scholten struct {
	sync.Mutex

	serv *server
	done chan struct{}

	peers map[int]string
	me    int

	ccp 		int
	children	[]int
	dad			int
	root 		bool

	// Persistent state on all servers:
	// logicalClock: Lamport logical clock
	currentState   *util.ProtectedString

	// Goroutine communication channels
    basicChan   chan *BasicArgs
    controlChan	chan *ControlArgs
}

// NewScholten create a new scholten object and return a pointer to it.
func NewScholten(peers map[int]string, me int, isroot bool) *Scholten {
	var err error

	// 0 is reserved to represent undefined vote/leader
	if me == 0 {
		panic(errors.New("Reserved instanceID('0')"))
	}

	scholten := &Scholten{
		done: make(chan struct{}),

		peers: peers,
		me:    me,

		ccp: 0,
		children: make(int, 0),
		dad: -1,
		root: isroot,

		currentState: util.NewProtectedString(),

		basicChan: make(chan *BasicArgs, 10*len(peers)),
		controlChan: make(chan *ControlArgs, 10*len(peers)),
	}

	scholten.serv, err = newServer(scholten, peers[me])
	if err != nil {
		panic(err)
	}

	if isroot {
		scholten.currentState.Set(active)
	} else {
		scholten.currentState.Set(passive)
	}

	go scholten.loop()

	return scholten
}

// Done returns a channel that will be used when the instance is done.
func (scholten *Scholten) Done() <-chan struct{} {
	return scholten.done
}

// All changes to Scholten structure should occur in the context of this routine.
func (scholten *Scholten) loop() {

	err := scholten.serv.startListening()
	if err != nil {
		panic(err)
	}

	for {
		switch scholten.currentState.Get() {
        case active:
            scholten.activeSelect()
        case passive:
            scholten.passiveSelect()
		}
	}
}

func (scholten *Scholten) activeSelect() {
	log.Println("[ACTIVE] Run Logic.")
	scholten.resetTerminationTimeout()
	for {
		select {
		case basic := <-scholten.basicChan:
			// ALUNO
		case control := <-scholten.controlChan:
			// ALUNO
		}
	}
}

func (scholten *Scholten) passiveSelect() {
	log.Println("[PASSIVE] Run Logic.")
	for {
		select {
		case basic := <-scholten.basicChan:
			// ALUNO
		case control := <-scholten.controlChan:
			// ALUNO
		}
	}
}
