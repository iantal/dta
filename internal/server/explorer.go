package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	btdprotos "github.com/iantal/btd/protos/btd"
	"github.com/iantal/dta/internal/files"
	"github.com/iantal/dta/internal/repository"
	"github.com/iantal/dta/internal/service"
	protos "github.com/iantal/dta/protos/dta"
	gpprotos "github.com/iantal/dta/protos/gradle-parser"
	mcdprotos "github.com/iantal/mcd/protos/mcd"
	"github.com/jinzhu/gorm"
)

type CommitExplorer struct {
	log      hclog.Logger
	explorer *service.Explorer
}

// NewCommitExplorer creates a new Analyzer
func NewCommitExplorer(l hclog.Logger, db *gorm.DB, basePath string, btdClient btdprotos.UsedBuildToolsClient, gradleClient gpprotos.GradleParseServiceClient, mcdClient mcdprotos.DownloaderClient, rmHost string, store files.Storage) *CommitExplorer {
	projectDB := repository.NewProjectDB(l, db)
	libraryDB := repository.NewLibraryDB(l, db)
	return &CommitExplorer{l, service.NewExplorer(l, projectDB, libraryDB, basePath, btdClient, gradleClient, mcdClient, rmHost, store)}
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