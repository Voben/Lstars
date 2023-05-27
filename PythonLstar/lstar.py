# custom made DFA class 
from myDfa import DFA

class Lstar:
    # will learn a DFA based on the examples and regular expression provided
    def learn_DFA(self, membership_query, equivalence_query, alphabet = ['0', '1']):
        self.num_mq = 0
        self.num_eq = 0
        self.alphabet = alphabet
        self.membership_query = membership_query
        self.equivalence_query = equivalence_query
        self.iter = 0

        qQ = [""]
        tT = [""]

        while True:
            # expand states and transition function 
            qQ, tf = self.close(qQ, tT)

            # get final states in q
            final_states = []
            for pos in qQ:
                self.num_mq += 1
                if self.membership_query(pos):
                    final_states.append(pos)

            # create hypothesis dfa
            self.hDFA = DFA(qQ, self.alphabet, "", tf, final_states)

            self.num_eq += 1
            # determine if dfa properly classifies the examples provided
            res = self.equivalence_query(self.hDFA)

            # if so return the hypothesis dfa
            if res == "None": # if res True then hDFA classifies L (or the examples provided)
                print("Number of Membership Queries:", self.num_mq)
                print("Number of Equivalence Queries:", self.num_eq)
                return self.hDFA
            
            # consider the counter example and create a new state and test word to be considered
            new_state, new_test_word = self.add_test_word(qQ, tf, res) # res should contain a counterexample that H does not classify
            
            if new_state not in qQ:
                qQ.append(new_state)

            if new_test_word not in tT:
                tT.append(new_test_word)

            self.iter += 1


    # function to check for t equivalence
    def are_indistinguishable(self, T, u, v):
        # go through the suffixes stored in T
        for t in T:
            self.num_mq += 2
            if self.membership_query( (u+t) ) != self.membership_query( (v+t)):
                return False

        return True

    # main loop sequence in l star
    def close(self, qQ, tT):
        # create a new empty tf, i to go through all states, create a copy of states to avoid infinite loop
        tf = {}
        i = 0
        newQ = qQ.copy()

        # go through all states
        while i < len(qQ):
            
            # go through the alphabet
            for a in self.alphabet:
                
                # get a new state
                q = qQ[i]

                # go through all states again
                for r in qQ:

                    # check if a state and a character in the alphabet lead to another state meaning that these two states are T-equivalent
                    if self.are_indistinguishable(tT, (q + a), r):
                        
                        # avoid key error if this is the first instance of q in tf
                        if q not in tf:
                            tf[q] = {}

                        # link the two states in the tf
                        tf[q][a] = r
                        break

                if q not in tf:
                    tf[q] = {}

                if a not in tf[q]:
                    newQ.append(q + a)
                    tf[q][a] = q + a

            i += 1
        return newQ, tf

    # used to treat a new counter example correctly
    # will go through the transition function accordingly depending on the counter example provided
    def add_test_word(self, Q, tf, ce):
        # q will store the string as we go through it and i will store the index
        q = ""
        i = 0

        # infinite loop
        while True:

            # try except that will trigger if the state is not yet in tf, same thing is carried out if triggered
            # check if the transition from the current state and the rest of string leads to the counter example
            try:
                self.num_mq += 2
                if self.membership_query( (tf[q][ce[i]] + ce[i+1:]) ) != self.membership_query(ce):
                    return (q + ce[i]), ce[i+1:]
                    
            except KeyError:
                return (q + ce[i]), ce[i+1:]

            # if not store q as the current part of the string and update the i
            else:
                q = tf[q][ce[i]]
                i += 1
