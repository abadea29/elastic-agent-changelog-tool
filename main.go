// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package main

import (
	"fmt"

	"github.com/elastic/elastic-agent-changelog-tool/internal/findPR"
	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/pkg/errors"
)

func main() {
	var repo string = "elastic/beats"
	var commitHash string = "191a0752b5ceddc7b7657a517b90ca76c1350f30"
	args := [...]string{repo, commitHash}
	fmt.Println("Find the original Pull Request")

	// Setup GitHub
	err := github.EnsureAuthConfigured()
	if err != nil {
		errors.Wrap(err, "GitHub auth configuration failed")
	}

	githubClient, err := github.Client()
	if err != nil {
		errors.Wrap(err, "creating GitHub client failed")
	}

	// GitHub user
	githubUser, err := github.User(githubClient)
	if err != nil {
		errors.Wrap(err, "fetching GitHub user failed")
	}
	fmt.Printf("Current GitHub user: %s\n", githubUser)

	if len(args) > 0 {
		repo = args[0]
	} else {
		repo = "elastic/beats"
	}
	fmt.Println("The repo is ", repo)
	fmt.Println("The commit hash we are looking for is", commitHash)

	// Find the original PR
	err = findPR.Find(githubUser, githubClient, repo, commitHash)
	if err != nil {
		errors.Wrap(err, "can't find the PR")
	}

	fmt.Println("Done")
}
