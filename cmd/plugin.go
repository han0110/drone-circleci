package main

import (
	"context"
	"regexp"
	"time"

	"github.com/han0110/drone-circleci/pkg/circleci"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Action defines supported action enum to dispatch
type Action string

const (
	// ActionWait will wait for pipeline's workflow to be success.
	ActionWait Action = "wait"
	// ActionTrigger will trigger a new pipeline.
	ActionTrigger Action = "trigger"
)

// IsValid check whether action is supported
func (a Action) IsValid() bool {
	switch a {
	case ActionWait:
		return true
	// TODO: Add ActionTrigger support
	case ActionTrigger:
		return false
	}
	return false
}

type (
	WaitConfig struct {
		SHA      string
		Branch   string
		Workflow string
		Interval time.Duration
	}
	Config struct {
		Repo     string
		APIToken string
		Action   Action
		Wait     WaitConfig
	}
	Plugin struct {
		Config Config
		client *circleci.Client
	}
)

func (p *Plugin) Exec() (err error) {
	logger.Debugf("plugin got config %+v", p.Config)

	if !p.Config.Action.IsValid() {
		return errors.Errorf("action %s is not supported yet", p.Config.Action)
	}

	if err = p.parseRepo(); err != nil {
		return errors.Wrap(err, "failed to parse repo")
	}
	if err = p.initClient(); err != nil {
		return errors.Wrap(err, "failed to init client")
	}

	// Wait for 5 second for unexpected time lag
	time.Sleep(5 * time.Second)

	switch p.Config.Action {
	case ActionWait:
		err = p.wait()
	default:
		panic(errors.Errorf("action %s not handled", p.Config.Action))
	}
	return errors.Wrapf(err, "failed to execute action %s", p.Config.Action)
}

var defaultRepoRegexp = regexp.MustCompile("(https://)(github|bitbucket)(.*)/(.+)/(.+)")

func (p *Plugin) parseRepo() error {
	switch {
	case defaultRepoRegexp.Match([]byte(p.Config.Repo)):
		p.Config.Repo = defaultRepoRegexp.ReplaceAllString(p.Config.Repo, "$2/$4/$5")
	default:
		return errors.Errorf("format of repo %s is not supported yet", p.Config.Repo)
	}
	return nil
}

func (p *Plugin) initClient() error {
	if p.Config.APIToken == "" {
		return errors.New("api token has to be set")
	}
	p.client = circleci.NewClient()
	if err := p.client.Authenticate(p.Config.APIToken); err != nil {
		return err
	}
	return nil
}

func (p *Plugin) wait() error {
	workflowNameRegexp, err := regexp.Compile(p.Config.Wait.Workflow)
	if err != nil {
		return errors.Wrapf(err, "failed to compile workflow name regexp %s", p.Config.Wait.Workflow)
	}
	pipelineIter := p.client.NewPipelineListIter(p.Config.Repo).SetBranch(p.Config.Wait.Branch)
	var pipelines circleci.Pipelines
	var targetPipeline *circleci.Pipeline
	for targetPipeline == nil {
		if ok := pipelineIter.Next(context.TODO(), &pipelines); !ok {
			if err := pipelineIter.Error(); err != nil {
				return errors.Wrap(err, "failed to get pipelines")
			}
			return errors.New("failed to find target pipeline")
		}
		for i := range pipelines {
			if pipelines[i].VCS.Revision == p.Config.Wait.SHA {
				targetPipeline = &pipelines[i]
				break
			}
		}
	}
	logger.WithFields(logrus.Fields{
		"number": targetPipeline.Number,
		"id":     targetPipeline.ID,
		"sha":    targetPipeline.VCS.Revision,
		"commit": targetPipeline.VCS.Commit.Subject,
	}).Info("target pipeline found")

	workflowIter := p.client.NewWorkflowListIter(targetPipeline.ID)
	var workflows circleci.Workflows
	for {
		allSuccess := true
		// Get all workflow of pipeline
		if err := workflowIter.All(context.TODO(), &workflows); err != nil {
			return errors.Wrapf(err, "failed to get all workflows of pipeline #%d", targetPipeline.Number)
		}
		// Filter desired workflow
		workflows = workflows.FilterByName(workflowNameRegexp)
		logger.Debugf("got workflows to check status, %+v", workflows)
	workflowsLoop:
		for i := range workflows {
			switch workflows[i].Status {
			case circleci.WFSSuccess:
				// Do nothing
			case circleci.WFSCanceled, circleci.WFSFailed:
				return errors.Errorf("pipeline has workflow %s in status %s", workflows[i].Name, workflows[i].Status)
			default:
				allSuccess = false
				logger.Infof("workflow %s is still in status %s", workflows[i].Name, workflows[i].Status)
				break workflowsLoop
			}
		}
		if allSuccess {
			break
		}
		// Sleep for checking interval
		time.Sleep(p.Config.Wait.Interval)
	}

	logger.Infof("all workflows are success")
	return nil
}
