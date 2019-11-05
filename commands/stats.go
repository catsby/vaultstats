package commands

import (
	"os"
	"strings"

	"github.com/mitchellh/cli"
)

// var wgPrs sync.WaitGroup

// var collaborators bool
// var all bool
// var includeUsers []string
// var filterUsers []string
// var tableFormat bool

// // StatReviewStatus maps to status, defined below
// type StatReviewStatus uint

// const (
// 	statusAll StatReviewStatus = iota
// 	statusWaiting
// 	statusComments
// 	statusChanges
// 	statusApproved
// )

// var filter StatReviewStatus

// vaultstats
//const vaultstatsID = 1836975

// vault team
const vaultstatsID = 1836984

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

	// 	ctx := context.Background()
	// 	ts := oauth2.StaticTokenSource(
	// 		&oauth2.Token{AccessToken: key},
	// 	)
	// 	tc := oauth2.NewClient(ctx, ts)
	// 	client := github.NewClient(tc)

	return 0
}
