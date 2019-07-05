package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("COMPLIANCE_CHECKER")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
		Organization struct {
			Name                            githubv4.String
			RequiresTwoFactorAuthentication githubv4.Boolean
		} `graphql:"organization(login: \"esdc-edsc\")"`
	}

	// Use client...
	err := client.Query(context.Background(), &query, nil)

	if err != nil {
		log.Fatal(err)
		// Handle error.
	}
	fmt.Println("    Login:", query.Viewer.Login)
	fmt.Println("CreatedAt:", query.Viewer.CreatedAt)
	fmt.Println("Org Name: ", query.Organization.Name)
	fmt.Println("2fa Required: ", query.Organization.RequiresTwoFactorAuthentication)

}
