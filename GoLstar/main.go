package main

import (
	"lstar/lstarPack"
	"fmt"
	"time"
)

// this is the learning main
func main() {
	// true_DFA := lstarPack.FromDfaGo("AbbadingoDFAs/eight_State.json")
	// true_DFA := lstarPack.FromDfaGo("AbbadingoDFAs/sixteen_State.json")
	// true_DFA := lstarPack.FromDfaGo("AbbadingoDFAs/thirtyTwo_State.json")
	true_DFA := lstarPack.FromDfaGo("AbbadingoDFAs/sixtyFour_State.json")

	// variables containg query functions
	member := func(s string) bool {
		return true_DFA.Label(s)
	}

	equiv := func(dfa lstarPack.Dfa) string {
		return true_DFA.BFS_ce(dfa)
	}

	// alphabet of the language
	a := [2]string{"0", "1"}

	// initiating a learning function
	learner := lstarPack.NewLstar(a, member, equiv)

	fmt.Println("----------------------------")

	// recording training time
	start := time.Now()
	ret_DFA := learner.Learn_dfa()
	taken := time.Since(start)

	fmt.Printf("\nTime Taken in Seconds (8 d.p.): %0.8f\n", taken.Seconds())
	fmt.Println("----------------------------")
	fmt.Println("Learnt DFA")
	fmt.Println("States:", len(ret_DFA.States))
	fmt.Println("Alphabet:", ret_DFA.Alphabet)
	fmt.Println("Start:", ret_DFA.Start)
	fmt.Println("Final states:", ret_DFA.Accepting)
	fmt.Println("----------------------------")
}