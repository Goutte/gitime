package main

import (
	"fmt"
	"github.com/tsuyoshiwada/go-gitlog"
	"log"
	"regexp"
	"strconv"
	"strings"
)

/*

Purpose
-------

Collect, addition and return all the `/spend` and `/spent` time-tracking directives in git commit messages.

This only looks at the `git log` of the currently checked out branch.


Usage
-----

	go run gitime.go


Dependencies
------------

	go get -u github.com/tsuyoshiwada/go-gitlog


*/

func main() {
	git := gitlog.New(&gitlog.Config{})

	commits, err := git.Log(nil, nil)
	if err != nil {
		log.Fatalln("Cannot read git log:", err)
	}

	ts := &TimeSpent{}
	for _, commit := range commits {
		ts.Add(CollectTimeSpent(commit.Subject))
		ts.Add(CollectTimeSpent(commit.Body))
	}

	fmt.Printf(ts.String() + "\n")
	fmt.Printf("%d minutes\n", ts.ToMinutes())
}

type TimeSpent struct {
	Months  float64
	Weeks   float64
	Days    float64
	Hours   float64
	Minutes float64
}

func (ts *TimeSpent) String() string {
	s := ""

	if ts.Months > 0.0 {
		s += fmt.Sprintf("%.1f month", ts.Months)
		if ts.Months >= 2.0 {
			s += "s"
		}
	}
	if ts.Weeks > 0.0 {
		if s != "" {
			s += " "
		}
		s += fmt.Sprintf("%.1f week", ts.Weeks)
		if ts.Weeks >= 2.0 {
			s += "s"
		}
	}
	if ts.Days > 0.0 {
		if s != "" {
			s += " "
		}
		s += fmt.Sprintf("%.1f day", ts.Days)
		if ts.Days >= 2.0 {
			s += "s"
		}
	}
	if ts.Hours > 0.0 {
		if s != "" {
			s += " "
		}
		s += fmt.Sprintf("%.1f hour", ts.Hours)
		if ts.Hours >= 2.0 {
			s += "s"
		}
	}
	if ts.Minutes > 0.0 {
		if s != "" {
			s += " "
		}
		s += fmt.Sprintf("%.1f minute", ts.Minutes)
		if ts.Minutes >= 2.0 {
			s += "s"
		}
	}

	return s
}

func (ts *TimeSpent) ToMinutes() uint64 {
	minutes := 0.0
	minutes += ts.Minutes
	minutes += ts.Hours * 60.0
	minutes += ts.Days * 8.0 * 60.0
	minutes += ts.Weeks * 5.0 * 8.0 * 60.0
	minutes += ts.Months * 4.0 * 5.0 * 8.0 * 60.0

	return uint64(minutes)
}

func (ts *TimeSpent) Add(other *TimeSpent) *TimeSpent {
	ts.Months += other.Months
	ts.Weeks += other.Weeks
	ts.Days += other.Days
	ts.Hours += other.Hours
	ts.Minutes += other.Minutes

	return ts
}

var sp = "^/spen[dt]\\s+"
var fl = "[0-9]+[.]?[0-9]*|[0-9]*[.]?[0-9]+"
var mi = "(?P<minutes>" + fl + ")\\s*(mi?|mins?|minutes?)?\\s*"
var ho = "(?P<hours>" + fl + ")\\s*(ho?|hours?)\\s*"
var da = "(?P<days>" + fl + ")\\s*(da?|days?)\\s*"
var we = "(?P<weeks>" + fl + ")\\s*(we?|weeks?)\\s*"
var mo = "(?P<months>" + fl + ")\\s*(mo|months?)\\s*"
var miP = "(" + mi + ")?"
var hoP = "(" + ho + ")?"
var daP = "(" + da + ")?"
var weP = "(" + we + ")?"
var moP = "(" + mo + ")?"

var spentAllRegex = regexp.MustCompile(sp + moP + weP + daP + hoP + miP)

// Keep these sorted by decreasing priority, since first match breaks.
var expressions = []*regexp.Regexp{
	spentAllRegex,
}

// CollectTimeSpent returns the TimeSpent that was collected from the message
// It reads the Gitlab /spend or /spent commands.
// Available time units: https://docs.gitlab.com/ee/user/project/time_tracking.html#available-time-units
// If no time unit is specified, minutes are assumed.
func CollectTimeSpent(message string) *TimeSpent {
	ts := &TimeSpent{}
	lines := strings.Split(message, "\n")

	for _, line := range lines {
		lineTs := extractTimeSpentFromLine(strings.TrimSpace(line))
		if lineTs == nil {
			continue
		}

		ts.Add(lineTs)
	}

	return ts
}

func extractTimeSpentFromLine(line string) *TimeSpent {
	for _, expression := range expressions {
		ts := extractTimeSpentUsingRegexp(line, expression)
		if ts != nil {
			return ts
		}
	}

	return nil
}

func extractTimeSpentUsingRegexp(line string, r *regexp.Regexp) *TimeSpent {
	matches := r.FindStringSubmatch(line)
	if len(matches) == 0 {
		return nil
	}

	months := extractTimeComponent(matches, r, "months")
	weeks := extractTimeComponent(matches, r, "weeks")
	days := extractTimeComponent(matches, r, "days")
	hours := extractTimeComponent(matches, r, "hours")
	minutes := extractTimeComponent(matches, r, "minutes")

	return &TimeSpent{
		Months:  months,
		Weeks:   weeks,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
	}
}

func extractTimeComponent(matches []string, r *regexp.Regexp, component string) float64 {
	componentIndex := r.SubexpIndex(component)
	componentString := "0"
	if componentIndex != -1 {
		if matches[componentIndex] != "" {
			componentString = matches[componentIndex]
		}
	}
	componentFloat, err := strconv.ParseFloat(componentString, 64)
	if err != nil {
		// this should never happen unless we fiddle with and break our regexes
		fmt.Println("cannot parse", component, componentString, r.String())
		return 0
	}

	return componentFloat
}
