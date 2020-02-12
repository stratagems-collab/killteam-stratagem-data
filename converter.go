package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Source struct {
	Id string `json:"id"`
    Title string `json:"title"`
    URL string `json:"url"`
}

type SourceReference struct {
	Id string `json:"id"`
    Page string `json:"page"`
    Data string `json:"data,omitempty"`
}

type Catalog struct {
	Factions []string `json:"factions"`
	Package StratagemPackage `json:"package"`
}

type StratagemPackage struct {
	PackageName string             `json:"packageName"`
	VersionName string             `json:"versionName"`
	VersionCode int                `json:"versionCode"`
	Factions    []StratagemFaction `json:"factions"`
	Sources []Source `json:"sources"`
}

type StratagemFaction struct {
	Code    string 	   `json:"code"`
	FactionKeyword string      `json:"faction"`
	Stratagems     []Stratagem `json:"tactics"`
}

type Stratagem struct {
	Title             string         `json:"title"`
	Subtitle          string         `json:"sub"`
	Description       string         `json:"desc"`
	Cost              int            `json:"cp"`
	Phases            StratagemPhase `json:"phases"`
	Keywords          []string       `json:"keywords,omitempty"`
	Equipment         []string       `json:"equipment,omitempty"`
	SpecialistKeyword string         `json:"specialist,omitempty"`
	SpecialistLevel   int            `json:"level,omitempty"`
	SourceReference	  *SourceReference `json:"source,omitempty"`
}

type StratagemPhase struct {
	Move     bool `json:"move,omitempty"`
	Psychic  bool `json:"psychic,omitempty"`
	Shoot    bool `json:"shoot,omitempty"`
	Fight    bool `json:"fight,omitempty"`
	Round    bool `json:"round,omitempty"`
	RoundOne bool `json:"round_one,omitempty"`
	Morale   bool `json:"morale,omitempty"`
	Event    bool `json:"event,omitempty"`
}

type Options struct {
	Filename     string
	UseGitRev	 bool
}

func main() {

	options := Options{}

	flag.BoolVar(&options.UseGitRev, "gitversion", false, "use the version number from the latest git-tag and commit count")
	flag.StringVar(&options.Filename, "catalog", "catalog.json", "the filename of the catalog to process")
	flag.Parse()

	catalogFile, err := os.Open(options.Filename)
	if err != nil {
	    fmt.Println(err)
	}
	defer catalogFile.Close()
	
	byteValue, _ := ioutil.ReadAll(catalogFile)
	var catalog Catalog
	json.Unmarshal(byteValue, &catalog)

	factions := make([]StratagemFaction, 0)
	
	for i:=0; i<len(catalog.Factions); i++ {
		filename := catalog.Factions[i]+".json"
		factionFile, err := os.Open(filename)
		if err != nil {
		    fmt.Println(err)
		}
		defer factionFile.Close()
		byteValue, _ := ioutil.ReadAll(factionFile)
		var faction StratagemFaction
		err = json.Unmarshal(byteValue, &faction)
		if jsonError, ok := err.(*json.SyntaxError); ok {
			line,char,_ := lineAndCharacter(string(byteValue), int(jsonError.Offset))
			fmt.Println("Error while parsing " + catalog.Factions[i]+".json" + " line: " + strconv.Itoa(line) + ", char: " + strconv.Itoa(char))
			fmt.Println(err)
		}
		for j, t := range faction.Stratagems {
			for k, ta := range faction.Stratagems {
				if j != k {
					if ta.Title == t.Title {
						panic("error duplicate tactic in " + filename + ":" + ta.Title)
					}
				}
			}
		}

		fmt.Println(faction.Code + " containing " + strconv.Itoa(len(faction.Stratagems)) + " tactics.")
		factions = append(factions, faction)
	}

	//Replace version info when applicable
	if options.UseGitRev {
		r, err := git.PlainOpen("./")
		if err != nil {
			panic(err)
		}
		tag, err := GetLatestTagFromRepository(r)
		if err != nil {
			panic(err)
		}
		p := strings.Split(tag, "/")
		catalog.Package.VersionName = p[len(p)-1]
		commits, err := GetNumberOfCommits(r)
		if err != nil {
			panic(err)
		}
		catalog.Package.VersionCode = commits
	}

	//serialize minimal header
	strataGemJson, _ := json.MarshalIndent(catalog.Package, "", "\t")
	err = ioutil.WriteFile("header.json", strataGemJson, 0644)
	if err != nil {
	    fmt.Println(err)
	}

	catalog.Package.Factions = factions
	
	//serialize full catalog
	strataGemJson, _ = json.MarshalIndent(catalog.Package, "", "\t")
	err = ioutil.WriteFile("data.json", strataGemJson, 0644)
	if err != nil {
	    fmt.Println(err)
	}
}

func lineAndCharacter(input string, offset int) (line int, character int, err error) {
	lf := rune(0x0A)
	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("Couldn't find offset %d within the input.", offset)
	}
	line = 1
	for i, b := range input {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}
	return line, character, nil
}

func GetNumberOfCommits(repository *git.Repository) (int, error) {
	ref, err := repository.Head()
	if err != nil {
		return -1, err
	}

	cIter, err := repository.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return -1, err
	}

	var cCount int
	err = cIter.ForEach(func(c *object.Commit) error {
		cCount++

		return nil
	})
	
	if err != nil {
		return -1, err
	}

	return cCount, nil
}

func GetLatestTagFromRepository(repository *git.Repository) (string, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return "", err
	}

	var latestTagCommit *object.Commit
	var latestTagName string
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := repository.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := repository.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return latestTagName, nil
}
