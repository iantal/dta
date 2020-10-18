package service

import (
	"context"

	"github.com/iantal/dta/internal/domain"
	mcdprotos "github.com/iantal/mcd/protos/mcd"
)

func (e *Explorer) sendMCDRequest(project *domain.Project) {
	projectID := project.ProjectID.String()
	commit := project.CommitHash

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
		e.db.UpdateProject(project, domain.McdFailure)
		return
	}

	if resp.CommitHash == commit && resp.ProjectID == projectID {
		e.db.UpdateProject(project, domain.McdSuccess)
	} else {
		e.db.UpdateProject(project, domain.McdSuccess)
	}
}
