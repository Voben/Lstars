from lstar import Lstar
import timeit
import myDfa

# import DfaGo DFAs from JSON files

dfa = myDfa.from_DfaGo_file("AbbadingoDFAs/eight_State.json")
# dfa = myDfa.from_DfaGo_file("AbbadingoDFAs/sixteen_State.json")
# dfa = myDfa.from_DfaGo_file("AbbadingoDFAs/thirtyTwo_State.json")
# dfa = myDfa.from_DfaGo_file("AbbadingoDFAs/sixtyFour_State.json")

# query functions to be passed as a parameter to the lstar algorithm 
def mq(string):    
    return  dfa.label(string)

def eq(hDFA):    
    return  dfa.BFS_ce(hDFA)

print("----------------------------")

# create Lstar object
lstar = Lstar()

# record the start time
start_time = timeit.default_timer()

# start learning
ret = lstar.learn_DFA(mq, eq, ['0', '1'])

# record the end of training time and calculate the training time 
end_time = timeit.default_timer()
training_time = end_time - start_time

# Print the training time taken
print("\nTraining Time: "+"{:.10f}".format(training_time) + " seconds")

print("----------------------------")

# Print the learnt DFA
print("Number of Iterations:", lstar.iter, '\n')

print("Number of States:", len(ret.states))
print("All States:", ret.states)

print("Alphabet:", ret.alphabet)

print("Starting State: '", ret.start,"'", sep='')

print("All Accepting States:", ret.accs)
print("Number of Accepting States:", len(ret.accs), '\n----------------------------')