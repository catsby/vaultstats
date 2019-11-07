package commands

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mitchellh/cli"
)

type CSVCommand struct {
	UI cli.Ui
}

func (c CSVCommand) Help() string {
	return strings.TrimSpace(`
Usage: vaultstats csv [options]
`)
}

func (c CSVCommand) Synopsis() string {
	return "Return a CSV formatted output of Vault issues"
}

func (c CSVCommand) Run(args []string) int {
	key, err := validateAPIKey()
	if err != nil {
		c.UI.Output(err.Error())
		return 1
	}

	issues, err := getGithubIssues(key)
	if err != nil {
		c.UI.Output(err.Error())
		return 1
	}

	if err := CreateCSV(issues); err != nil {
		return 1
	}

	return 0
}

func labelExists(issue github.Issue, validateExists string) bool {
	for _, l := range issue.Labels {
		if *l.Name == validateExists {
			return true
		}
	}

	return false
}

func CreateCSV(issues []github.Issue) error {
	records := [][]string{{"title", "enhancement", "bug", "metadata", "url"}}
	for _, i := range issues {
		enhancement := labelExists(i, "enhancement")
		bug := labelExists(i, "bug")

		records = append(records, []string{*i.Title, fmt.Sprintf("%v", enhancement), fmt.Sprintf("%v", bug), fmt.Sprintf("%#v", ghLabels(i.Labels)), *i.URL})
	}

	w := csv.NewWriter(os.Stdout)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()

	return w.Error()
}

// ghLabels is a convienence type so we can nicely print all the GitHub labels
// with GoStringer
type ghLabels []github.Label

func (ghl ghLabels) GoString() string {
	labels := make([]string, len(ghl))
	for i, l := range ghl {
		labels[i] = *l.Name
	}

	sort.Strings(labels)
	return strings.Join(labels, ", ")
}
