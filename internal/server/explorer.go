package server

import (
	"context"

	btdprotos "github.com/iantal/btd/protos/btd"
	"github.com/iantal/dta/internal/files"
	"github.com/iantal/dta/internal/repository"
	"github.com/iantal/dta/internal/service"
	"github.com/iantal/dta/internal/utils"
	protos "github.com/iantal/dta/protos/dta"
	gpprotos "github.com/iantal/dta/protos/gradle-parser"
	mcdprotos "github.com/iantal/mcd/protos/mcd"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type CommitExplorer struct {
	log      *utils.StandardLogger
	explorer *service.Explorer
}

// NewCommitExplorer creates a new Analyzer
func NewCommitExplorer(l *utils.StandardLogger, db *gorm.DB, basePath string, btdClient btdprotos.UsedBuildToolsClient, gradleClient gpprotos.GradleParseServiceClient, mcdClient mcdprotos.DownloaderClient, rmHost string, store files.Storage) *CommitExplorer {
	projectDB := repository.NewProjectDB(l, db)
	libraryDB := repository.NewLibraryDB(l, db)
	return &CommitExplorer{l, service.NewExplorer(l, projectDB, libraryDB, basePath, btdClient, gradleClient, mcdClient, rmHost, store)}
}

// ExploreCommit performs the analysis for a specific commitId and projectId
func (c *CommitExplorer) ExploreCommit(ctx context.Context, rr *protos.ExploreCommitRequest) (*protos.ExploreCommitResponse, error) {
	c.log.WithFields(logrus.Fields{
		"projectID": rr.ProjectID,
		"commit":    rr.CommitHash,
	}).Info("Handle request for ExploreCommit")

	go func() {
		c.explorer.Explore(rr.GetProjectID(), rr.GetCommitHash())
	}()

	return &protos.ExploreCommitResponse{}, nil
}
