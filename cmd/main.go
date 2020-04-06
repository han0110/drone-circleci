package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	version = "development"

	logger = logrus.WithField("prefix", "drone-circleci")
)

func main() {
	// Load env file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		if err := godotenv.Load(env); err != nil {
			logger.Warnf("failed to load env from file %s, err: %v", env, err)
		}
	}

	app := cli.NewApp()
	app.Name = "Drone CircleCI"
	app.Usage = "Drone plugin for CircleCI integration"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo",
			Usage:  "repo to take action",
			EnvVar: "PLUGIN_REPO,DRONE_REPO_LINK",
		},
		cli.StringFlag{
			Name:   "api-token",
			Usage:  "api token of circleci",
			EnvVar: "PLUGIN_API_TOKEN",
		},
		cli.StringFlag{
			Name:   "action",
			Usage:  "action to dispatch",
			EnvVar: "PLUGIN_ACTION",
			Value:  "wait",
		},
		cli.StringFlag{
			Name:   "wait.sha",
			Usage:  "commit sha to wait",
			EnvVar: "PLUGIN_WAIT_SHA,DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "wait.branch",
			Usage:  "commit branch to wait",
			EnvVar: "PLUGIN_WAIT_BRANCH,DRONE_SOURCE_BRANCH",
		},
		cli.StringFlag{
			Name:   "wait.workflow",
			Usage:  "workflows to wait (in regex format)",
			EnvVar: "PLUGIN_WAIT_WORKFLOW",
			Value:  ".+",
		},
		cli.Int64Flag{
			Name:   "wait.interval",
			Usage:  "interval to check status (in second)",
			EnvVar: "PLUGIN_WAIT_INTERVAL",
			Value:  15,
		},
	}

	app.Before = func(ctx *cli.Context) error {
		prefixedTextFormatter := new(prefixed.TextFormatter)
		prefixedTextFormatter.TimestampFormat = "2006-01-02 15:04:05"
		prefixedTextFormatter.FullTimestamp = true
		logrus.SetFormatter(prefixedTextFormatter)
		logrus.SetLevel(logrus.InfoLevel)
		if debug := os.Getenv("DEBUG"); debug == "true" {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			Repo:     c.String("repo"),
			APIToken: c.String("api-token"),
			Action:   Action(c.String("action")),
			Wait: WaitConfig{
				SHA:      c.String("wait.sha"),
				Branch:   c.String("wait.branch"),
				Workflow: c.String("wait.workflow"),
				Interval: time.Duration(c.Int64("wait.interval")) * time.Second,
			},
		},
	}
	return plugin.Exec()
}
