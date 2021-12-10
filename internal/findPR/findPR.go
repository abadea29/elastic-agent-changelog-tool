// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package findPR

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/google/go-github/v32/github"
)

// For the case when a commit SHA belongs to multiple PRs, we will choose the PR that was merged first
func Find(githubClient *github.Client, owner string, repo string, commitHash string) (*int, error) {
	rawListOfPRs, _, err := githubClient.PullRequests.ListPullRequestsWithCommit(context.Background(), owner, repo, commitHash, nil)
	listOfPRs := []*github.PullRequest{}

	for _, pr := range rawListOfPRs {
		if pr.MergedAt != nil {
			listOfPRs = append(listOfPRs, pr)
		}
	}

	sort.Slice(listOfPRs, func(i, j int) bool {
		return listOfPRs[i].MergedAt.Before(*listOfPRs[j].MergedAt)
	})

	return getOriginalPR(listOfPRs[0], githubClient, owner, repo).Number, err
}

func getOriginalPR(childPR *github.PullRequest, githubClient *github.Client, owner string, repo string) *github.PullRequest {
	if !containsBackportLabel(childPR.Labels) {
		return childPR
	} else {
		return getOriginalPR(getParentPR(childPR, githubClient, owner, repo), githubClient, owner, repo)
	}
}

func containsBackportLabel(labels []*github.Label) bool {
	const backportLabel = "backport"
	for _, label := range labels {
		if *label.Name == backportLabel {
			fmt.Println("PR contains label 'backport'")
			return true
		}
	}

	return false
}

// Making use of heuristics to find the 'parent' PR
func getParentPR(childPR *github.PullRequest, githubClient *github.Client, owner string, repo string) *github.PullRequest {
	var isAMatch bool
	var parentPR *github.PullRequest
	var prNumber int
	var regexMatch string

	isAMatch, _ = regexp.MatchString(`backport #\d+`, *childPR.Title)
	if isAMatch {
		regexMatch = regexp.MustCompile(`backport #\d+`).FindString(*childPR.Title)
	}

	isAMatch, _ = regexp.MatchString(`Cherry-pick #\d+`, *childPR.Title)
	if isAMatch {
		regexMatch = regexp.MustCompile(`Cherry-pick #\d+`).FindString(*childPR.Title)
	}

	isAMatch, _ = regexp.MatchString(`cherry-pick of #\d+`, *childPR.Body)
	if isAMatch {
		regexMatch = regexp.MustCompile(`cherry-pick of #\d+`).FindString(*childPR.Body)
	}

	isAMatch, _ = regexp.MatchString(`Cherry-pick of PR #\d+`, *childPR.Body)
	if isAMatch {
		regexMatch = regexp.MustCompile(`Cherry-pick of PR #\d+`).FindString(*childPR.Body)
	}

	isAMatch, _ = regexp.MatchString(`cherry-pick of PR [a-z]+#\d+`, *childPR.Body)
	if isAMatch {
		regexMatch = regexp.MustCompile(`cherry-pick of PR [a-z]+#\d+`).FindString(*childPR.Body)
	}

	prNumber, _ = strconv.Atoi(regexp.MustCompile(`\d+`).FindString(regexMatch))
	parentPR, _, _ = githubClient.PullRequests.Get(context.Background(), owner, repo, prNumber)
	return parentPR
}
