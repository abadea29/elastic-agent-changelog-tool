// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package findPR

import (
	"github.com/google/go-github/v32/github"
)

// /repos/{owner}/{repo}/commits/{ref}  or /repos/{owner}/{repo}/git/commits/{commit_sha} -> long hash
func Find(githubUser string, githubClient *github.Client, repo string, commitHash string) error {

	return nil

}
