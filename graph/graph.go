package graph

import "sort"

// WordGraph is a representation of a body of text as a graph of words. Each
// node contains a string key which maps to a map of all the words that follow
// immediately after, and their counts.
type Children map[string]int
type WordGraph map[string]Children

func CreateWordGraphFromInputText(inputText []string) WordGraph {
	wg := make(WordGraph)

	// Store each occurrence of a word in the text into a graph.
	for _, word := range inputText {
		if _, present := wg[word]; !present {
			wg[word] = make(Children)
		}
	}

	// For each word found, increment its count seen immediately following
	// lastWord.
	lastWord := ""
	for _, word := range inputText {
		if lastWord != "" {
			wg[lastWord][word] += 1
		}
		lastWord = word
	}
	return wg
}

type MarkovNode struct {
	Key    string
	Weight float64
}
type MarkovNodeList []MarkovNode
type MarkovChain map[string][]MarkovNode

func (b MarkovNodeList) Len() int {
	return len(b)
}

func (b MarkovNodeList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b MarkovNodeList) Less(i, j int) bool {
	return b[i].Weight < b[j].Weight
}

func CreateMarkovChainFromWordGraph(wg WordGraph) MarkovChain {
	mc := make(MarkovChain)
	for word, children := range wg {
		mnl := make(MarkovNodeList, 0)
		totalCount := 0
		for _, childCount := range children {
			totalCount += childCount
		}
		for childKey, childCount := range children {
			mnl = append(mnl, MarkovNode{
				Key:    childKey,
				Weight: float64(childCount) / float64(totalCount) * 100,
			})
		}
		sort.Sort(sort.Reverse(mnl))
		lastWeight := 0.0
		for i, _ := range mnl {
			mnl[i].Weight += lastWeight
			lastWeight = mnl[i].Weight
		}
		mc[word] = mnl
	}
	return mc
}
