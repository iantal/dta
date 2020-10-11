package domain

// Status represents the status of a job
type Status string

const (
	DownloadSuccess       Status = "DOWNLOAD_SUCCESS"
	DownloadFailure       Status = "DOWNLOAD_FAILURE"
	UnbundleSuccess       Status = "UNBUNDLE_SUCCESS"
	UnbundleFailure       Status = "UNBUNDLE_FAILURE"
	BuildToolSuccess      Status = "BUILD_TOOL_SUCCESS"
	BuildToolFailure      Status = "BUILD_TOOL_FAILURE"
	DependencyTreeSuccess Status = "DEPENDENCY_TREE_SUCCESS"
	DependencyTreeFailure Status = "DEPENDENCY_TREE_FAILURE"
	ParseSuccess          Status = "PARSE_SUCCESS"
	ParseFailure          Status = "PARSE_FAILURE"
)

func (s Status) String() string {
	return string(s)
}
