package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func validateAPIKey() (string, error) {
	key := os.Getenv("GITHUB_API_TOKEN")
	if key == "" {
		return key, fmt.Errorf("Missing API Token!")
	}

	return key, nil
}

func getGithubIssues(key string) ([]github.Issue, error) {
	var issues []github.Issue
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

	for {
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

	return issues, nil
}

func getClosedGithubIssues(key string) ([]github.Issue, error) {
	var issues []github.Issue
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

	for {
		sresults, resp, err := client.Search.Issues(ctx, fmt.Sprintf("state:closed %s closed:>=2019-11-01", repoStr), sopt)
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

	return issues, nil
}
