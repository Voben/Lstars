import json

# DFA class
class DFA:
    # init class that does not require any parameters
    def __init__(self, states = [], alphabet = [], start = "", tf = {}, accs = []):
        self.add_states(states)
        self.add_alphabet(alphabet)
        self.add_start(start)
        self.add_tf(tf)
        self.add_accepting(accs)

    # the following functions all populate the different attributes of the class
    # they are used in the init func
    def add_states(self, states):
        self.states = states

    def add_alphabet(self, alphabet):
        self.alphabet = alphabet

    def add_start(self, start):
        self.start = start

        if self.start not in self.states:
            self.states.append(self.start)

    def add_tf(self, tf):
        self.tf = tf

    def add_accepting(self, accs):
        self.accs = accs

        for acc in self.accs:
            if acc not in self.states:
                self.states.append(acc)

    # output whether a string is in the DFA's language or not
    def label(self, string):
        curString = ""

        for char in string:
            
            if curString not in self.tf:
                return False
            if char not in self.tf[curString]:
                return False

            curString = self.tf[curString][char]

        if curString in self.accs:
            return True
        else:
            return False

    # export the DFA in a dot file
    def export_dot(self, string):
        state_names = {}
        for state in self.states:
            state_names[state] = str(len(state_names) + 1)

        with open(string + ".dot", "w") as fp:
            fp.write("digraph G{\n")
            fp.write("0 [label=\"\", shape=point];\n")
            fp.write("0 -> 1;\n")

            for state in self.states:

                if state in self.accs:
                    nodeShape = "doublecircle"
                else:
                    nodeShape = "circle"
    
                fp.write(state_names[state] + " [label=\""+ state_names[state] + "\", shape=" + nodeShape + "];\n")

                if state in self.tf:
                    fp.write(state_names[state] + " -> " + state_names[self.tf[state][self.alphabet[0]]] + " [label=" + self.alphabet[0] +"];\n")
                    fp.write(state_names[state] + " -> " + state_names[self.tf[state][self.alphabet[1]]] + " [label=" + self.alphabet[1] +"];\n")
                else:
                    fp.write(state_names[state] + " -> " + state_names[state] + " [label=" + self.alphabet[0] +"];\n")
                    fp.write(state_names[state] + " -> " + state_names[state] + " [label=" + self.alphabet[1] +"];\n")

            fp.write("}")

    # export the DFA as a JSON string
    def export_JSON(self):
        DFA_dict = {"states": str(self.states), "alphabet": str(self.alphabet), "start": str(self.start), "tf": str(self.tf), "accepting": str(self.accs)}
       
        jsonString = json.dumps(DFA_dict)
        
        return jsonString