package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/eshyong/markovchain/graph"
	"github.com/eshyong/markovchain/parse"
)

func main() {
	debug := flag.Bool("debug", false, "Set this to print out debug information, including a JSON representation of the Markov chain. By default, that file will be written to \"markovchain.json\"")
	flag.Parse()
	if flag.Arg(0) == "" {
		fmt.Println("Usage: markovchain filename [runlength")
		os.Exit(0)
	}

	fileContentsAsBytes, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	fileContentsAsString := string(fileContentsAsBytes)
	parsedStrings := parse.ParseInputString(fileContentsAsString)

	// Perform processing on the input text to convert it into a Markov chain.
	wordGraph := graph.CreateWordGraphFromInputText(parsedStrings)
	markovChain := graph.CreateMarkovChainFromWordGraph(wordGraph)

	if *debug {
		outputJson, err := json.MarshalIndent(markovChain, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("markovchain.json", outputJson, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Press enter to print the next string: ")

	rand.Seed(time.Now().Unix())
	nextString := parsedStrings[rand.Intn(len(parsedStrings))]
	stdin := bufio.NewReader(os.Stdin)
	var nextWords []graph.MarkovNode
	for {
		// ReadLine here is only used to prompt the next string.
		_, _, err := stdin.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		// Print the currently selected string, then select the next string by
		// rolling a dice and picking based on the weights of the nodes in the
		// list.
		fmt.Printf("%s ", nextString)
		nextWords = markovChain[nextString]
		for _, node := range nextWords {
			if weight := float64(rand.Intn(100)); weight < node.Weight {
				nextString = node.Key
				break
			}
		}
	}
}
