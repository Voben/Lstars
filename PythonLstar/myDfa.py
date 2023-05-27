import json

# will extract a DFA from a JSON file containing a DfaGo DFA
def from_DfaGo_file(filename):
    with open(filename) as f:
        d = json.load(f)
    
    ret_DFA = DFA()

    for i in d["Alphabet"]:
        ret_DFA.alphabet.append(str(i))

    ret_DFA.start = str(d["StartingState"])

    for idx1 in range(0, len(d["States"])):
        ret_DFA.states.append(str(idx1))

        if d["States"][idx1]["Label"] == 0:
            ret_DFA.accs.append(str(idx1))

        for idx2 in range(0, len(d["States"][idx1]["Next"])):
            if str(idx1) not in ret_DFA.tf:
                ret_DFA.tf[str(idx1)] = {}
            ret_DFA.tf[str(idx1)][str(idx2)] = str(d["States"][idx1]["Next"][idx2])
            
    ret_DFA.states = ret_DFA.states[1:]       

    return ret_DFA   

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

    # this function expects a string to be processed by the DFA
    # returns a boolean depending on if the word is part of the
    # DFA's language or not
    def label(self, string):
        curNode = self.start

        for char in string:
            
            if curNode not in self.tf:
                return False
            if char not in self.tf[curNode]:
                return False

            # curString = self.tf[curString][char]
            curNode = self.tf[curNode][char]

        if curNode in self.accs:
            return True
        else:
            return False

    # export a Dfa from Python to a dot file
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
    
    # export a JSON tring 
    def export_JSON(self):
        jsonString = {"states": self.states, "alphabet": self.alphabet, "start": self.start, "tf": self.tf, "accepting": self.accs}
        return jsonString

    # export a DFA as a DfaGo JSON
    def export_DfaGo_JSON(self):
        ret = "{\"Alphabet\":["

        for a in self.alphabet:
            ret += a + ","

        ret = ret[:-1] + "],\"StartingState\":" + str(self.states.index(self.start)) + ",\"States\":["

        for s in self.states:
            curNodeStr = "{\"Label\":"
            
            if s in self.accs:
                curNodeStr += "0,\"Next\":["
            else:
                curNodeStr += "2,\"Next\":["

            for v in self.tf[s].values():
                curNodeStr += str(self.states.index(v)) + ","

            curNodeStr = curNodeStr[:-1] + "],\"depth\":0,\"order\":0},"
        
            ret += curNodeStr

        ret = ret[:-1] + "],\"dateCreated\":\"None\",\"depth\":0,\"dirty\":false,\"docType\":\"DfaGo/DFA\",\"version\":1}"
        
        print("inside of export")
        print(ret)

        return ret

    # Breadth First Search Counter-Exampling Method
    def BFS_ce(self, h_DFA):
        visited, queue = {}, [["", self.start]]

        if self.label("") != h_DFA.label(""):
            return ""

        while queue:
            vertex_lst = queue.pop(0)
            
            if vertex_lst[1] not in visited:
                visited[vertex_lst[1]] = 1
            else:
                visited[vertex_lst[1]] += 1


            if visited[vertex_lst[1]] <= 2:

                for sym, next_state in self.tf[vertex_lst[1]].items():
    
                    if self.label(vertex_lst[0] + sym) != h_DFA.label(vertex_lst[0] + sym):
                        return vertex_lst[0] + sym

                    queue.append([vertex_lst[0] + sym, next_state])
        
        return "None"