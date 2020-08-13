package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"sync"
	"time"
	"bufio"
	"strings"
	"net/http"
	"io/ioutil"
)

var target string
var pathList string
var verbose bool
var wait int
var NoContext bool
var wg sync.WaitGroup
var count int

func init() {
	flag.StringVar(&target, "target", "", "Target Url")
	flag.StringVar(&pathList, "l", "", "Admin List")
	flag.BoolVar(&verbose, "v", false, "Verbose Output")
	flag.IntVar(&wait, "t", 0, "Time to wait between requests in ms")
	flag.BoolVar(&NoContext, "no-context", false, "Disable Context Display")
	flag.Parse()

	if target == "" || (!(strings.HasPrefix(target, "http://")) && !(strings.HasPrefix(target, "https://"))) {
		log.Fatal("Incorrect or Missing Target Url")
	}

	if pathList == "" {
		log.Fatal("Missing Admin List")
	}

	if !NoContext {
		fmt.Println(`
           __            __               __                    __         
 .---.-.--|  |.--------.|__|.-----.______|  |.-----.----.---.-.|  |_.-----.
 |  _  |  _  ||        ||  ||     |______|  ||  _  |  __|  _  ||   _|  -__|
 |___._|_____||__|__|__||__||__|__|      |__||_____|____|___._||____|_____|   

                      -- Made with <3 by luxunator --                             
			`)
	}

}

func request(target_url string, path string) error {
	wg.Add(1)

	response, err := http.Get(target_url)
	if err != nil {
		fmt.Println("Request Error: ", target_url, err)
		return err
	}
	defer response.Body.Close()
	
	if response.StatusCode != 404 {

		if verbose {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error Reading Body Length")
				return err
			}
			fmt.Println(target_url, count+1, response.StatusCode, len(body))
		} else {
			fmt.Println(target_url, count+1)
		}		
	}
	count+=1
	wg.Done()
	return nil
}

func main() {

	if !NoContext {
		fmt.Println("Scanning - ", target, "\n")
	}
	
	paths, err := os.Open(pathList)
    if err != nil {
        log.Fatal("Error Reading Input File")
    }
    defer paths.Close()

	scanner := bufio.NewScanner(paths)

	start := time.Now()

	for scanner.Scan() {
		go request(target + scanner.Text(), scanner.Text())
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}

	wg.Wait()

	elapsed := time.Since(start)

	if !NoContext {
		fmt.Println("\nScan Time:", elapsed)
	}
}