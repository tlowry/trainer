// TrainBayesian project main.go
package main

import (
	"bufio"
	"github.com/jbrukh/bayesian"
	"log"
	"os"
	"strings"
)

const (
	Wanted   bayesian.Class = "Wanted"
	Unwanted bayesian.Class = "UnWanted"
)

var wanted = []string{""}
var unWanted = []string{""}

func LearnFile(classifier *bayesian.Classifier, name string, class bayesian.Class) {
	file, err := os.OpenFile(name, os.O_RDONLY, 0666)
	if err != nil {
		panic("could not open file")
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		words := strings.Split(string(line), " ")
		classifier.Learn(words, class)
	}
}

func testClassifier(c *bayesian.Classifier, doc []string) bayesian.Class {

	_, inx, _ := c.ProbScores(doc)
	class := c.Classes[inx]

	return class
}

func TestFile(classifier *bayesian.Classifier, name string, class bayesian.Class) int {
	file, err := os.OpenFile(name, os.O_RDONLY, 0666)
	if err != nil {
		panic("could not open file")
	} else {
		log.Println("Parsing file ", name, "as", class)
	}
	var score, tested int
	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		words := strings.Split(string(line), " ")
		cls := testClassifier(classifier, words)

		if cls == class {
			score++
		} else {
			// wrong, try again
			log.Println("Classifier wrong ", words, "!=", cls, " expected ", class)
			classifier.Learn(words, class)
		}

		tested++
	}

	pc := (score * 100) / tested
	return pc
}

func main() {

	c := bayesian.NewClassifier(Wanted, Unwanted)

	var err error = nil
	/*
		c, err := bayesian.NewClassifierFromFile("files/gfx.ebay.classifier")
		if err != nil {
			panic(err)
		}
	*/

	wantedFile := "files/wanted.txt"
	LearnFile(c, wantedFile, Wanted)

	unWantedFile := "files/unwanted.txt"
	LearnFile(c, unWantedFile, Unwanted)

	log.Printf("classifier is trained: %d documents read\n", c.WordCount())

	accuracy := 0
	for accuracy < 100 {
		accuracy = TestFile(c, wantedFile, Wanted)
	}

	accuracy = 0
	for accuracy < 100 {
		accuracy = TestFile(c, unWantedFile, Unwanted)
	}

	err = c.WriteToFile("files/gfx.ebay.classifier")
	if err != nil {
		panic(err)
	}

}
