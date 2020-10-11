package service

import (
	"context"
	"fmt"
	"path/filepath"

	btdprotos "github.com/iantal/btd/protos/btd"
	"github.com/iantal/dta/internal/domain"
)

func (e *Explorer) detectBuildTools(destPath string, project *domain.Project) {

	projectID := project.ProjectID.String()
	commit := project.CommitHash

	r := &btdprotos.BuildToolRequest{
		ProjectID:  projectID,
		CommitHash: project.CommitHash,
	}

	resp, err := e.btdClient.GetBuildTools(context.Background(), r)
	if err != nil {
		e.log.Error("Could not get build tool for", "projectID", projectID, "commit", commit)
		e.db.UpdateProject(project, domain.BuildToolFailure)
		return
	}

	

	for _, bt := range resp.BuildTools {
		if bt == "gradle" {
			e.db.UpdateProjectBuildToolAndStatus(project, domain.BuildToolSuccess, "gradle")
			
			gw := NewGradlew(e.log, filepath.Join(destPath, commit))
			projects, err := gw.Projects()
			if err != nil {		
				e.db.UpdateProject(project, domain.DependencyTreeFailure)		
				e.log.Error("Error generating gradle dependency tree")
			}
			e.db.UpdateProject(project, domain.DependencyTreeSuccess)		
			fmt.Println(projects)
		} else if bt == "maven" {
			e.db.UpdateProjectBuildToolAndStatus(project, domain.BuildToolSuccess, "maven")
		} else {
			fmt.Println("UNKNOWN")
		}
	}
}
