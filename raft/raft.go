package raft

import (
	"errors"
	"labRaft/util"
	"log"
	"sync"
	"time"
)

// Raft is the struct that hold all information that is used by this instance
// of raft.
type Raft struct {
	sync.Mutex

	serv *server
	done chan struct{}

	peers map[int]string
	me    int

	// Persistent state on all servers:
	// currentTerm latest term server has seen (initialized to 0 on first boot, increases monotonically)
	// votedFor candidateId that received vote in current term (or 0 if none)
	currentState *util.ProtectedString
	currentTerm  int
	votedFor     int

	// Goroutine communication channels
	electionTick    <-chan time.Time
	requestVoteChan chan *RequestVoteArgs
	appendEntryChan chan *AppendEntryArgs
}

// NewRaft create a new raft object and return a pointer to it.
func NewRaft(peers map[int]string, me int) *Raft {
	var err error

	// 0 is reserved to represent undefined vote/leader
	if me == 0 {
		panic(errors.New("Reserved instanceID('0')"))
	}

	raft := &Raft{
		done: make(chan struct{}),

		peers: peers,
		me:    me,

		currentState: util.NewProtectedString(),
		currentTerm:  0,
		votedFor:     0,

		requestVoteChan: make(chan *RequestVoteArgs, 10*len(peers)),
		appendEntryChan: make(chan *AppendEntryArgs, 10*len(peers)),
	}

	raft.serv, err = newServer(raft, peers[me])
	if err != nil {
		panic(err)
	}

	go raft.loop()

	return raft
}

// Done returns a channel that will be used when the instance is done.
func (raft *Raft) Done() <-chan struct{} {
	return raft.done
}

// All changes to Raft structure should occur in the context of this routine.
// This way it's not necessary to use synchronizers to protect shared data.
// To send data to each of the states, use the channels provided.
func (raft *Raft) loop() {

	err := raft.serv.startListening()
	if err != nil {
		panic(err)
	}

	raft.currentState.Set(follower)
	for {
		switch raft.currentState.Get() {
		case follower:
			raft.followerSelect()
		case candidate:
			raft.candidateSelect()
		case leader:
			raft.leaderSelect()
		}
	}
}

// followerSelect implements the logic to handle messages from distinct
// events when in follower state.
func (raft *Raft) followerSelect() {
	log.Println("[FOLLOWER] Run Logic.")
	raft.resetElectionTimeout()
	for {
		select {
		case <-raft.electionTick:
			log.Println("[FOLLOWER] Election timeout.")
			raft.currentState.Set(candidate)
			return

		case rv := <-raft.requestVoteChan:

			reply := &RequestVoteReply{
				Term: raft.currentTerm,
			}

			// Quem mandou quer ser leader num term velho. Então nego o voto.
			if rv.Term < raft.currentTerm {
				reply.VoteGranted = false

			// Quem mandou quer ser o líder num novo term, então me atualizo.
			} else if rv.Term > raft.currentTerm {
				log.Printf("[FOLLOWER] Update old term '%v' to new term '%v' from '%v'.\n", raft.currentTerm, rv.Term, raft.peers[rv.CandidateID])
				raft.currentTerm = rv.Term
				raft.votedFor = 0
			}

			// Quem mandou quer ser o líder no meu term ou no novo (já atualizado
			// acima). Dou meu voto, se eu não votei ou se já votei nesse mesmo.
			if rv.Term == raft.currentTerm && (raft.votedFor == 0 || raft.votedFor == rv.CandidateID) {
				log.Printf("[FOLLOWER] Vote granted to '%v' for term '%v'.\n", raft.peers[rv.CandidateID], raft.currentTerm)
				raft.votedFor = rv.CandidateID
				raft.resetElectionTimeout()
				reply.VoteGranted = true
			} else {
				log.Printf("[FOLLOWER] Vote denied to '%v' for term '%v'.\n", raft.peers[rv.CandidateID], raft.currentTerm)
				reply.VoteGranted = false
			}

			rv.replyChan <- reply

		case ae := <-raft.appendEntryChan:

			reply := &AppendEntryReply{
				Term: raft.currentTerm,
			}

			// Quem mandou acha que é líder, mas não é. Então ‘rejeito’ essa msg.
			if ae.Term < raft.currentTerm {
				log.Printf("[FOLLOWER] Rejected AppendEntry from '%v'.\n", raft.peers[ae.LeaderID])
				reply.Success = false
				ae.replyChan <- reply
				continue

			// Quem mandou é o líder de um novo term, então me atualizo.
			} else if ae.Term > raft.currentTerm {
				raft.currentTerm = ae.Term
				raft.votedFor = 0
			}

			// Como existe um líder do meu term ou de um novo term (já atualizado
			// acima), então ‘limpo’ meu timeout
			raft.resetElectionTimeout()

			log.Printf("[FOLLOWER] Accept AppendEntry from '%v'.\n", raft.peers[ae.LeaderID])
			reply.Success = true
			ae.replyChan <- reply

		}
	}
}

// candidateSelect implements the logic to handle messages from distinct
// events when in candidate state.
func (raft *Raft) candidateSelect() {
	log.Println("[CANDIDATE] Run Logic.")
	// Candidates (§5.2):
	// Increment currentTerm, vote for self
	raft.currentTerm++
	raft.votedFor = raft.me
	voteCount := 1

	log.Printf("[CANDIDATE] Running for term '%v'.\n", raft.currentTerm)
	// Reset election timeout
	raft.resetElectionTimeout()
	// Send RequestVote RPCs to all other servers
	replyChan := make(chan *RequestVoteReply, 10*len(raft.peers))
	raft.broadcastRequestVote(replyChan)

	for {
		select {
		case <-raft.electionTick:
			// If election timeout elapses: start new election
			log.Println("[CANDIDATE] Election timeout.")
			raft.currentState.Set(candidate)
			return
		case rvr := <-replyChan:

			if rvr.VoteGranted {
				log.Printf("[CANDIDATE] Vote granted by '%v'.\n", raft.peers[rvr.peerIndex])
				voteCount++
				log.Println("[CANDIDATE] VoteCount: ", voteCount)
				// break

				// O canal replyChan recebe respostas dos votos. Basicamente tem que
				// relatar cada voto e depois fazer: “If votes received from majority of servers:
				// become leader”.
				if voteCount >= 1 + int(len(raft.peers)/2) {
					log.Printf("[CANDIDATE] I am the new elected leader.\n")
					raft.currentState.Set(leader)
					return
				}

			} else {
				log.Printf("[CANDIDATE] Vote denied by '%v'.\n", raft.peers[rvr.peerIndex])
			}

		case rv := <-raft.requestVoteChan:

			reply := &RequestVoteReply{
				Term: raft.currentTerm,
			}

			// Quem mandou quer ser leader num term velho. Então nego o voto.
			if rv.Term < raft.currentTerm {
				reply.VoteGranted = false

			// Quem mandou quer ser o líder num novo term, então me atualizo e
			// volto a ser follower.
			} else if rv.Term < raft.currentTerm {
				log.Printf("[CANDIDATE] Update old term '%v' to new term '%v' from '%v'.\n", raft.currentTerm, rv.Term, raft.peers[rv.CandidateID])
				raft.currentTerm = rv.Term
				raft.votedFor = 0
				// step down
				log.Printf("[CANDIDATE] Stepping down.\n")
				raft.currentState.Set(follower)
				raft.requestVoteChan <- rv
				return
			}

			// Quem mandou quer ser o líder no meu term. Dou meu voto, se eu não
			// votei ou se já votei nesse mesmo. Daí eu adio a minha própria eleição
			// (‘resetando’ o timeout).
			if rv.Term == raft.currentTerm && (raft.votedFor == 0 || raft.votedFor == rv.CandidateID) {
				log.Printf("[CANDIDATE] Vote granted to '%v' for term '%v'.\n", raft.peers[rv.CandidateID], raft.currentTerm)
				raft.votedFor = rv.CandidateID
				raft.resetElectionTimeout()
				reply.VoteGranted = true
			} else {
				log.Printf("[CANDIDATE] Vote denied to '%v' for term '%v'.\n", raft.peers[rv.CandidateID], raft.currentTerm)
				reply.VoteGranted = false
			}

			rv.replyChan <- reply

		case ae := <-raft.appendEntryChan:

			reply := &AppendEntryReply{
				Term: raft.currentTerm,
			}

			// Quem mandou acha que é líder, mas não é. Então ‘rejeito’ essa msg.
			if ae.Term < raft.currentTerm {
				log.Printf("[CANDIDATE] Rejected AppendEntry from '%v'.\n", raft.peers[ae.LeaderID])
				reply.Success = false
				ae.replyChan <- reply
				continue

			// Quem mandou é o líder de um novo term, então me atualizo.
			} else if ae.Term > raft.currentTerm {
				log.Printf("[CANDIDATE] Update old term '%v' to new term '%v' from '%v'.\n", raft.currentTerm, ae.Term, raft.peers[ae.LeaderID])
				raft.currentTerm = ae.Term
				raft.votedFor = 0
			}

			// Como existe um líder do meu term ou de um novo term (já atualizado
			// acima), então volto a ser follower
			// step down
			log.Printf("[CANDIDATE] Stepping down.\n")
			raft.currentState.Set(follower)
			raft.appendEntryChan <- ae
			return

		}
	}
}

// leaderSelect implements the logic to handle messages from distinct
// events when in leader state.
func (raft *Raft) leaderSelect() {
	log.Println("[LEADER] Run Logic.")
	replyChan := make(chan *AppendEntryReply, 10*len(raft.peers))
	raft.broadcastAppendEntries(replyChan)

	heartbeat := time.NewTicker(raft.broadcastInterval())
	defer heartbeat.Stop()

	broadcastTick := make(chan time.Time)
	defer close(broadcastTick)

	go func() {
		for t := range heartbeat.C {
			broadcastTick <- t
		}
	}()

	for {
		select {
		case <-broadcastTick:
			raft.broadcastAppendEntries(replyChan)
		case aet := <-replyChan:

			// O canal replyChan aqui recebe respostas dos AppendEntry
			// enviados por mim (líder). Basicamente faço o seguinte:
			// Se receber falso, é porque existe outro líder mais atual,
			// então faço Step down!
			if ! aet.Success {
				// step down
				log.Printf("[LEADER] Rejected from '%v'.\n", raft.peers[aet.peerIndex])
				log.Printf("[LEADER] Stepping down.\n")
				raft.currentState.Set(follower)
				return
			}

			_ = aet

		case rv := <-raft.requestVoteChan:

			reply := &RequestVoteReply{
				Term: raft.currentTerm,
			}

			// Quem mandou quer ser leader num term velho. Então nego o voto.
			if rv.Term < raft.currentTerm {
				reply.VoteGranted = false

			// Quem mandou quer ser o líder num novo term, então me atualizo e
			// volto a ser follower.
			} else if rv.Term < raft.currentTerm {
				log.Printf("[LEADER] Update old term '%v' to new term '%v' from '%v'.\n", raft.currentTerm, rv.Term, raft.peers[rv.CandidateID])
				raft.currentTerm = rv.Term
				raft.votedFor = 0
				// step down
				log.Printf("[LEADER] Stepping down.\n")
				raft.currentState.Set(follower)
				raft.requestVoteChan <- rv
				return
			}

			// Quem mandou quer ser o líder no meu term. Dou meu voto, se eu não
			// votei ou se já votei nesse mesmo. Daí eu adio a minha própria eleição
			// (‘resetando’ o timeout).
			if rv.Term == raft.currentTerm && (raft.votedFor == 0 || raft.votedFor == rv.CandidateID) {
				log.Printf("[LEADER] Vote granted to '%v' for term '%v'.\n", raft.peers[rv.CandidateID], raft.currentTerm)
				raft.votedFor = rv.CandidateID
				raft.resetElectionTimeout()
				reply.VoteGranted = true
			} else {
				log.Printf("[LEADER] Vote denied to '%v' for term '%v'.\n", raft.peers[rv.CandidateID], raft.currentTerm)
				reply.VoteGranted = false
			}

			rv.replyChan <- reply

		case ae := <-raft.appendEntryChan:

			reply := &AppendEntryReply{
				Term: raft.currentTerm,
			}

			// Quem mandou acha que é líder, mas não é. Então ‘rejeito’ essa msg.
			if ae.Term < raft.currentTerm {
				log.Printf("[LEADER] Rejected AppendEntry from '%v'.\n", raft.peers[ae.LeaderID])
				reply.Success = false
				ae.replyChan <- reply
				continue

			// Quem mandou é o líder de um novo term, então me atualizo.
			} else if ae.Term > raft.currentTerm {
				log.Printf("[LEADER] Update old term '%v' to new term '%v' from '%v'.\n", raft.currentTerm, ae.Term, raft.peers[ae.LeaderID])
				raft.currentTerm = ae.Term
				raft.votedFor = 0
			}

			// Como existe um líder do meu term ou de um novo term (já atualizado
			// acima), então volto a ser follower
			// step down
			log.Printf("[LEADER] Stepping down.\n")
			raft.currentState.Set(follower)
			raft.appendEntryChan <- ae
			return

		}
	}
}
