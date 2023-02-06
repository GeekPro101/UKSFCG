package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

// Due to the Go regex engine not supporting look-forwards and look-behinds, capturing groups have been used instead

const (
	numberRe       string = `^\d{1,2}\.\s`                                                         // should match "12."
	airacRe        string = `AIRAC \(\d{4}\)`                                                      // should match "AIRAC (2012)
	airacNumRe     string = `\d{4}`                                                                // should match 2012
	airacMessageRe string = `(?:-\s)([\S\s]*)`                                                     // should match the "- message" part, but only captures the message
	contribNameRe  string = `(?:\-\sthanks\sto\s\@[A-Za-z0-9-]+\s\()([A-Za-z]+\s?[A-Za-z]*)(?:\))` // matches the thanks to part, but only captures the names
	contribEndRe   string = `\s-\sthanks\sto\s[^\n]*$`                                             // matches the thanks to part
)

var opts struct {
	InputFile  string `long:"input" description:"Input file name" optional:"yes" default:"changelog.md"`
	OutputFile string `long:"output" description:"Output file name" optional:"yes" default:"output.txt"`
	Url        string `long:"url" description:"Set url - specify to use default, or provide value to use that" optional:"yes" optional-value:"https://raw.githubusercontent.com/VATSIM-UK/UK-Sector-File/main/.github/CHANGELOG.md"`
}

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
	_, err := flags.Parse(&opts)
	if err != nil {
		if strings.Contains(err.Error(), "help") {
			os.Exit(0)
		}
		log.Panicln("Error while parsing args")
	}
	var filebytes []byte
	if opts.Url != "" {
		fmt.Println("Using online file - be aware this may take some time (10s+)")
		fmt.Println("Reading from: " + opts.Url)
		filebytes = GetWebChangelog(opts.Url)
	} else {
		fmt.Println("Using local file")
		fmt.Println("Reading from: " + opts.InputFile)
		filebytes, err = os.ReadFile(opts.InputFile)
		if err != nil {
			log.Panicln("Unable to read from " + opts.InputFile)
		}
	}
	changelog := Changelog{}
	changelog.Changes = GetChanges(filebytes)
	changelog.AIRACList, changelog.Other = changelog.ChangesSorter()
	changelog.AIRACMap, changelog.AIRACs = changelog.AIRACMapGen()
	changelog.Contributors = changelog.ContribGen()
	file := CreateFile(opts.OutputFile)
	Output(file, changelog)
	fmt.Println("Output to: " + opts.OutputFile)
	timeElapsed := time.Since(start)
	fmt.Println("Time taken:", timeElapsed)
}

func GetWebChangelog(urls string) []byte {
	resp, err := http.Get(urls)
	if err != nil {
		log.Panicln("Could not retrieve changelog")
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("Error reading response body")
	}
	return response
}

func GetChanges(filebytes []byte) []string {
	changes := []string{}
	split := strings.Split(string(filebytes), "\n") // split the document into newlines
	airacnumbers := 0
	numberReComp := regexp.MustCompile(numberRe)
	for _, s := range split {
		if airacnumbers >= 2 {
			break // if we go through multiple AIRACs then stop
		}
		if strings.Contains(s, "#") {
			airacnumbers = airacnumbers + 1
		}
		b := numberReComp.MatchString(s) // only match changelog entries
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

func OutputAIRAC(f io.Writer, c Changelog) {
	f.Write([]byte("--- AIRACs: ---" + "\n"))
	for _, key := range c.AIRACs {
		value := c.AIRACMap[fmt.Sprint(key)]
		f.Write([]byte(fmt.Sprint(key) + ":\n"))
		for _, y := range value {
			f.Write([]byte(y + "\n"))
		}
	}
}

func OutputOther(f io.Writer, c Changelog) {
	f.Write([]byte("--- Other: ---" + "\n"))
	for _, value := range c.Other {
		f.Write([]byte(value + "\n"))
	}
}

func OutputContribs(f io.Writer, c Changelog) {
	f.Write([]byte("--- Contributors: ---" + "\n"))
	for _, value := range c.Contributors {
		f.Write([]byte(value + "\n"))
	}
}

func Output(f io.Writer, c Changelog) {
	OutputAIRAC(f, c)
	f.Write([]byte("\n"))
	OutputOther(f, c)
	f.Write([]byte("\n"))
	OutputContribs(f, c)

}
