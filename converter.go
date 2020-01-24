package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Source struct {
	Id string `json:"id"`
    Title string `json:"title"`
    URL string `json:"url"`
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

func main() {

	flag.Parse()

	catalogFile, err := os.Open(flag.Arg(0))
	if err != nil {
	    fmt.Println(err)
	}
	defer catalogFile.Close()
	
	byteValue, _ := ioutil.ReadAll(catalogFile)
	var catalog Catalog
	json.Unmarshal(byteValue, &catalog)

	factions := make([]StratagemFaction, 0)
	
	for i:=0; i<len(catalog.Factions); i++ {
		
		factionFile, err := os.Open(catalog.Factions[i]+".json")
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
		fmt.Println(faction.Code + " containing " + strconv.Itoa(len(faction.Stratagems)) + " tactics.")
		factions = append(factions, faction)
	}

	catalog.Package.Factions = factions
	
	strataGemJson, _ := json.MarshalIndent(catalog.Package, "", "\t")
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
