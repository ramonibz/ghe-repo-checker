package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"fmt"
	"os"
	"ghe-repo-checker/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
)

func main() {

	createDb :=false

	if len(os.Args) > 1 {
		if "initdb" == os.Args[1] {
			createDb = true
		}else {
			panic("Unexpected argument, only valid argument is: 'initdb' an it's used to initializa the db")
		}
	}

	gheRepos := collectAllRepos()

	if createDb {
		fmt.Println("Initializing db with all repos. .  .   .")
		aws.CreateItems(gheRepos)
	}

	dbReposMap := aws.ScanTable("ghe-repositories")

	deletedRepos := getDeletedRepos(gheRepos, dbReposMap)
	newRepos := getNewRepos(gheRepos, dbReposMap)

	if len(deletedRepos) > 0 {
		fmt.Println("Deleted")
		fmt.Println(deletedRepos)
		aws.SendNotification("Deleted MS Repos", strings.Join(deletedRepos[:],","))
	}

	if len(newRepos) > 0 {
		fmt.Println("Created")
		fmt.Println(newRepos)
		aws.SendNotification("Created New MS Repos", strings.Join(newRepos[:],","))
	}
	fmt.Println("Done!")
	//dynamo.ListAllTables()
}


func collectAllRepos()map[string]github.Repository {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GHE_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	m := make(map[string]github.Repository)
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg("scmspain", opt)
		if err != nil {
			panic(err)
		}

		for _, element := range repos {
			m[*element.Name] = *element
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	return m
}

func getDeletedRepos(gheRepos map[string]github.Repository, dbRepos map[string]map[string]*dynamodb.AttributeValue) []string{
	var retVal []string

	for _, element := range dbRepos {
		var found = false
		for _, subElement := range gheRepos {
			if *subElement.Name == *element["Name"].S {
				found = true
			}
		}
		if !found {
			retVal = append(retVal, *element["Name"].S)
		}
	}

	return retVal
}

func getNewRepos(gheRepos map[string]github.Repository, dbRepos map[string]map[string]*dynamodb.AttributeValue) []string {
	var retVal []string

	for _, element := range gheRepos {
		var found = false
		if !strings.HasPrefix(*element.Name, "ms-") {
			continue
		}
		for _, subElement := range dbRepos {
			//log.Println("------------")
			//log.Println(*subElement.Name)
			//log.Println(*element["Name"].S)
			//log.Println("------------")
			if *subElement["Name"].S == *element.Name {
				found = true
			}
		}
		if !found {
			retVal = append(retVal, *element.Name)
		}
	}

	return retVal
}