package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	Input: cat a list of URLs/IPs via stdin

	Aim: Perform a GET request to the current target and report it's status code
		- Allow the user to specify the output to only show the requested status code results (i.e -s 404)

	Reason: Not really sure how this differs to just using gobuster or ffuf in this intended way, but here we are. Smaller tool usage I guess?
*/

var out io.Writer = os.Stdout
var quietMode bool

var targetStatusCodes = []int{}

func main() {
	var targetStatusCodesFlag string
	flag.StringVar(&targetStatusCodesFlag, "s", "", "Only look for these Status Codes (i.e -s 200 or -s 200,404")

	quietModeFlag := flag.Bool("q", false, "Only output the data we care about")
	flag.Parse()

	quietMode = *quietModeFlag

	// check if we have any target status codes
	if targetStatusCodesFlag != "" {
		if !quietMode {
			fmt.Println("Status code flag supplied. Parsing into array")
			fmt.Println("Input:", targetStatusCodesFlag)
		}

		splitArg := strings.Split(targetStatusCodesFlag, ",")

		for _, sa := range splitArg {
			saInt, err := strconv.Atoi(sa)
			if err != nil {
				fmt.Println("Couldn't convert to int, skipping: ", sa)
				fmt.Println("err:", err)
				continue
			}
			targetStatusCodes = append(targetStatusCodes, saInt)
		}

		if !quietMode {
			fmt.Println("Only looking for status codes:")
			for _, sc := range targetStatusCodes {
				fmt.Println(sc)
			}
		}
	} else {
		if !quietMode {
			fmt.Println("No status codes supplied.")
		}
		// this will output everything regardless of status code
	}

	if !quietMode {
		banner()
		fmt.Println("")
	}

	writer := bufio.NewWriter(out)
	targetDomains := make(chan string, 1)
	var wg sync.WaitGroup

	ch := readStdin()
	go func() {
		//translate stdin channel to domains channel
		for u := range ch {
			targetDomains <- u
		}
		close(targetDomains)
	}()

	// flush to writer periodically
	t := time.NewTicker(time.Millisecond * 500)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				writer.Flush()
			}
		}
	}()

	for u := range targetDomains {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if !quietMode {
				fmt.Println("Checking:", url)
			}

			// we assume the URL passed in already has it's http/https protocol set up ready.
			statusCode := makeRequest(url, quietMode)
			if statusCode == -1 {
				// we didn't get a status code (i.e was an error), report it anyway for our users if they want all
				if !quietMode {
					fmt.Println("No status code returned for:", url)
				}
			} else {
				if len(targetStatusCodes) > 0 {
					if contains(targetStatusCodes, statusCode) {
						fmt.Printf("[%d] %s\n", statusCode, url)
					}
				} else {
					fmt.Printf("[%d] %s\n", statusCode, url)
				}
			}
		}(u)
	}

	wg.Wait()

	// just in case anything is still in buffer
	writer.Flush()
}

func banner() {
	fmt.Println("---------------------------------------------------")
	fmt.Println("StatusQuode -> Crawl3r")
	fmt.Println("Makes a request to the target URL and spits out the status codes")
	fmt.Println("")
	fmt.Println("Run again with -q for cleaner output")
	fmt.Println("---------------------------------------------------")
}

func readStdin() <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			url := strings.ToLower(sc.Text())
			if url != "" {
				lines <- url
			}
		}
	}()
	return lines
}

func makeRequest(url string, quietMode bool) int {
	resp, err := http.Get(url) // TODO: shall we HEAD this instead?! Lesser load than GET? We don't care for source just the code
	if err != nil {
		if !quietMode {
			fmt.Println("[error] performing the request to:", url)
		}
		return -1
	}
	defer resp.Body.Close() // TODO: do we need this?

	return resp.StatusCode
}

func contains(arr []int, a int) bool {
	for _, i := range arr {
		if i == a {
			return true
		}
	}
	return false
}
