package scholten

import (
	"TerminationDetection/util"
	"errors"
	"sync"
	"log"
)

// Scholten is the struct that hold all information that is used by this instance
// of scholten.
type Scholten struct {
	sync.Mutex

	serv *server
	done chan int

	peers map[int]string
	me    int

	ccp      int
	children []int
	dad      int
	root     bool

	works 	 int
	// Persistent state on all servers:
	currentState *util.ProtectedString

	// Goroutine communication channels
	basicChan   chan *BasicArgs
	controlChan chan *ControlArgs
	terminationChan  chan *TerminationArgs
}

// NewScholten create a new scholten object and return a pointer to it.
func NewScholten(peers map[int]string, me int, isroot bool) *Scholten {
	var err error

	// 0 is reserved to represent undefined vote/leader
	if me == 0 {
		panic(errors.New("Reserved instanceID('0')"))
	}

	scholten := &Scholten{
		done: make(chan int),

		peers: peers,
		me:    me,

		ccp:      0,
		children: make([]int, 10),
		dad:      -1,
		root:     isroot,

		currentState: util.NewProtectedString(),

		basicChan:   make(chan *BasicArgs, 10*len(peers)),
		controlChan: make(chan *ControlArgs, 10*len(peers)),
    terminationChan:  make(chan *TerminationArgs, 10*len(peers)),
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
func (scholten *Scholten) Done() <-chan int {
	return scholten.done
}

// All changes to Scholten structure should occur in the context of this routine.
func (scholten *Scholten) loop() {

	err := scholten.serv.startListening()
	if err != nil {
		panic(err)
	}
	if scholten.root {
		go scholten.doWork()
	}

	for {
		select {
		case basic := <-scholten.basicChan:
			// ALUNO
			if scholten.root || scholten.dad != -1 {
				scholten.sendControl(basic.Sender)
			} else {
				scholten.dad = basic.Sender
			}

			go scholten.doWork()

		case control := <-scholten.controlChan:
			// ALUNO
			for i, child := range scholten.children {
				if child == control.Sender {
					scholten.removeChild(child, i)
					break
				}
			}

			if scholten.ccp == 0 && scholten.currentState.Get() == passive {
				scholten.leaveTree()
			}

		case termination := <- scholten.terminationChan:
			// ALUNO
			log.Println("Received termination message from", termination.Sender,"!")
			scholten.done <- 0
		}
	}
}
