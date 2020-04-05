package circleci

import (
	"context"
	"strconv"
	"time"
)

type (
	PipelineError struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	PipelineTrigger struct {
		Type  string `json:"type"`
		Actor struct {
			Login     string `json:"login"`
			AvatarURL string `json:"avatar_url"`
		} `json:"actor"`
		ReceivedAt time.Time `json:"received_at"`
	}
	PipelineVCS struct {
		OriginRepositoryURL string `json:"origin_repository_url"`
		TargetRepositoryURL string `json:"target_repository_url"`
		Revision            string `json:"revision"`
		ProviderName        string `json:"provider_name"`
		Commit              struct {
			Body    string `json:"body"`
			Subject string `json:"subject"`
		} `json:"commit"`
		Branch string `json:"branch"`
	}
)

// Pipeline defines struct for item of api response of pipeline.
type Pipeline struct {
	Number      uint            `json:"number"`
	ID          string          `json:"id"`
	ProjectSlug string          `json:"project_slug"`
	State       string          `json:"state"`
	Errors      []PipelineError `json:"errors"`
	VCS         PipelineVCS     `json:"vcs"`
	Trigger     PipelineTrigger `json:"trigger"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedAt   time.Time       `json:"created_at"`
}

// pipelineResponse defines struct of single pipeline response.
type pipelineResponse Pipeline

// GetPipeline get pipeline by number
func (c *Client) GetPipeline(ctx context.Context, projectSlug string, pipelineNumber uint) (Pipeline, error) {
	var res pipelineResponse
	pathParams := map[string]string{
		"projectSlug":    projectSlug,
		"pipelineNumber": strconv.Itoa(int(pipelineNumber)),
	}
	if err := c.get(
		&res, APIPath[c.apiVersion]["pipeline"],
		withPathParams(pathParams),
	); err != nil {
		return Pipeline{}, err
	}
	return Pipeline(res), nil
}

// Pipelines defines plural of Pipeline.
type Pipelines []Pipeline

// PipelineListIterator extends listIterator.
type PipelineListIterator struct {
	listIterator
	projectSlug string
	branch      string
}

// NewPipelineListIter create a instance of pipeline list iterator.
func (c *Client) NewPipelineListIter(projectSlug string) *PipelineListIterator {
	client := c.C()
	client.client.SetPathParams(map[string]string{"projectSlug": projectSlug})
	return &PipelineListIterator{
		listIterator: newListIterator(client, APIPath[c.apiVersion]["pipelines"]),
		projectSlug:  projectSlug,
	}
}

// SetBranch set branch query parameter for iterator.
func (i *PipelineListIterator) SetBranch(branch string) *PipelineListIterator {
	i.branch = branch
	i.extraOptions = append(i.extraOptions, withQueryParams(map[string]string{
		"branch": branch,
	}))
	return i
}

// Next get data of next page.
func (i *PipelineListIterator) Next(ctx context.Context, result *Pipelines) bool {
	return i.listIterator.Next(ctx, result)
}

// All get data of all the pages.
func (i *PipelineListIterator) All(ctx context.Context, result *Pipelines) error {
	return i.listIterator.All(ctx, result)
}
