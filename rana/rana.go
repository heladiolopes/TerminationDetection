package rana

import (
	"errors"
	"TerminationDetection/util"
	"log"
	"sync"
	"time"
)

// Rana is the struct that hold all information that is used by this instance
// of rana.
type Rana struct {
	sync.Mutex

	serv *server
	done chan struct{}

	peers map[int]string
	me    int

	ackToWait int

	// Persistent state on all servers:
	// logicalClock: Lamport logical clock
	currentState   *util.ProtectedString
	logicalClock   int
	startActive		 bool
	// Goroutine communication channels
	waveTick    <-chan time.Time		// ?
	terminationTick <-chan time.Time	// Detect termination of a process
    waveChan    chan *WaveArgs
    basicChan   chan *BasicArgs
    ackChan     chan *AckArgs
}

// NewRana create a new rana object and return a pointer to it.
func NewRana(peers map[int]string, me int, activeStart bool) *Rana {
	var err error

	// 0 is reserved to represent undefined vote/leader
	if me == 0 {
		panic(errors.New("Reserved instanceID('0')"))
	}

	rana := &Rana{
		done: make(chan struct{}),

		peers: peers,
		me:    me,

		ackToWait: 0,

		currentState: util.NewProtectedString(),
		logicalClock:  0,

		waveChan: make(chan *WaveArgs, 10*len(peers)),
		basicChan: make(chan *BasicArgs, 10*len(peers)),
		ackChan: make(chan *AckArgs, 10*len(peers)),
	}

	rana.serv, err = newServer(rana, peers[me])
	if err != nil {
		panic(err)
	}
	if activeStart {
		rana.currentState.Set(active)
	} else {
		rana.currentState.Set(quiet)
	}

	go rana.loop()

	return rana
}

// Done returns a channel that will be used when the instance is done.
func (rana *Rana) Done() <-chan struct{} {
	return rana.done
}

// All changes to Rana structure should occur in the context of this routine.
func (rana *Rana) loop() {

	err := rana.serv.startListening()
	if err != nil {
		panic(err)
	}

  // TODO: lembrar de colocar parâmetro
	// para usuário decidir se esta ativo

	for {
		switch rana.currentState.Get() {
        case active:
            rana.activeSelect()
        case passive:
            rana.passiveSelect()
		case quiet:
			rana.quietSelect()
		}
	}
}

func (rana *Rana) activeSelect() {
	log.Println("[ACTIVE] Run Logic.")
	rana.resetTerminationTimeout()
	for {
		select {
		case <-rana.TerminationTick:
			// ALUNO
		case wave := <-rana.waveChan:
			// ALUNO
		case basic := <-rana.basicChan:
			// ALUNO
		case ack := <-rana.ackChan:
			// ALUNO
		}
	}
}

func (rana *Rana) passiveSelect() {
	log.Println("[PASSIVE] Run Logic.")
	for {
		select {
		case wave := <-rana.waveChan:
			// ALUNO
		case basic := <-rana.basicChan:
			// ALUNO
		case ack := <-rana.ackChan:
			// ALUNO
		}
	}
}

func (rana *Rana) quietSelect() {
	log.Println("[QUIET] Run Logic.")
	rana.resetWaveTimeout()

	// TODO: enviar wave

	for {
		select {
		case <-rana.waveTick:
			// ALUNO
		case wave := <-rana.waveChan:
			// ALUNO
		case basic := <-rana.basicChan:
			// ALUNO
	}
}
