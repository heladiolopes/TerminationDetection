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

	// ackToWait: number of messages that still need to be receveid so the process can become quiet
	ackToWait int

	// currentState: persistent state on all servers
	currentState *util.ProtectedString

	// logicalClock: Lamport logical clock
	logicalClock int

	// startActive: tells if the process is initially active
	startActive bool

	// Goroutine communication channels
	// terminationTick: detect termination of a process
	terminationTick <-chan time.Time
    waveChan    chan *WaveArgs
    basicChan   chan *BasicArgs
    ackChan     chan *AckArgs
	termChan	chan *TerminationArgs
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
		termChan: make(chan *TerminationArgs, 10*len(peers)),
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

// activeSelect implements the logic to handle active processes
func (rana *Rana) activeSelect() {
	log.Println("[ACTIVE] Run Logic.")

	// ALUNO

	// timer related to active period of the current process
	resetTerminationTimeout()

	// process sends basic messages to others
	broadcastBasic()

	/*
	// timer related to active period of the current process
	rana.resetTerminationTimeout()

	// process sends basic messages to others
	rana.broadcastBasic()
	*/

	for {
		select {
		case <-rana.terminationTick:
			// ALUNO

			// timeout
			// if process is not waiting for remaining acks
			if (ackToWait == 0)
				// process becomes quiet
				currentState = quiet

				// send wave message to all processes
				broadcastWave()
				return 
			// process is waiting for remaining acks
			if (ackToWait != 0)
				// process becomes passive
				currentState = passive
				return

			/*
			if rana.ackToWait == 0 {
				rana.currentState.Set(quiet)
				log.Println("[ACTIVE] Changing to quiet.")
				// process starts wave
				rana.broadcastWave()
				return
			} else { // else process become passive
				rana.currentState.Set(passive)
				log.Println("[ACTIVE] Changing to passive.")
				return
			}
			*/
		case wave := <-rana.waveChan:
			// ALUNO

			// receiving wave
			// process is active: descard wave

			// update logical clock, if necessary
			if (logicalClock < waveClock)
				logicalClock = waveClock

			/*
			// receiving wave
			// process is active: descard wave
			log.Printf("[ACTIVE] Descarting wave from %v.", wave.Initiator)
			// update logical clock, if necessary
			if rana.logicalClock < wave.Clock {
				log.Printf("[ACTIVE] Updating clock from %v to %v.", rana.logicalClock, wave.Clock)
				rana.logicalClock = wave.Clock
			}
			*/
		case basic := <-rana.basicChan:
			// ALUNO

			// receiving basic message
			// process already active

			// send ack message
			args := &AckArgs{
				(...)
			}
			sendAck((...), args)

			/*
			// receiving basic message
			// process already active
			args := &AckArgs{
				Clock: rana.logicalClock,
				Sender: rana.me,
			}
			// send ack message
			rana.sendAck(basic.Sender, args)
			log.Printf("[ACTIVE] Reseting active state.")
			break
			*/
		case ack := <-rana.ackChan:
			// ALUNO

			// update ack to wait
			ackToWait--

			// update logical clock, if necessary
			if (logicalClock < ackClock + 1)
				logicalClock = ackClock + 1

			/*
			// update ack to wait
			rana.ackToWait--
			// update logical clock, if necessary
			if rana.logicalClock < ack.Clock + 1 {
				rana.logicalClock = ack.Clock + 1
			}
			log.Printf("[ACTIVE] Receiving ack from %v.", ack.Sender)
			*/
		}
	}
}

func (rana *Rana) passiveSelect() {
	log.Println("[PASSIVE] Run Logic.")
	for {
		select {
		case wave := <-rana.waveChan:
			// ALUNO

			// receiving wave
			// process is passive: descard wave

			// update logical clock, if necessary
			if (logicalClock < waveClock)
				logicalClock = waveClock

			/*
			// receiving wave
			// process is passive: descard wave
			log.Println("[PASSIVE] Descarting wave from", wave.Initiator)
			// update logical clock, if necessary
			if rana.logicalClock < wave.Clock {
				log.Printf("[PASSIVE] Updating clock from %v to %v", rana.logicalClock, wave.Clock)
				rana.logicalClock = wave.Clock
			}
			*/
		case basic := <-rana.basicChan:
			// ALUNO

			// receiving basic message

			// send ack message
			args := &AckArgs{
				(...)
			}
			sendAck((...), args)

			// change state to active
			currentState = active

			/*
			// receiving basic message
			log.Println("[PASSIVE] Receving basic, activating again.")
			args := &AckArgs{
				Clock: rana.logicalClock,
				Sender: rana.me,
			}
			// send ack message
			rana.sendAck(basic.Sender, args)
			// change state to active
			rana.currentState.Set(active)
			break
			*/
		case ack := <-rana.ackChan:
			// ALUNO

			// update ack to wait
			ackToWait--

			// update logical clock, if necessary
			if (logicalClock < ackClock + 1)
				logicalClock = ackClock + 1

			// if process is not waiting for remaining acks
			if (ackToWait == 0)
				// process becomes quiet
				currentState = quiet

				// send wave message to all processes
				broadcastWave()
				break 

			/*
			// update ack to wait
			rana.ackToWait--
			// update logical clock, if necessary
			if rana.logicalClock < ack.Clock + 1 {
				rana.logicalClock = ack.Clock + 1
			}
			log.Println("[PASSIVE] Receiving ack from", ack.Sender)
			// if process is not waiting for remaining acks
			// process becomes quiet
			if rana.ackToWait == 0 {
				log.Println("[PASSIVE] Changing to quiet.")
				rana.currentState.Set(quiet)
				// process starts wave
				rana.broadcastWave()
				break
			}
			*/
		}
	}
}

func (rana *Rana) quietSelect() {
	log.Println("[QUIET] Run Logic.")
	for {
		select {
		case wave := <-rana.waveChan:
			// ALUNO

			// receiving wave
			// if wave logical clock smaller smaller than process logical clock
			if (waveClockc < logicalClock)
				// discard wave

			// if logical clock smaller or equal than wave logical clock
			else
				// update logical clock
				logicalClock = waveClock

				// if process already signed boolean array
				if (sign)
					// if process is the initiator
					if (initiator = me)
						// if array is signed by all process
						if (Signature == true)
							// send termination message to all processes
							broadcastTermination()


				// if process did not sign boolean array
				else
					// check and print processes that have already signed

					// sign boolean array
					Signature[me] = true

					// process continue wave
					// for all other processes, send wave
					go func(peer int) {
						sendWave(peer, wave)
					}()

			/*
			// receiving wave
			// if wave logical clock smaller smaller than process logical clock
			// discard wave
			if wave.Clock < rana.logicalClock {
				log.Printf("[QUIET] Rejecting wave from %v. WaveClock: %v < MyClock: %v", wave.Initiator, wave.Clock, rana.logicalClock)
			} else {
				// update logical clock
				rana.logicalClock = wave.Clock
				sign := wave.Signature[rana.me]
				// if process already signed boolean array
				if sign {
					// process is the initiator of the wave
					if wave.Initiator == rana.me {
						everybody := true
						// check if array is signed by all processes
						for peerIndex := range rana.peers {
							if wave.Signature[peerIndex] == false {
								everybody = false
								break
							}
						}
						// array signed by all processes
						if everybody {
							// log.Println("[QUIET] Termination deteced by me.")
							// process send termination message to all processes
							rana.broadcastTermination()
							// os.Exit(0)
						}
					}
				} else { // process has not signed boolean array
					// check processes that have already signed
					signedBy := make([]int, 0)
					for peerIndex := range rana.peers {
						if wave.Signature[peerIndex] {
							signedBy = append(signedBy, peerIndex)
						}
					}
					log.Printf("[QUIET] Acepting wave '%v' from %v.", signedBy, wave.Initiator)
					// sign array
					wave.Signature[rana.me] = true
					// process continue wave
					for peerIndex := range rana.peers {
						if peerIndex != rana.me {
							go func(peer int) {
								rana.sendWave(peer, wave)
							}(peerIndex)
						}
					}
				}
			}
			*/
		case basic := <-rana.basicChan:
			// ALUNO

			// receiving basic message

			// send ack message
			args := &AckArgs{
				(...)
			}
			sendAck((...), args)

			// change state to active
			currentState = active

			/*
			// receiving basic message
			log.Printf("[QUIET] Receving basic from %v, activating again.", basic.Sender)
			args := &AckArgs{
				Clock: rana.logicalClock,
				Sender: rana.me,
			}
			// send ack message
			rana.sendAck(basic.Sender, args)
			// change state to active
			rana.currentState.Set(active)
			return
			*/
		case termination := <-rana.termChan:
			// ALUNO

			// receiving termination message
			// exit processes
			Exit()

			/*
			// receiving termination message
			log.Printf("[QUIET] Termination deteced by %v.", termination.Initiator)
			// exit process
			os.Exit(0)
			*/
		}
	}
}
