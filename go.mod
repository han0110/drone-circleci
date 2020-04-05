module github.com/han0110/drone-circleci

go 1.13

require (
	github.com/go-resty/resty/v2 v2.2.0
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.5.0
	github.com/urfave/cli v0.0.0-00010101000000-000000000000
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
)

replace github.com/urfave/cli => github.com/bradrydzewski/cli v1.19.2-0.20170424184348-0d51abd87c77
