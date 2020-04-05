package circleci

// HeaderAPIToken defines key of circle api token in header.
const HeaderAPIToken = "Circle-Token"

// APIEndpoint defines circleci api endpoint.
const APIEndpoint = "https://circleci.com/api"

// APIVersion provides enum of circleci api version.
type APIVersion string

const (
	V1 APIVersion = "v1.1"
	V2 APIVersion = "v2"
)

// DefaultAPIVersion is currently v2.
const DefaultAPIVersion = V2

// IsSupported return whether this version is supported.
func (v APIVersion) IsSupported() bool {
	switch v {
	case V2:
		return true
	case V1:
		// TODO: Remove me until v1 implemented.
		return false
	}
	return false
}

// APIPath defines map of api path by api version.
var APIPath = map[APIVersion]map[string]string{
	V2: {
		"me":               "/me",
		"collaborations":   "/me/collaborations",
		"pipeline":         "/project/{projectSlug}/pipeline/{pipelineNumber}",
		"pipelines":        "/project/{projectSlug}/pipeline",
		"pipelinesofmine":  "/project/{projectSlug}/pipeline/mine",
		"pipelineconfig":   "/pipeline/{pipelineID}/config",
		"pipelineworkflow": "/pipeline/{pipelineID}/workflow",
	},
}

// ItemsNumPerPage defines number of items on each iteration of api response
// which is constant and not abled to configured currently.
const ItemsNumPerPage = 20

// NullPageToken provide a magic value to tell iterator to stop.
const NullPageToken = "NULL"
