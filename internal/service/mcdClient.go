package service

import (
	"context"

	"github.com/iantal/dta/internal/domain"
	mcdprotos "github.com/iantal/mcd/protos/mcd"
	"github.com/sirupsen/logrus"
)

func (e *Explorer) sendMCDRequest(project *domain.Project) {
	projectID := project.ProjectID.String()
	commit := project.CommitHash
	e.log.WithFields(logrus.Fields{
		"projectID": project.ProjectID,
		"commit": project.CommitHash,
	}).Info("Sending request to mcd")
	libraries := e.libraryDB.GetLibraryByIDAndCommit(projectID, commit)

	libs := []*mcdprotos.Library{}
	for _, l := range libraries {
		lib := &mcdprotos.Library{
			Name: l.Name,
			Type: l.Packaging,
		}
		libs = append(libs, lib)
	}

	r := &mcdprotos.DownloadRequest{
		ProjectID:  projectID,
		CommitHash: commit,
		Libraries:  libs,
	}

	resp, err := e.mcdClient.Download(context.Background(), r)
	if err != nil {
		e.log.WithFields(logrus.Fields{
			"projectID": project.ProjectID,
			"commit": project.CommitHash,
			"error": err,
		}).Error("Failed to get data from mcd")
		e.db.UpdateProject(project, domain.McdFailure)
		return
	}

	if resp.CommitHash == commit && resp.ProjectID == projectID {
		e.db.UpdateProject(project, domain.McdSuccess)
	} else {
		e.log.WithFields(logrus.Fields{
			"projectID": project.ProjectID,
			"commit": project.CommitHash,
			"response": resp,
		}).Warn("Response from mcd is invalid")
		e.db.UpdateProject(project, domain.McdFailure)
	}
}
