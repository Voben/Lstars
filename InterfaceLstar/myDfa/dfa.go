package myDfa

import (
	"encoding/json"
	"strings"
	"strconv"
	"os"
	"io/ioutil"
	"fmt"
)

// these 2 data types are obtained from DfaGo and are used 
// to obtain a custom DFA object from a DfaGo JSON 
type State struct {
	Label int8 // Label of the state.
	Next  []int      // Index is the alphabet symbol. Value is index of adjacent state (-1 if there is no transition for symbol).

	// Internal.
	depth int // The depth of the state in the DFA (only if DFA is not dirty).
	order int // The canonical order of the state in the DFA (only if DFA is not dirty).
}

type DfaGo struct {
	States        []State // The states in the DFA.
	Alphabet      []int   // A slice containing the ordered alphabet of the DFA.
	StartingState int     // The index of the starting state of the DFA.
	depth         int     // The depth of the DFA.
	dirty         bool    // Flag to tell us whether the depth and order of the DFA and its states is computed.
}

// this is the custom DFA struct
type Dfa struct {
	States []string
	Alphabet [2]string
	Start string
	Tf map[string]map[string]string
	Accepting []string
}

// obtain a DFA from a DfaGo JSON
func FromDfaGo(filename string) Dfa {
	jsonFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var myOwn Dfa
	var fromJSON DfaGo

	json.Unmarshal(byteValue, &fromJSON)

	myOwn.Tf = map[string]map[string]string{}
	myOwn.Start = strconv.Itoa(fromJSON.StartingState)

	for idx, char := range fromJSON.Alphabet{
		myOwn.Alphabet[idx] = strconv.Itoa(char)
	}

	for idx, state := range fromJSON.States{
		myOwn.States = append(myOwn.States, strconv.Itoa(idx))

		if state.Label == 0 {
			myOwn.Accepting = append(myOwn.Accepting, strconv.Itoa(idx))
		}

		myOwn.Tf[strconv.Itoa(idx)] = map[string]string{}

		for idx2, next := range state.Next{
			myOwn.Tf[strconv.Itoa(idx)][strconv.Itoa(idx2)] = strconv.Itoa(next)
		}
	}

	return myOwn 
}

// change a Python list saved as a string to a slice of strings
func listToSlice(s string) []string{
	if strings.Contains(s, "'") == false{
		var nilSlice []string
		return nilSlice
	} else {
		s = strings.ReplaceAll(s, "[", "")
		s = strings.ReplaceAll(s, "]", "")
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "'", "")
		return strings.Split(s, ",")
	}
}

// change alphabet list to slice
func alphaListToSlice(s string) [2]string{
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, " ", "")
	s_slice := strings.Split(s, ",")

	arr := [2]string{}
	copy(arr[:], s_slice[:1])

	return arr
}

// change a string of a transition function (tf)
// to a nested map acting as a tf
func strToTF(s string) map[string]map[string]string {
	ret_Tf := map[string]map[string]string{}
	var last_key string


	s = strings.ReplaceAll(s, "}", "")
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, ",", ":")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, " ", "")

	ret := strings.Split(s, ":")

	// for idx, ele := range ret{
	for i := 0; i < len(ret) - 1; i++{
		if (i % 5) == 0 || i == 0{
			ret_Tf[ret[i]] = map[string]string{}
			last_key = ret[i]
		} else {
			ret_Tf[last_key][ret[i]] = ret[i+1]
			i++
			i++
			ret_Tf[last_key][ret[i]] = ret[i+1]
			i++
		}
	}

	return ret_Tf
}

// check if a element contains a slice or not
func slice_contains(s []string, ele string) bool{
	for _, element := range s{
		if ele == element{
			return true
		}
	}
	return false
}

// check if a nested map contains a string
func nest_map_contains(s map[string]map[string]string, ele string) bool{
	_, val := s[ele]
	return val 
}

// check if a map contains a string
func map_contains(s map[string]string, ele string) bool{
	_, val := s[ele]
	return val 
}

// change a JSON string to a DFA
func JSONtoDFA(s string) Dfa{
	var x map[string]string
	json.Unmarshal([]byte(s), &x)

	alphabet := [2]string{"0", "1"}

	return Dfa{States: listToSlice(x["states"]), Alphabet: /*alphaListToSlice(x["alphabet"])*/alphabet, Start: "", Tf: strToTF(x["tf"]), Accepting: listToSlice(x["accepting"])}
}

// use a DFA to process a given word
func (dfa Dfa) Label(s string) bool {
	curNode := dfa.Start

	for _, ele := range s {
		if nest_map_contains(dfa.Tf, curNode) == false{
			return false
		} else if map_contains(dfa.Tf[curNode], string(ele)) == false{
			return false
		}

		curNode = dfa.Tf[curNode][string(ele)]
	}

	if slice_contains(dfa.Accepting, curNode) == true{
		return true
	} else {
		return false
	}
}

// Breadth First Search Counter-Exampling method
func (dfa Dfa) BFS_ce(h_dfa Dfa) string{
	visited := map[string]int{}

	queue := [][]string{
		{"", dfa.Start},
	}

	vertex_slice := []string{}

	for ;;{
		if len(queue) <= 0 {
			return "None"
		}

		vertex_slice, queue = queue[0], queue[1:]

		val, exists := visited[vertex_slice[1]]

		if exists {
			visited[vertex_slice[1]] ++
		} else {
			visited[vertex_slice[1]] = 1
		}

		if val <= 2{
			for sym, next_state := range(dfa.Tf[vertex_slice[1]]){
				if dfa.Label(vertex_slice[0] + sym) != h_dfa.Label(vertex_slice[0] + sym){
					return vertex_slice[0] + sym
				}

				queue = append(queue, []string{vertex_slice[0] + sym, next_state})
			}
		}
	}
}