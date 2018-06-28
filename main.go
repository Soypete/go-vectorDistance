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

func getWeightedAverage(info []weighting) float64 {
	var sumtop float64
	var sumbottom float64

	for _, sw := range info {
		sumtop += (sw.score * float64(sw.weight))
		sumbottom += float64(sw.weight)
	}

	return sumtop / sumbottom
}

// GetDandbMatchScore returns a score representing the likelyhood that a business already exists.
func GetDandbMatchScore(sName, rName, sCity, rCity, sState, rState, sZip, rZip, sPhone, rPhone, sBin, rBin, sDuns, rDuns, ad1, ad2 string) float64 {
	var weights []weighting
	address1 := ad1
	address2 := ad2

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
	address1 = strings.ToLower(address1)
	address2 = strings.ToLower(address2)

	address1Slice := strings.SplitAfter(address1, "")
	for i := range address1Slice {
		if wordReplacements[address1Slice[i]] != "" {
			address1Slice[i] = fmt.Sprintf("%s ", wordReplacements[address1Slice[i]])
		}
	}

	address2Slice := strings.SplitAfter(address2, "")
	for i := range address2Slice {
		if wordReplacements[address2Slice[i]] != "" {
			address2Slice[i] = fmt.Sprintf("%s ", wordReplacements[address2Slice[i]])
		}
	}

	address1 = strings.Join(address1Slice, "")
	address2 = strings.Join(address2Slice, "")

	if ad1 != "" {
		jaro := smetrics.Jaro(address1, address2)
		jaroWinkler := smetrics.JaroWinkler(address1, address2, 0.7, 4)
		fmt.Printf("JaroW: %v \n", jaroWinkler)
		wagnerFisher := smetrics.WagnerFischer(address1, address2, 1, 2, 4)

		wagnerFisher = (100 - wagnerFisher) / 100
		score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
		weights = append(weights, weighting{score: score, weight: 1})
	}
	if sName != "" {
		jaro := smetrics.Jaro(sName, rName)
		jaroWinkler := smetrics.JaroWinkler(sName, rName, 0.7, 4)
		fmt.Printf("JaroW: %v\n", jaroWinkler)
		wagnerFisher := smetrics.WagnerFischer(sName, rName, 1, 2, 4)

		wagnerFisher = (100 - wagnerFisher) / 100
		score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
		weights = append(weights, weighting{score: score, weight: 1})
	}

	if sCity != "" {
		jaro := smetrics.Jaro(sCity, rCity)
		jaroWinkler := smetrics.JaroWinkler(sCity, rCity, 0.7, 4)
		fmt.Printf("JaroW: %v \n", jaroWinkler)
		wagnerFisher := smetrics.WagnerFischer(sCity, rCity, 1, 2, 4)

		wagnerFisher = (100 - wagnerFisher) / 100
		score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
		weights = append(weights, weighting{score: score, weight: 1})
	}

	if sState != "" {
		jaro := smetrics.Jaro(sState, rState)
		jaroWinkler := smetrics.JaroWinkler(sState, rState, 0.7, 4)
		fmt.Printf("JaroW: %v\n", jaroWinkler)
		wagnerFisher := smetrics.WagnerFischer(sState, rState, 1, 2, 4)

		wagnerFisher = (100 - wagnerFisher) / 100
		score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
		weights = append(weights, weighting{score: score, weight: 1})
	}

	if sZip != "" {
		if strings.Contains(rZip, sZip) {
			score := float64(1)
			weights = append(weights, weighting{score: score, weight: 2})
		} else {
			jaro := smetrics.Jaro(sZip, rZip)
			jaroWinkler := smetrics.JaroWinkler(sZip, rZip, 0.7, 4)
			fmt.Printf("JaroW: %v\n", jaroWinkler)
			wagnerFisher := smetrics.WagnerFischer(sZip, rZip, 1, 2, 4)

			wagnerFisher = (100 - wagnerFisher) / 100
			score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
			weights = append(weights, weighting{score: score, weight: 2})
		}
	}

	if sPhone != "" {
		jaro := smetrics.Jaro(sPhone, rPhone)
		jaroWinkler := smetrics.JaroWinkler(sPhone, rPhone, 0.7, 4)
		fmt.Printf("JaroW: %v\n", jaroWinkler)
		wagnerFisher := smetrics.WagnerFischer(sPhone, rPhone, 1, 2, 4)

		wagnerFisher = (100 - wagnerFisher) / 100
		score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
		weights = append(weights, weighting{score: score, weight: 1})
	}

	if sBin != "" {
		jaro := smetrics.Jaro(sBin, rBin)
		jaroWinkler := smetrics.JaroWinkler(sBin, rBin, 0.7, 4)
		fmt.Printf("JaroW: %v\n", jaroWinkler)
		wagnerFisher := smetrics.WagnerFischer(sBin, rBin, 1, 2, 4)

		wagnerFisher = (100 - wagnerFisher) / 100
		score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
		weights = append(weights, weighting{score: score, weight: 10})
	}

	if sDuns != "" {
		jaro := smetrics.Jaro(sDuns, rDuns)
		jaroWinkler := smetrics.JaroWinkler(address1, address2, 0.7, 4)
		fmt.Printf("JaroW: %v\n", jaroWinkler)
		wagnerFisher := smetrics.WagnerFischer(address1, address2, 1, 2, 4)

		wagnerFisher = (100 - wagnerFisher) / 100
		score := (jaro + jaroWinkler + float64(wagnerFisher)) / 3
		weights = append(weights, weighting{score: score, weight: 10})
	}

	return getWeightedAverage(weights)
}

func main() {
	in, err := os.Open("/Users/mpeterson/code/pair_programming/distance/20180622/20180622_results_searches.csv")

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
		// fmt.Println(data.SearchStreet)
		// fmt.Println(data.ResultStreet)
		time.Sleep(500 * time.Millisecond)
		score := GetDandbMatchScore(data.SearchName, data.ResultName, data.SearchCity, data.ResultCity, data.SearchState, data.ResultState, data.SearchZip, data.ResultZip, data.SearchPhone, data.ResultPhone, data.SearchBin, data.ResultBin, data.SearchDuns, data.ResultDuns, data.SearchStreet, data.ResultStreet)
		fmt.Println(score)
	}
}
