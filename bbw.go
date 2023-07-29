package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const url = "https://pentester.land/writeups.json"

type WriteupList struct {
	Data []Writeup `json:"data"`
}

type Writeup struct {
	Links           []WriteupTitle `json:"Links"`
	Authors         []string       `json:"Authors"`
	Programs        []string       `json:"Programs"`
	Bugs            []string       `json:"Bugs"`
	Bounty          string         `json:"Bounty"`
	PublicationDate string         `json:"PublicationDate"`
	AddedDate       string         `json:"AddedDate"`
}

type WriteupTitle struct {
	Title string `json:"Title"`
	Link  string `json:"Link"`
}

func main() {
	printBanner()
	bbw_config, err := readConfigFile()
	writeups := GetWriteUps()

	if len(os.Args) > 1 {
		numOfWriteups, err := strconv.Atoi(os.Args[1])
		checkError(err)
		if len(writeups.Data) < numOfWriteups {
			fmt.Println("That's a huge number of writeups!")
		} else {
			for i := 0; i < numOfWriteups; i++ {
				printWriteupDetails(writeups.Data[i])
			}
		}
		os.Exit(0)
	}

	if err != nil && os.IsNotExist(err) {
		printWriteupDetails(writeups.Data[0])
		saveLatestWriteup(writeups.Data[0])
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}

	var lastReadWriteup Writeup
	json.Unmarshal([]byte(bbw_config), &lastReadWriteup)

	if isSameWriteup(writeups.Data[0], lastReadWriteup) {
		fmt.Println("No new writeups available.")
		os.Exit(0)
	}

	for _, writeup := range writeups.Data {
		if isSameWriteup(writeup, lastReadWriteup) {
			saveLatestWriteup(writeups.Data[0])
			os.Exit(0)
		} else {
			printWriteupDetails(writeup)
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func isSameWriteup(writeup1, writeup2 Writeup) bool {
	writeup1_signature := getWriteupSignature(writeup1)
	writeup2_signature := getWriteupSignature(writeup2)

	for i := 0; i < 4; i++ {
		if writeup1_signature[i] != writeup2_signature[i] {
			return false
		}
	}

	return true
}

func getWriteupSignature(writeup Writeup) []string {
	var titles []string
	var links []string

	for _, metadata := range writeup.Links {
		titles = append(titles, metadata.Title)
		links = append(links, metadata.Link)
	}

	return []string{strings.Join(titles, " | "), strings.Join(links, " | "), strings.Join(writeup.Authors, " | "), writeup.PublicationDate}
}

func saveLatestWriteup(writeup Writeup) {
	bbwgo_config := getConfigFilename()
	raw_bytes, err := json.Marshal(writeup)
	checkError(err)
	err = os.WriteFile(bbwgo_config, raw_bytes, 0666)
	checkError(err)
}

func printWriteupDetails(writeup Writeup) {
	var titles []string
	var links []string

	for _, metadata := range writeup.Links {
		titles = append(titles, metadata.Title)
		links = append(links, metadata.Link)
	}

	fmt.Println("-------------------------------------------------------")
	fmt.Println("Title :", strings.Join(titles, " | "))
	fmt.Println("Link(s) :", strings.Join(links, " | "))
	fmt.Println("Author(s) :", strings.Join(writeup.Authors, " | "))
	fmt.Println("Program(s) :", strings.Join(writeup.Programs, " | "))
	fmt.Println("Bug(s) :", strings.Join(writeup.Bugs, " "))
	fmt.Println("Bounty :", writeup.Bounty)
	fmt.Println("Publication date :", writeup.PublicationDate)
	fmt.Println("Added date :", writeup.AddedDate)
}

func GetWriteUps() WriteupList {
	response, err := http.Get(url)
	checkError(err)
	defer response.Body.Close()

	byteBody, err := io.ReadAll(response.Body)
	checkError(err)

	var writeups WriteupList
	json.Unmarshal(byteBody, &writeups)

	return writeups
}

func readConfigFile() (string, error) {
	bbwgo_config := getConfigFilename()

	bbwgo_config_data, err := os.ReadFile(bbwgo_config)
	if os.IsNotExist(err) {
		return "", err
	}

	return string(bbwgo_config_data), nil
}

func getConfigFilename() string {
	home_dir, err := os.UserHomeDir()
	checkError(err)
	bbwgo_config := home_dir + "/.bbwgo.json"

	return bbwgo_config
}

func printBanner() {
	fmt.Println(`
╔╗ ╔╗ ╦ ╦
╠╩╗╠╩╗║║║
╚═╝╚═╝╚╩╝	

- Get list of latest bug bounty writeups
- Author: n3hal_ (github.com/Nehal-Zaman)
	`)
}
