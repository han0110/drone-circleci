# Drone CircleCI

[![Build Status](https://cloud.drone.io/api/badges/han0110/drone-circleci/status.svg?ref=refs/heads/master)](https://cloud.drone.io/han0110/drone-circleci)

Drone plugin for CircleCI integration

## Usage

### Wait

```shell
docker run --rm \
  -e PLUGIN_ACTION=wait \
  -e PLUGIN_API_TOKEN=${PLUGIN_API_TOKEN} \
  -e DRONE_REPO_LINK=${DRONE_REPO_LINK} \
  -e DRONE_COMMIT_SHA=${DRONE_COMMIT_SHA} \
  -e DRONE_COMMIT_BRANCH=${DRONE_COMMIT_BRANCH} \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  han0110/docker-circleci
```
