package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func main() {
	const (
		repoBlob = "cloud-platform"
		org      = "ministryofjustice"
	)
	// GitHub Client creation
	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if os.Getenv("GITHUB_OAUTH_TOKEN") == "" {
		log.Println("you must have the GITHUB_OAUTH_TOKEN env var")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		Sort:        "full_name",
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allRepos []*github.Repository

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			fmt.Println("Fail:", err)
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var list []*github.Repository
	for _, repo := range allRepos {
		if strings.Contains(*repo.FullName, repoBlob) {
			list = append(list, repo)
		}
	}

	for _, l := range list {
		fmt.Println(*l.Name)
	}
}
