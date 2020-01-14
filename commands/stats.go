package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/mitchellh/cli"
)

// vault team
// const vaultstatsID = 1836984

// StatsCommand command for querying Stats and status by team, person, etc
type StatsCommand struct {
	UI cli.Ui
}

// Help lists usage syntax
func (c StatsCommand) Help() string {
	helpText := `
Usage: vaultstats stats [options] 

Options:

	--summarize, -s        Count issues by parent group, e.g. "core" instead of "core/api"
	--bugs, -b        catagorize by bugs

`
	return strings.TrimSpace(helpText)
}

// Synopsis shows a synopsis in the top level help
func (c StatsCommand) Synopsis() string {
	return "List Stats for the hashicorp/vault repo"
}

// Run Stats query with args
func (c StatsCommand) Run(args []string) int {
	key, err := validateAPIKey()
	if err != nil {
		c.UI.Output(err.Error())
		return 1
	}
	var summarize, bugs bool
	if len(args) > 0 {
		for _, a := range args {
			if a == "--summarize" || a == "-s" {
				summarize = true
			}
			if a == "--bugs" || a == "-b" {
				bugs = true
			}
		}
	}
	_ = bugs

	n := time.Now()
	c.UI.Output(n.Format(time.RFC1123))
	c.UI.Output("Collecting stats...")

	var issues []github.Issue
	if bugs {
		issues, err = getGitHubIssuesByBugs(key)
	} else {
		issues, err = getGitHubIssues(key)
	}
	if err != nil {
		c.UI.Output(err.Error())
		return 1
	}

	c.UI.Output(fmt.Sprintf("Total count: %d\n", len(issues)))

	var issueCount, prCount, unlabeled int
	labelMap := make(map[string]int)
	for _, i := range issues {
		// add up the labels
		if len(i.Labels) == 0 {
			unlabeled++
		}

		for _, l := range i.Labels {
			key := *l.Name
			if summarize {
				parts := strings.Split(key, "/")
				if parts[0] != "version" {
					key = parts[0]
				}
			}
			labelMap[key]++
		}

		if i.PullRequestLinks != nil {
			prCount++
			continue
		}
		issueCount++
	}

	if !bugs {
		c.UI.Output(fmt.Sprintf("  Open Issue count: %d", issueCount))
		c.UI.Output(fmt.Sprintf("  Open PR count: %d", prCount))
		c.UI.Output(fmt.Sprintf("  Unlabeled count: %d\n", unlabeled))
	} else {
		c.UI.Output(fmt.Sprintf("  Open bug count: %d\n", issueCount))
		delete(labelMap, "bug")
	}

	// sort label names
	var labelNames []string
	for name := range labelMap {
		labelNames = append(labelNames, name)
	}

	sort.Strings(labelNames)
	c.UI.Output("Count by label:")
	for _, name := range labelNames {
		c.UI.Output(fmt.Sprintf("  %s: %d", name, labelMap[name]))
	}

	// query closed since November 1, 2019
	// clear the search options page, if set by the above search
	closedIssues, err := getClosedGithubIssues(key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("\n\nError getting issue count: %d\n", err))
	}
	c.UI.Output(fmt.Sprintf("\n\nClosed since November 1, 2019: %d\n", closedIssues))

	return 0
}
