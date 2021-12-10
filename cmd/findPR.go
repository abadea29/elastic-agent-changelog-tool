// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package cmd

import (
	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/elastic/elastic-agent-changelog-tool/internal/cobraext"
	"github.com/elastic/elastic-agent-changelog-tool/internal/findPR"
)

const (
	findPRLongDescription = `Use this command to find the original PR that included the commit in the repository.`
	owner                 = "elastic"
)

// Initial setup of the command
func setupFindPRCommand() *cobraext.Command {
	cmd := &cobra.Command{
		Use:   "find-pr",
		Short: "Find the original Pull Request",
		Long:  findPRLongDescription,
		RunE:  findPRCommandAction,
	}

	cmd.Flags().String(cobraext.RepositoryFlagName, "beats", cobraext.RepositoryFlagDescription)

	return cobraext.NewCommand(cmd, cobraext.ContextPackage)
}

// Command action
func findPRCommandAction(cmd *cobra.Command, args []string) error {
	cmd.Println("Find the original Pull Request")

	var commitHash string

	repo, err := cmd.Flags().GetString(cobraext.RepositoryFlagName)
	if err != nil {
		return errors.Wrapf(err, "can't read %s flag:", cobraext.RepositoryFlagName)
	}

	// Setup GitHub
	err = github.EnsureAuthConfigured()
	if err != nil {
		return errors.Wrap(err, "GitHub auth configuration failed")
	}

	githubClient, err := github.Client()
	if err != nil {
		return errors.Wrap(err, "creating GitHub client failed")
	}

	// GitHub user
	githubUser, err := github.User(githubClient)
	if err != nil {
		return errors.Wrap(err, "fetching GitHub user failed")
	}
	cmd.Printf("Current GitHub user: %s\n", githubUser)

	// Find the original PR
	originalPRNumber, err := findPR.Find(githubClient, owner, repo, commitHash)
	if err != nil {
		return errors.Wrap(err, "can't find the PR")
	} else {
		cmd.Printf("Original PR number is", originalPRNumber)
	}

	cmd.Println(commitHash, ":", originalPRNumber)
	cmd.Println("Done")
	return nil
}
