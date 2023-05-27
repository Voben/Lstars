package lstarPack

import (
	"fmt"
)

// lstar struct
type lstar struct{
	alphabet [2]string
	mq mem_query
	eq equiv_query
	num_mq int
	num_eq int
}

// query type
type mem_query func(string) bool

type equiv_query func(Dfa) string

// create a new lstar object
func NewLstar(alpha [2]string, in_mq mem_query, in_eq equiv_query) lstar{
	return lstar{alphabet: alpha, mq: in_mq, eq: in_eq}
}

// learn function from DFA
func (l lstar) Learn_dfa() Dfa{
	l.num_eq = 0
	l.num_mq = 0
	Q := []string{""}
	T := []string{""}
	tf := map[string]map[string]string{}
	iter := 1

	for {
		Q, tf = l.close(Q, T)

		final_states := []string{}
		for _, state := range Q {
			l.num_mq = l.num_mq + 1
			if l.mq(state) == true {
				final_states = append(final_states, state)
			}
		}

		// create hyopthesis dfa
		h_DFA := Dfa{States: Q, Alphabet: l.alphabet, Start: "", Tf: tf, Accepting: final_states}

		l.num_eq = l.num_eq + 1
		ret := l.eq(h_DFA)

		if ret == "None"{
			fmt.Println("Number of Membership Queries:", l.num_mq)
			fmt.Println("Number of Equivalence Queries:", l.num_eq)
			return h_DFA
		}

		new_state, new_test_word := l.add_test_word(tf, ret)

		if slice_contains(Q, new_state) == false{
			Q = append(Q, new_state)
		}

		if slice_contains(T, new_test_word) == false{
			T = append(T, new_test_word)
		}

		iter++
	}
}

// check if two strings are indistinguishable given a set
func (l lstar) are_indistinguishable(t_set []string, u string, v string) bool{
	for _, t := range t_set{
		l.num_mq = l.num_mq + 2
		if l.mq(u + t) != l.mq(v + t){
			return false
		}
	}
	return true
}

// close the two provided sets of strings
func (l lstar) close(q_set []string, t_set []string) ([]string, map[string]map[string] string){
	tf := map[string]map[string]string{}
    q_copy := make([]string, len(q_set))
	copy(q_copy, q_set)
	
	for i := 0; i < len(q_set); i++ {

		for _, a := range l.alphabet {
			
			q := q_set[i]

			for _, r := range q_set {

				if l.are_indistinguishable(t_set, (q + a), r) == true{

					if nest_map_contains(tf, q) == false {
						tf[q] = map[string]string{}
					} 
					
					tf[q][a] = r
					break
				}	
			}

			if nest_map_contains(tf, q) == false {
				tf[q] = map[string]string{}
			}

			if map_contains(tf[q], a) == false{
				q_copy = append(q_copy, q+a)
				tf[q][a] = q+a
			}
		}
	}

	return q_copy, tf
}

// create a new test word from a given counter-exampling
func (l lstar) add_test_word(tf map[string]map[string] string, ce string) (string, string){
	q := ""
	i := 0

	for {
		ele, exists := tf[q][string(ce[i])]
		l.num_mq = l.num_mq + 2

		if exists {
			
			if l.mq(ele + ce[i+1:]) != l.mq(ce){
				return (q + string(ce[i])), ce[i+1:]				
			}

		} else {
			return (q + string(ce[i])), ce[i+1:]
		}

		q = ele
		i++
	}
}

