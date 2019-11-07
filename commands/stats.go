package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

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

	n := time.Now()
	c.UI.Output(n.Format(time.RFC1123))
	c.UI.Output("Collecting stats...")

	issues, err := getGithubIssues(key)
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
		} else {
			for _, l := range i.Labels {
				labelMap[*l.Name]++
			}
		}

		if i.PullRequestLinks != nil {
			prCount++
			continue
		}
		issueCount++
	}
	c.UI.Output(fmt.Sprintf("  Open Issue count: %d", issueCount))
	c.UI.Output(fmt.Sprintf("  Open PR count: %d", prCount))
	c.UI.Output(fmt.Sprintf("  Unlabeled count: %d\n", unlabeled))

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
	c.UI.Output(fmt.Sprintf("\n\nClosed since November 1, 2019: %d\n", len(closedIssues)))

	return 0
}
