package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestHTMLToResult(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		result   Result
	}{
		//	{name: "Bug #1", filename: "test_para2.html", result: Result{0, []string{}, "Sorry, but the terms do not match anything in the table."}},
		{name: "No results", filename: "test_noresult.html", result: Result{0, []string{}, "Sorry, but the terms do not match anything in the table."}},
		{name: "Bad Query", filename: "test_badquery.html", result: Result{0, []string{}, "Sorry, the page you requested was not found."}},
		{name: "Truncated file", filename: "test_half.html", result: Result{11196, []string{
			"The positive integers. Also called the natural numbers, the whole numbers or the counting numbers, but these terms are ambiguous.",
			"Digital sum (i.e., sum of digits) of n; also called digsum(n).",
			"Powers of primes. Alternatively, 1 and the prime powers (p^k, p prime, k &gt;= 1).",
			"The nonnegative integers.",
		}, ""}},
		{name: "Lots of results", filename: "test_1_2_3_4.html", result: Result{11196, []string{
			"The positive integers. Also called the natural numbers, the whole numbers or the counting numbers, but these terms are ambiguous.",
			"Digital sum (i.e., sum of digits) of n; also called digsum(n).",
			"Powers of primes. Alternatively, 1 and the prime powers (p^k, p prime, k &gt;= 1).",
			"The nonnegative integers.",
			"Palindromes in base 10.",
			//"Numbers k such that (266*10^k + 1)/3 is prime.",
		}, ""}},
		{name: "Four (4) results", filename: "test_parasitic.html", result: Result{4, []string{
			"Numbers k such that k and 4*k are anagrams.",
			"Numbers m with the property that shifting the rightmost digit of m to the left end multiplies the number by 4.",
			"Numbers that are proper divisors of the number you get by rotating digits right once.",
			"Non-repdigit numbers n such that n divides A045876(n).",
		}, ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ioutil.ReadFile("testdata/" + tt.filename)
			if err != nil {
				t.Errorf("Missing test file %v", err)
				return
			}
			if gotResult := HTMLToResult(string(data)); !reflect.DeepEqual(gotResult, tt.result) {
				t.Errorf("HTMLToResult(%v) = %v, want %v", tt.filename, gotResult, tt.result)
			}
		})
	}
}

func TestCreateQuery(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "No argument", args: args{args: []string{}}, want: "https://oeis.org/search?q=&sort=&language=&go=Search"},
		{name: "One argument", args: args{args: []string{"1"}}, want: "https://oeis.org/search?q=1%20&sort=&language=&go=Search"},
		{name: "Two argument", args: args{args: []string{"1", "2"}}, want: "https://oeis.org/search?q=1%202%20&sort=&language=&go=Search"},
		{name: "Mixed Type", args: args{args: []string{"1", "A", "2"}}, want: "https://oeis.org/search?q=1%202%20&sort=&language=&go=Search"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateQuery(tt.args.args); got != tt.want {
				t.Errorf("CreateQuery(%v) = %v, want %v", tt.args.args, got, tt.want)
			}
		})
	}
}

func TestFetchResults(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name       string
		args       args
		wantResult Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := FetchResults(tt.args.query); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FetchResults(%v) = %v, want %v", tt.args.query, gotResult, tt.wantResult)
			}
		})
	}
}

func TestGetTopFiveResults(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		results  []string
	}{
		{name: "", filename: "test_1_2_3_4.html", results: []string{
			"The positive integers. Also called the natural numbers, the whole numbers or the counting numbers, but these terms are ambiguous.",
			"Digital sum (i.e., sum of digits) of n; also called digsum(n).",
			"Powers of primes. Alternatively, 1 and the prime powers (p^k, p prime, k &gt;= 1).",
			"The nonnegative integers.",
			"Palindromes in base 10.",
			//"Numbers k such that (266*10^k + 1)/3 is prime.",
		}},
		{name: "", filename: "test_half.html", results: []string{
			"The positive integers. Also called the natural numbers, the whole numbers or the counting numbers, but these terms are ambiguous.",
			"Digital sum (i.e., sum of digits) of n; also called digsum(n).",
			"Powers of primes. Alternatively, 1 and the prime powers (p^k, p prime, k &gt;= 1).",
			"The nonnegative integers.",
		}},
		{name: "", filename: "test_parasitic.html", results: []string{
			"Numbers k such that k and 4*k are anagrams.",
			"Numbers m with the property that shifting the rightmost digit of m to the left end multiplies the number by 4.",
			"Numbers that are proper divisors of the number you get by rotating digits right once.",
			"Non-repdigit numbers n such that n divides A045876(n).",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ioutil.ReadFile("testdata/" + tt.filename)
			if err != nil {
				t.Errorf("Missing test file %v", err)
				return
			}
			if gotResults := GetTopFiveResults(string(data)); !reflect.DeepEqual(gotResults, tt.results) {
				t.Errorf("GetTopFiveResults(%v) size %d = %v, want %v", tt.filename, len(gotResults), gotResults, tt.results)
			}
		})
	}
}

func TestPrettyPrint(t *testing.T) {
	tests := []struct {
		name    string
		results Result
	}{
		{name: "No results", results: Result{0, []string{}, ""}},
		{name: "Four (4) Results", results: Result{4, []string{"", "", "", ""}, ""}},
		{name: "Six (6) Results", results: Result{6, []string{"", "", "", "", "", ""}, ""}},
		{name: "Error", results: Result{6, []string{"", "", "", "", "", ""}, "Error message"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrettyPrint(tt.results)
		})
	}
}
