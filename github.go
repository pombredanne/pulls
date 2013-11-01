package pulls

import (
	"fmt"
	gh "github.com/crosbymichael/octokat"
	"strconv"
)

// Top level type that manages a repository
type Maintainer struct {
	repo   gh.Repo
	client *gh.Client
}

func NewMaintainer(client *gh.Client, org, repo string) (*Maintainer, error) {
	return &Maintainer{
		repo:   gh.Repo{Name: repo, UserName: org},
		client: client,
	}, nil
}

func (m *Maintainer) Repository() (*gh.Repository, error) {
	return m.client.Repository(m.repo, nil)
}

// Return all pull requests
func (m *Maintainer) GetPullRequests(state string) ([]*gh.PullRequest, error) {
	o := &gh.Options{}
	o.QueryParams = map[string]string{
		"state":    state,
		"per_page": "100",
	}
	prevSize := -1
	page := 1
	allPRs := []*gh.PullRequest{}
	for len(allPRs) != prevSize {
		o.QueryParams["page"] = strconv.Itoa(page)
		if prs, err := m.client.PullRequests(m.repo, o); err != nil {
			return nil, err
		} else {
			prevSize = len(allPRs)
			allPRs = append(allPRs, prs...)
			page += 1
		}
		fmt.Printf(".")
	}
	return allPRs, nil
}

// Return a single pull request
func (m *Maintainer) GetPullRequest(number string) (*gh.PullRequest, error) {
	return m.client.PullRequest(m.repo, number, nil)
}

// Return all comments for an issue or pull request
func (m *Maintainer) GetComments(number string) ([]gh.Comment, error) {
	return m.client.Comments(m.repo, number, nil)
}

type ShowFilters struct {
	NoMerge  bool
	FromUser string
}

// Filter pull requests
func (m *Maintainer) FilterPullRequests(prs []*gh.PullRequest, filters *ShowFilters) ([]*gh.PullRequest, error) {
	out := []*gh.PullRequest{}
	for _, pr := range prs {
		fullPr, err := m.GetPullRequest(strconv.Itoa(pr.Number))
		fmt.Printf(".")
		if err != nil {
			return nil, err
		}
		if (!filters.NoMerge || (filters.NoMerge && !fullPr.Mergeable)) &&
			(filters.FromUser == "" || (filters.FromUser == fullPr.User.Login)) {
			out = append(out, fullPr)
		}
	}
	return out, nil
}

// Add a comment to an existing pull request
func (m *Maintainer) AddComment(number, comment string) (gh.Comment, error) {
	return m.client.AddComment(m.repo, number, comment)
}

// Merge a pull request
func (m *Maintainer) MergePullRequest(number, comment string) (gh.Merge, error) {
	o := &gh.Options{}
	o.Params = map[string]string{
		"commit_message": comment,
	}
	return m.client.MergePullRequest(m.repo, number, o)
}
