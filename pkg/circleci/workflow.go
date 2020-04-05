package circleci

import (
	"context"
	"regexp"
	"time"
)

// WorkflowStatus defines workflow status enum.
type WorkflowStatus string

const (
	WFSRunning  WorkflowStatus = "running"
	WFSSuccess  WorkflowStatus = "success"
	WFSFailed   WorkflowStatus = "failed"
	WFSCanceled WorkflowStatus = "canceled"
	WFSOnHold   WorkflowStatus = "on_hold"
)

// Workflow defines struct for item of api response of workflow.
type Workflow struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Status         WorkflowStatus `json:"status"`
	WorkflowID     string         `json:"Workflow_id"`
	PipelineNumber uint           `json:"pipeline_number"`
	ProjectSlug    string         `json:"project_slug"`
	CreatedAt      time.Time      `json:"created_at"`
	StoppedAt      time.Time      `json:"stopped_at"`
}

// Workflows defines plural of Workflow.
type Workflows []Workflow

// WorkflowFilterFn deinfes type of filter function of workflows.
type WorkflowFilterFn func(Workflow) bool

// FilterByFn filter out those whose return value by WorkflowFilterFn is false.
func (s Workflows) FilterByFn(fn WorkflowFilterFn) Workflows {
	workflows := make(Workflows, len(s))
	num := 0
	for _, w := range s {
		if fn(w) {
			workflows[num] = w
			num++
		}
	}
	return workflows[:num]
}

// FilterByName filter by workflow's name using regular expression.
func (s Workflows) FilterByName(nameRegexp *regexp.Regexp) Workflows {
	return s.FilterByFn(func(w Workflow) bool {
		return nameRegexp.Match([]byte(w.Name))
	})
}

// WorkflowListIterator extends listIterator.
type WorkflowListIterator struct {
	listIterator
	pipelineID string
}

// NewWorkflowListIter create a instance of Workflow list iterator.
func (c *Client) NewWorkflowListIter(pipelineID string) *WorkflowListIterator {
	client := c.C()
	client.client.SetPathParams(map[string]string{"pipelineID": pipelineID})
	return &WorkflowListIterator{
		listIterator: newListIterator(client, APIPath[c.apiVersion]["pipelineworkflow"]),
		pipelineID:   pipelineID,
	}
}

// Next get data of next page.
func (i *WorkflowListIterator) Next(ctx context.Context, result *Workflows) bool {
	return i.listIterator.Next(ctx, result)
}

// All get data of all the pages.
func (i *WorkflowListIterator) All(ctx context.Context, result *Workflows) error {
	return i.listIterator.All(ctx, result)
}
