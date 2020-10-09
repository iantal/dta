package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	btdprotos "github.com/iantal/btd/protos/btd"
	"github.com/iantal/dta/internal/files"
	"github.com/iantal/dta/internal/service"
	protos "github.com/iantal/dta/protos/dta"
)

type CommitExplorer struct {
	log       hclog.Logger
	explorer *service.Explorer
}

// NewCommitExplorer creates a new Analyzer
func NewCommitExplorer(l hclog.Logger, basePath string, btdClient btdprotos.UsedBuildToolsClient, rmHost string, store files.Storage) *CommitExplorer {
	return &CommitExplorer{l, service.NewExplorer(l, basePath, btdClient, rmHost, store)}
}

// ExploreCommit performs the analysis for a specific commitId and projectId 
func (c *CommitExplorer) ExploreCommit(ctx context.Context, rr *protos.ExploreCommitRequest) (*protos.ExploreCommitResponse, error) {
	c.log.Info("Handle request for EcploreCommit", "projectID", rr.GetProjectID(), "commit", rr.GetCommitHash())

	err := c.explorer.Explore(rr.GetProjectID(), rr.GetCommitHash())
	if err != nil {
		return nil, err
	}

	return &protos.ExploreCommitResponse{}, nil
}
