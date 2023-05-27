package main

import (
	"net"
	"log"
	"bufio"
	"strings"

	"fmt" // print statements

	"github.com/go-python/cpy3" // creates bindings between python and go (github.com/go-python/cpy3)

	"path/filepath" // used to create directories (path/filepath)

	"os" // used to get the cli arguments (os)

	"time" // used to track training time

	"main/myDfa" // pakcage for custom DFA
)

const (
    HOST = "localhost"
    PORT = "8000"
    TYPE = "tcp"
)

// Uncomment the desired DFA to be used for learning 

// var true_DFA = myDfa.FromDfaGo("AbbadingoDFAs/eight_State.json")
// var true_DFA = myDfa.FromDfaGo("AbbadingoDFAs/sixteen_State.json")
var true_DFA = myDfa.FromDfaGo("AbbadingoDFAs/thirtyTwo_State.json")
// var true_DFA = myDfa.FromDfaGo("AbbadingoDFAs/sixtyFour_State.json")

// membership query function
func mq(s string) bool {
	ret := true_DFA.Label(s)	
	return ret
}

// equivalence query function
func eq(dfa myDfa.Dfa) string {
	ret := true_DFA.BFS_ce(dfa)	
	return ret
}

// starts a tcp server and use the given integer parameter
// as the size of the buffer
func start_tcp(bufferSize int){
	listener, err := net.Listen(TYPE, ":" + PORT)

	if err != nil {
		log.Fatal(err)	
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		
		if err != nil {
			log.Fatal(err)
		}

		go handlerByte(conn, bufferSize)
	}
}

// handler function that works by utilising a byte array
// for the buffer 
func handlerByte(conn net.Conn, bufferSize int) {

	for {
		buffer := make([]byte, bufferSize)
		bufLen, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Read returns error...")
			log.Fatal(err)
		}

		query := string(buffer[:bufLen])

		if query[:2] == "mq" {
			word := query[3:]

			if mq(word) == true {
				conn.Write([]byte("true"))
			} else {
				conn.Write([]byte("false"))
			}
		} else if query[:2] == "eq" {
			hypo_dfa := myDfa.JSONtoDFA(query[3:])
			ret := eq(hypo_dfa)
			conn.Write([]byte(ret))
		} else if query[:3] == "end" {
			conn.Close()
			break
		} else {
			conn.Write([]byte("Not a query"))
		}
	}

	conn.Close()
}

func handler(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
		}

		query := strings.Split(message, "q ")
		fmt.Println("Query as Slice:", query, "")

		if query[0] == "m" {
			word := query[1]

			if mq(word) == true {
				conn.Write([]byte("true"))
			} else {
				conn.Write([]byte("false"))
			}
		} else if query[0] == "e" {
			hypo_dfa := myDfa.JSONtoDFA(query[1])
			ret := eq(hypo_dfa)
			conn.Write([]byte(ret))
		} else if query[0] == "end" {
			conn.Close()
			break
		} else {
			conn.Write([]byte("Not a query"))
		}

		// conn.Close()
	}

}

func main(){
	// 2048 is enough for DFAs up to 32 states, 8192 is required for 64 states
	go start_tcp(2048)

	// finalize and initializer a python interpreter
	defer python3.Py_Finalize()
	python3.Py_Initialize()

	// getting the current path cwd 
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// setting the sys variable in python to the current directory to import the python module
	python3.PyRun_SimpleString("import random\nimport re\nimport time\nimport sys\nsys.path.append(\"" + dir + "\")")
	PyFolder := python3.PyImport_ImportModule("PyFolder")

	// checking to make sure the module import was succesfull
	if PyFolder == nil {
		panic("PyFolder  not found!")
	}

	// // get a dictionary with all atributes in the module
	oDict := python3.PyModule_GetDict(PyFolder)

	// // Gathering all the functions required for lstar
	init_func := python3.PyDict_GetItemString(oDict, "init_lstar")
	learn_func := python3.PyDict_GetItemString(oDict, "learn")

	// create a tuple with 2 len
	// create python bytes objects from the string of
	// the ip and port required for python to ask queries
	pyTuple_s3 := python3.PyTuple_New(3)
	ip_py := python3.PyBytes_FromString(HOST)
	port_py := python3.PyBytes_FromString(PORT)
	bufferSize_py := python3.PyBytes_FromString("8192")

	// set these python objects into the tuples
	python3.PyTuple_SetItem(pyTuple_s3, 0, ip_py)
	python3.PyTuple_SetItem(pyTuple_s3, 1, port_py)
	python3.PyTuple_SetItem(pyTuple_s3, 2, bufferSize_py)

	// get the lstar object return
	py_lstar := init_func.CallObject(pyTuple_s3)

	// create a new tuple with 1 len
	// and place the lstar object inside of it
	pyTuple_s1 := python3.PyTuple_New(1)
	python3.PyTuple_SetItem(pyTuple_s1, 0, py_lstar)	
	
	fmt.Println("---------------------------------")

	// start recording time
	start := time.Now()

	// call learn func and pass the pytupple containing the lstar object
	learn_ret := learn_func.CallObject(pyTuple_s1)

	// stop recording time and print
	taken := time.Since(start)

	JSON_String := python3.PyBytes_AsString(learn_ret)
	learnt_dfa := myDfa.JSONtoDFA(JSON_String)	

	fmt.Println("\nTime Taken:", taken.Seconds())
	fmt.Println("---------------------------------")
	fmt.Println("Number of States in DFA Learnt:", len(learnt_dfa.States))
	fmt.Println("---------------------------------")

	// clearing up memory
	defer learn_ret.DecRef()
	defer port_py.DecRef()
	defer ip_py.DecRef()
	defer pyTuple_s1.DecRef()
	defer pyTuple_s3.DecRef()
	defer learn_func.DecRef()
	defer init_func.DecRef()
	defer PyFolder.DecRef()
}