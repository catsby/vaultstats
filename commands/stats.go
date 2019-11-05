package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/mitchellh/cli"
	"golang.org/x/oauth2"
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
	key := os.Getenv("GITHUB_API_TOKEN")
	if key == "" {
		c.UI.Error("Missing API Token!")
		return 1
	}

	n := time.Now()
	c.UI.Output(n.Format(time.RFC1123))
	c.UI.Output("Collecting stats...")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: key},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// by default, only show issues
	repoFilter := []string{
		"hashicorp/vault",
		// "hashicorp/vault-plugin-auth-kubernetes",
	}

	// query and count issues/prs
	sopt := &github.SearchOptions{Sort: "updated"}

	repoStr := "repo:"
	repoStr = repoStr + strings.Join(repoFilter, " repo:")

	var issues []github.Issue

	for {
		// sresults, resp, err := client.Search.Issues(ctx, fmt.Sprintf("state:open no:label %s %s", s, filter), sopt)
		sresults, resp, err := client.Search.Issues(ctx, fmt.Sprintf("state:open %s", repoStr), sopt)
		if err != nil {
			log.Printf("Error Searching Issues: %s", err)
			break
		}
		issues = append(issues, sresults.Issues...)
		if resp.NextPage == 0 {
			break
		}
		sopt.Page = resp.NextPage
	}

	c.UI.Output(fmt.Sprintf("Total count: %d\n", len(issues)))

	var issueCount, prCount, unlabeled int
	for _, i := range issues {
		if len(i.Labels) == 0 {
			unlabeled++
		}

		if i.PullRequestLinks != nil {
			prCount++
			continue
		}
		issueCount++
	}
	c.UI.Output(fmt.Sprintf("  Issue count: %d", issueCount))
	c.UI.Output(fmt.Sprintf("  PR count: %d", prCount))
	c.UI.Output(fmt.Sprintf("  Unlabeled count: %d", unlabeled))

	// query closed since November 1, 2019
	str := fmt.Sprintf("state:closed %s closed:>=2019-11-01", repoStr)
	sresults, resp, err := client.Search.Issues(ctx, fmt.Sprintf("state:closed %s closed:>=2019-11-01", repoStr), sopt)
	if err != nil {
		log.Printf("Error Searching closed Issues: %s", err)
	}

	c.UI.Output(fmt.Sprintf("\n  Closed since November 1, 2019: %d\n", sresults.Total))

	return 0
}
