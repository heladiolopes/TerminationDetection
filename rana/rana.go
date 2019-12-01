package rana

import (
	"errors"
	"TerminationDetection/util"
	"log"
	"sync"
	"time"
	"os"
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

	rana.broadcastBasic()

	for {
		select {
		case <-rana.terminationTick:
			// ALUNO
			if rana.ackToWait == 0 {
				rana.currentState.Set(quiet)
				log.Println("[ACTIVE] Changing to quiet.")
				rana.broadcastWave()
				return
			} else {
				rana.currentState.Set(passive)
				log.Println("[ACTIVE] Changing to passive.")
				return
			}
		case wave := <-rana.waveChan:
			// ALUNO
			if rana.logicalClock < wave.Clock {
				rana.logicalClock = wave.Clock
				log.Printf("[ACTIVE] Updating clock to %v", rana.logicalClock)
			}

			log.Println("[ACTIVE] Descarting wave from", wave.Initiator)
		case basic := <-rana.basicChan:
			// ALUNO

			args := &AckArgs{
				Clock: rana.logicalClock,
				Sender: rana.me,
			}

			rana.sendAck(basic.Sender, args)

			log.Println("[ACTIVE] Reseting active state")
			break

		case ack := <-rana.ackChan:
			// ALUNO
			rana.ackToWait--
			if rana.logicalClock < ack.Clock + 1 {
				rana.logicalClock = ack.Clock + 1
			}
			log.Println("[ACTIVE] Receiving ack from", ack.Sender)
		}
	}
}

func (rana *Rana) passiveSelect() {
	log.Println("[PASSIVE] Run Logic.")
	for {
		select {
		case wave := <-rana.waveChan:
			// ALUNO

			if rana.logicalClock < wave.Clock {
				rana.logicalClock = wave.Clock
				log.Printf("[PASSIVE] Updating clock to %v", rana.logicalClock)
			}

			log.Println("[PASSIVE] Descarting wave from", wave.Initiator)
		case basic := <-rana.basicChan:
			// ALUNO
			log.Println("[PASSIVE] Receving basic, activating again.")

			args := &AckArgs{
				Clock: rana.logicalClock,
				Sender: rana.me,
			}

			rana.sendAck(basic.Sender, args)

			rana.currentState.Set(active)
			break

		case ack := <-rana.ackChan:
			// ALUNO
			rana.ackToWait--
			if rana.logicalClock < ack.Clock + 1 {
				rana.logicalClock = ack.Clock + 1
			}
			log.Println("[PASSIVE] Receiving ack from", ack.Sender)
			if rana.ackToWait == 0 {
				rana.currentState.Set(quiet)
				rana.broadcastWave()
				break
			}
		}
	}
}

func (rana *Rana) quietSelect() {
	log.Println("[QUIET] Run Logic.")

	for {
		select {
		case wave := <-rana.waveChan:
			// ALUNO
			log.Println(wave)
			// os.Exit(3)

			if wave.Clock < rana.logicalClock {
				log.Println("[QUIET] Rejecting wave from", wave.Initiator, "because clock smaller them mine.")
			} else {
				rana.logicalClock = wave.Clock
				sign := wave.Signature[rana.me]
				if sign {
					if wave.Initiator == rana.me {
						everybody := true
						for peerIndex := range rana.peers {
							if wave.Signature[peerIndex] == false {
								everybody = false
								break
							}
						}
						if everybody {
							log.Println("[QUIET] Termination deteced.")
							os.Exit(0)
							// avisar aos outros processos
						}
					}
				} else {
					log.Printf("[QUIET] Acepting wave from %v.", wave.Initiator)
					wave.Signature[rana.me] = true

					for peerIndex := range rana.peers {
						if peerIndex != rana.me {
							go func(peer int) {
								rana.sendWave(peer, wave)
							}(peerIndex)
						}
					}

				}
			}

		case basic := <-rana.basicChan:
			// ALUNO
			log.Printf("[QUIET] Receving basic from %v, activating again.", basic.Sender)

			args := &AckArgs{
				Clock: rana.logicalClock,
				Sender: rana.me,
			}

			rana.sendAck(basic.Sender, args)

			rana.currentState.Set(active)
			return

		}
	}
}
