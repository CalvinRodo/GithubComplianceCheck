package main

import (
	"context"
	"fmt"
	"log"
	"os"

	. "github.com/logrusorgru/aurora"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func getNumRepositories(client *githubv4.Client, vars map[string]interface{}) githubv4.Int {
	var query struct {
		Organization struct {
			Repositories struct {
				TotalCount githubv4.Int
			}
		} `graphql:"organization(login: $org)"`
	}

	err := client.Query(context.Background(), &query, vars)

	if err != nil {
		log.Fatal(err)
	}

	return query.Organization.Repositories.TotalCount

}

func getClient() *githubv4.Client {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("COMPLIANCE_CHECKER")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	return githubv4.NewClient(httpClient)
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Need organization name to validate")
	}

	org := os.Args[1]
	fmt.Println("Checking settings for: ", org)

	variables := map[string]interface{}{
		"org": githubv4.String(org),
	}

	client := getClient()

	numRepos := getNumRepositories(client, variables)
	fmt.Println("Number of Repositories found in organization: ", numRepos)

	var query struct {
		Organization struct {
			Name                            string
			RequiresTwoFactorAuthentication bool
		} `graphql:"organization(login: $org)"`
	}

	// Use client...
	err := client.Query(context.Background(), &query, variables)

	if err != nil {
		log.Fatal(err)
		// Handle error.
	}

	fmt.Println("Full Org Name: ", query.Organization.Name)
	if query.Organization.RequiresTwoFactorAuthentication {
		fmt.Println(Green("IA2 Success: Two Factor Authentication Required"))
	} else {
		fmt.Println(Red("IA2 Failed: Two Factgor Authentication not Required for members"))
		fmt.Println(Red("https://help.github.com/en/articles/requiring-two-factor-authentication-in-your-organization"))
	}

	variables = map[string]interface{}{
		"org":           githubv4.String(org),
		"commentCursor": (*githubv4.String)(nil), // Null after arg to get first page
	}

	var repoQuery struct {
		Organization struct {
			Repositories struct {
				Nodes struct {
					Name                  githubv4.String
					BranchProtectionRules struct {
						RequiredApprovingReviewCount githubv4.Int
						RequiresApprovingReviews     bool
						RequiresCommitSignatures     bool
						RestrictPushes               bool
						RequiresStrictStatusChecks   bool
					} `graphql:"branchProtectionRules(first: 100)"`
				}
				PageInfo struct {
					EndCursor   string
					HasNextPage bool
				}
			} `graphql:"repositories(first:1, after: $repoCursor)"`
		} `graphql:"organization(login: $org)"`
	}

	for {
		// Use client...
		err = client.Query(context.Background(), &repoQuery, variables)

		if err != nil {
			log.Fatal(err)
			// Handle error.
		}
		//Validate BranchProtectionRules
		fmt.Println(repoQuery.Organization.Repositories.Nodes.Name)
		if !repoQuery.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
	}
}
