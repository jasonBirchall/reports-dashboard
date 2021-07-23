package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v35/github"
)

func main() {
	client := *github.NewClient(nil)
	repoBlob := "cloud-platform"
	org := "ministryofjustice"

	ctx := context.Background()
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

	for _, repo := range list {
		fmt.Println(*repo.FullName)
	}
}
