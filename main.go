package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xrash/smetrics"
)

type CsvLine struct {
	DataSource       string
	CreatedAt        string
	IsSelectedMatch  string
	ResultConfidence string
	ResultName       string
	SearchName       string
	ResultCity       string
	SearchCity       string
	ResultState      string
	SearchState      string
	ResultZip        string
	SearchZip        string
	ResultPhone      string
	SearchPhone      string
	ResultBin        string
	SearchBin        string
	colbool1         bool
	colbool2         bool
	colbool3         bool
	colbool4         bool
	ResultDuns       string
	SearchDuns       string
	ResultStreet     string
	SearchStreet     string
}

type weighting struct {
	score  float64
	weight int
}

//getVectorDistance uses Jaro distance Jaro-Winkler distance and Levenstein distance calculated
// by the WagnerFisher method. They are averaged to give the most accurate vector distance
func getVectorDistance(searchVal, returnVal string) float64 {
	jaro := smetrics.Jaro(searchVal, returnVal)
	jaroWinkler := smetrics.JaroWinkler(searchVal, returnVal, 0.7, 4)
	wagnerFisher := smetrics.WagnerFischer(searchVal, returnVal, 1, 2, 4)

	wagnerFisher = (100 - wagnerFisher) / 100
	score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
	return score
}

func getWeightedAverage(info []weighting) float64 {
	var sumtop float64
	var sumbottom float64

	for _, sw := range info {
		sumtop += (sw.score * float64(sw.weight))
		sumbottom += float64(sw.weight)
	}
	// formula is ∑w_n*x_n/∑w_n where w is the weight and x is the score
	return sumtop / sumbottom
}

// GetDandbMatchScore returns a score representing the likelyhood that a business already exists.
func GetDandbMatchScore(sName, rName, sCity, rCity, sState, rState, sZip, rZip, sPhone, rPhone, sBin, rBin, sDuns, rDuns, ad1, ad2 string) float64 {
	var weights []weighting
	sAddress := ad1
	rAddress := ad2

	fmt.Printf(" search name: %s found name: %s sZip: %s rZIp: %s sBin: %s rBin: %s \n", sName, rName, sZip, rZip, sBin, rBin)

	//TODO: convert provided string into most general format
	wordReplacements := map[string]string{
		"apartment": "apt",
		"avenue":    "ave",
		"boulevard": "blvd",
		"court":     "ct",
		"drive":     "dr",
		"east":      "e",
		"highway":   "hwy",
		"lane":      "ln",
		"north":     "n",
		"road":      "rd",
		"south":     "s",
		"street":    "st",
		"suite":     "ste",
		"west":      "w",
	}

	// simple step to eliminate typing error
	sAddress = strings.ToLower(sAddress)
	rAddress = strings.ToLower(rAddress)

	sAddressSlice := strings.SplitAfter(sAddress, "")
	for i := range sAddressSlice {
		if wordReplacements[sAddressSlice[i]] != "" {
			sAddressSlice[i] = fmt.Sprintf("%s ", wordReplacements[sAddressSlice[i]])
		}
	}

	rAddressSlice := strings.SplitAfter(rAddress, "")
	for i := range rAddressSlice {
		if wordReplacements[rAddressSlice[i]] != "" {
			rAddressSlice[i] = fmt.Sprintf("%s ", wordReplacements[rAddressSlice[i]])
		}
	}

	sAddress = strings.Join(sAddressSlice, "")
	rAddress = strings.Join(rAddressSlice, "")

	if sAddress != "" {
		score := getVectorDistance(sAddress, rAddress)
		weights = append(weights, weighting{score: score, weight: 1})
	}

	if sName != "" {
		score := getVectorDistance(sName, rName)
		weights = append(weights, weighting{score: score, weight: 1})
	}

	if sCity != "" {
		score := getVectorDistance(sCity, rCity)
		weights = append(weights, weighting{score: score, weight: 1})
	}

	if sState != "" {
		score := getVectorDistance(sState, rState)
		weights = append(weights, weighting{score: score, weight: 1})
	}

	if sZip != "" {
		if strings.Contains(rZip, sZip) {
			score := float64(1)
			weights = append(weights, weighting{score: score, weight: 2})
		} else {
			score := getVectorDistance(sZip, rZip)
			weights = append(weights, weighting{score: score, weight: 2})
		}
	}

	return getWeightedAverage(weights)
}

func main() {
	in, err := os.Open("20180622/20180622_results_searches.csv")

	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(in).ReadAll()
	if err != nil {
		panic(err)
	}
	count := 1
	// Loop through lines & turn into object
	for _, line := range lines {
		if count == 1 {
			count++
			continue
		}
		data := CsvLine{
			DataSource:   line[0],
			ResultName:   line[4],
			SearchName:   line[5],
			ResultCity:   line[6],
			SearchCity:   line[7],
			ResultState:  line[8],
			SearchState:  line[9],
			ResultZip:    line[10],
			SearchZip:    line[11],
			ResultPhone:  line[12],
			SearchPhone:  line[13],
			ResultBin:    line[14],
			SearchBin:    line[15],
			ResultDuns:   line[20],
			SearchDuns:   line[21],
			ResultStreet: line[22],
			SearchStreet: line[23],
		}
		fmt.Println(data.SearchStreet)
		fmt.Println(data.ResultStreet)
		time.Sleep(500 * time.Millisecond)
		score := GetDandbMatchScore(data.SearchName, data.ResultName, data.SearchCity, data.ResultCity, data.SearchState, data.ResultState, data.SearchZip, data.ResultZip, data.SearchPhone, data.ResultPhone, data.SearchBin, data.ResultBin, data.SearchDuns, data.ResultDuns, data.SearchStreet, data.ResultStreet)
		fmt.Println(score)
	}
}
