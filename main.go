package main

import (
	"flag"
	"fmt"
	"github.com/grokify/html-strip-tags-go"
	"io/ioutil"
	//"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Result struct, populated with the number of results
type Result struct {
	NumResults int
	Results    []string
	Error      string
}

func main() {
	flag.Parse()
	args := flag.Args()

	query := CreateQuery(args)

	results := FetchResults(query)

	PrettyPrint(results)
}

// PrettyPrint prints Result in a user readable way
func PrettyPrint(results Result) {
	if results.Error != "" {
		fmt.Fprintf(os.Stderr, "Error: %s\n", results.Error)
		return
	}
	fmt.Printf("Found %d results.", results.NumResults)
	if results.NumResults > 5 {
		fmt.Printf(" Showing first five:")
	}
	fmt.Println("")
	for i, elem := range results.Results {
		if i > 4 {
			return
		}
		fmt.Printf("%d. %s\n", i+1, elem)
	}
}

// FetchResults Queries the OEIS and returns it's answer
func FetchResults(query string) (result Result) {
	resp, err := http.Get(query)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Network error\n")
		//log.Panic("Error fetching oeis data",err)
		os.Exit(-1)
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Data returned from oeis.org not readable\n")
		//log.Panic("Error reading oeis data: ", err)
		os.Exit(-2)
	}

	result = HTMLToResult(string(html))

	return
}

// HTMLToResult turns HTML into a Result struct
func HTMLToResult(data string) (result Result) {
	result = Result{
		NumResults: 0,
		Results:    []string{},
	}

	badqueryresult := "Sorry, the page you requested was not found."
	if strings.Contains(data, badqueryresult) {
		result.Error = badqueryresult
		return
	}

	noresregexp := "Sorry, but the terms do not match anything in the table."
	// No Results
	if strings.Contains(data, noresregexp) {
		result.Error = noresregexp
		return
	}
	regex := regexp.MustCompile("Displaying .* of (.*) results found.")

	m := regex.FindAllString(data, 1)
	for _, elem := range m {
		result.NumResults, _ = strconv.Atoi(strings.Split(elem, " ")[3])
	}

	result.Results = GetTopFiveResults(data)

	return
}

// GetTopFiveResults gets 5 or less titles from the OEIS data
func GetTopFiveResults(data string) (results []string) {
	//	regex := regexp.MustCompile(`(?s)<td valign=top align=left>(.*?)<td`)
	regex := regexp.MustCompile(`(?m)<td valign=top align=left>\n(.*)\n`)
	regexres := regex.FindAllIndex([]byte(data), -1)
	for i, elems := range regexres {
		if i > 4 {
			return
		}
		mod := strings.TrimSpace(data[elems[0]+27 : elems[1]])
		mod = strip.StripTags(mod)
		mod = fmt.Sprintf("%s", strings.TrimSpace(mod))
		results = append(results, mod)
	}

	return
}

// CreateQuery turns the slice of arguments into a query, it ignores all not int64 arguments
func CreateQuery(args []string) string {
	query := "https://oeis.org/search?q="
	for _, elem := range args {
		if _, err := strconv.ParseInt(elem, 10, 64); err == nil {
			query += elem
			query += "%20"
		}
	}
	query += "&sort=&language=&go=Search"
	return query
}
