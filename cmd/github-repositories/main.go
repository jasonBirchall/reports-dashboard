package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	str, _ := fetchRepoDescription("ministryofjustice", "cloud-platform-environments")
	fmt.Println(str)
}

// fetchRepoDescription fetches description of repo with owner and name.
func fetchRepoDescription(owner, name string) (string, error) {
	var q struct {
		Repository struct {
			Name  string
			Url   string
			Owner struct {
				Login string
			}
			DefaultBranchRef struct {
				Name string
			}
			BranchProtectionRules struct{} `graphql:"branchProtectionRules(first: 50)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_OAUTH_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	err := client.Query(context.Background(), &q, variables)
	return q.Repository.Owner.Login, err
}
