package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Due to the Go regex engine not supporting look-forwards and look-behinds, capturing groups have been used instead

const (
	numberRe          string = `^\d{1,2}\.\s`                                                         // should match "12."
	airacRe           string = `AIRAC \(\d{4}\)`                                                      // should match "AIRAC (2012)
	airacNumRe        string = `\d{4}`                                                                // should match 2012
	airacMessageRe    string = `(?:-\s)([\S\s]*)`                                                     // should match the "- message" part, but only captures the message
	contribNameRe     string = `(?:\-\sthanks\sto\s\@[A-Za-z0-9-]+\s\()([A-Za-z]+\s?[A-Za-z]*)(?:\))` // matches the thanks to part, but only captures the names
	contribEndRe      string = `\s-\sthanks\sto\s[^\n]*$`                                             // matches the thanks to part
	defaultInputFile  string = "changelog.md"
	defaultOutputFile string = "output.txt"
)

type Changelog struct {
	Changes      []string
	AIRACList    []string
	AIRACMap     map[string][]string
	AIRACs       []int
	Other        []string
	Contributors []string
}

func main() {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()
	// Get input + output  from flags
	inputfile := flag.String("in", defaultInputFile, "Set the input file")
	outputfile := flag.String("out", defaultOutputFile, "Set the output file")
	flag.Parse()
	// Open + read file
	fmt.Println("Reading from: " + *inputfile)
	filebytes, err := os.ReadFile(*inputfile)
	if err != nil {
		log.Panicln("Unable to read from " + *inputfile)
	}
	changelog := Changelog{}
	changelog.Changes = GetChanges(filebytes)
	changelog.AIRACList, changelog.Other = changelog.ChangesSorter()
	changelog.AIRACMap, changelog.AIRACs = changelog.AIRACMapGen()
	changelog.Contributors = changelog.ContribGen()
	file := CreateFile(*outputfile)
	OutputAIRAC(file, changelog)
	OutputOther(file, changelog)
	OutputContribs(file, changelog)
	fmt.Println("Output to: " + *outputfile)
	timeElapsed := time.Since(start)
	fmt.Println("Time taken:", timeElapsed)
}

func GetChanges(filebytes []byte) []string {
	changes := []string{}
	split := strings.Split(string(filebytes), "\n") // split the document into newlines
	numberReComp := regexp.MustCompile(numberRe)
	for _, s := range split {
		b := numberReComp.MatchString(s) // ensures that it isn't the AIRAC title
		if b {
			s = numberReComp.ReplaceAllString(s, "") // removes the number from the start
			changes = append(changes, s)
		}
	}
	return changes
}

func (changelog *Changelog) ChangesSorter() ([]string, []string) {
	airacReComp := regexp.MustCompile(airacRe)
	contribEndReComp := regexp.MustCompile(contribEndRe)
	airacList := []string{}
	otherList := []string{}
	for _, s := range changelog.Changes {
		contribLoc := contribEndReComp.FindStringIndex(s) // find the location of where the "thanks to" part begins
		if len(contribLoc) == 0 {
			contribLoc = []int{0}
			contribLoc[0] = len(s) // if there is no "thanks to" part then just keep the whole message
		}
		contribLocBeg := contribLoc[0]
		b := airacReComp.MatchString(s)
		if b {
			airacList = append(airacList, s[:contribLocBeg]) // remove the thanks to part
		} else {
			otherList = append(otherList, s[:contribLocBeg]) // ditto
		}
	}
	return airacList, otherList
}

func (changelog *Changelog) AIRACMapGen() (map[string][]string, []int) {
	airacmap := map[string][]string{}
	airacNumReComp := regexp.MustCompile(airacNumRe)
	airacMessageReComp := regexp.MustCompile(airacMessageRe)
	for _, s := range changelog.AIRACList {
		num := airacNumReComp.FindString(s)
		message := airacMessageReComp.FindStringSubmatch(s)
		if len(message) == 0 { // if it can't find the message, something has gone wrong
			fmt.Println("Malformed message string in " + s)
			continue
		}
		airacval := airacmap[num]
		airacval = append(airacval, message[1])
		airacmap[num] = airacval
	}
	keylist := []int{}
	for d := range airacmap {
		d, err := strconv.Atoi(d) // converts to int as that makes the sorting easier
		if err != nil {
			log.Panicln("Unable to convert AIRAC string to int")
			return nil, nil
		}
		keylist = append(keylist, d)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keylist))) // reverses the sort order so that it goes newest airac first
	return airacmap, keylist
}

func (changelog *Changelog) ContribGen() []string {
	contribNameReComp := regexp.MustCompile(contribNameRe)
	contribmap := map[string]bool{}
	contriblist := []string{}
	for _, y := range changelog.Changes {
		submatch := contribNameReComp.FindStringSubmatch(y)
		if len(submatch) == 0 {
			continue
		}
		if _, ok := contribmap[submatch[1]]; !ok { // checking if the contributor already exists using a map - if they don't, then it sets the key as true and continues
			contribmap[submatch[1]] = true
			contriblist = append(contriblist, submatch[1])
		}
	}
	return contriblist
}

func CreateFile(outputFile string) *os.File {
	newfile, err := os.Create(outputFile)
	if err != nil {
		log.Panicln("Could not create new file")
	}
	return newfile
}

func OutputAIRAC(f *os.File, c Changelog) {
	f.WriteString("--- AIRACs: ---" + "\n")
	for _, key := range c.AIRACs {
		value := c.AIRACMap[fmt.Sprint(key)]
		f.WriteString(fmt.Sprint(key) + ":\n")
		for _, y := range value {
			f.WriteString(y + "\n")
		}
		f.WriteString("\n")
	}
}

func OutputOther(f *os.File, c Changelog) {
	f.WriteString("--- Other: ---" + "\n")
	for _, value := range c.Other {
		f.WriteString(value + "\n")
	}
	f.WriteString("\n")
}

func OutputContribs(f *os.File, c Changelog) {
	f.WriteString("--- Contributors: ---" + "\n")
	for _, value := range c.Contributors {
		f.WriteString(value + "\n")
	}
}
