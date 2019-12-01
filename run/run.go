package main

import (
	"TerminationDetection/rana"
	// "TerminationDetection/scholten"
	"flag"
	"hash/fnv"
	"log"
	"math/rand"
	"time"
)

// Command line parameters
var (
	instanceID = flag.Int("id", 1, "ID of the instance")
	seed       = flag.String("seed", "", "Seed for RNG")
	//RNG = Random Number Generator
	//parameters
	activeStart = flag.Bool("active", false, "True if and only if the process should start itself active")
	initial     = flag.Bool("init", false, "True if the process should start itself active and start both the detection and the basic algorithm")
	algorithm   = flag.String("alg", "rana", "Name of the algorithm. Can be either \"rana\" or \"scholten\"")
)

func main() {
	flag.Parse()

	if *seed != "" {
		h := fnv.New32a()
		h.Write([]byte(*seed))
		rand.Seed(int64(h.Sum32()))
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	peers := make(map[int]string)
	peers[1] = "localhost:3001"
	peers[2] = "localhost:3002"
	peers[3] = "localhost:3003"
	peers[4] = "localhost:3004"

	if _, ok := peers[*instanceID]; !ok {
		log.Fatalf("[MAIN] Invalid instance id.\n")
	}

	if *algorithm == "rana" {
		rana := rana.NewRana(peers, *instanceID, *initial)
		<-rana.Done()
	// } else if *algorithm == "scholten" {
	// 	scholten := scholten.newScholten(peers, *instanceID)
	// 	<-scholten.Done()
	} else {
		log.Fatalf("[MAIN] Invalid algorithm")
	}
}
