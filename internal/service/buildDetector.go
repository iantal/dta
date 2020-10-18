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
			e.handleGradle(destPath, commit, project)
			e.db.UpdateProject(project, domain.DependencyTreeSuccess)
			e.sendMCDRequest(project)
		} else if bt == "maven" {
			e.db.UpdateProjectBuildToolAndStatus(project, domain.BuildToolSuccess, "maven")
		} else {
			fmt.Println("UNKNOWN")
		}
	}
}

func (e *Explorer) handleGradle(destPath, commit string, project *domain.Project) {
	gw := NewGradlew(e.log, filepath.Join(destPath, commit))

	projects, err := gw.Projects()
	if err != nil {
		e.db.UpdateProject(project, domain.DependencyTreeFailure)
		e.log.Error("Error generating gradle dependency tree")
		return
	}

	fmt.Println(projects)

	if len(projects) == 1 {
		e.db.UpdateProjectName(project, projects[0])
		tree := gw.Dependencies(projects[0], false)
		results, err := e.parseGradle(tree)
		if err != nil {
			e.db.UpdateProject(project, domain.DependencyTreeFailure)
			e.log.Error("Could not generate dependency tree for " + projects[0])
			return
		}
		for _, proj := range results {
			for _, library := range proj.Libraries {
				e.libraryDB.AddLibrary(domain.NewLibrary(project.ProjectID, commit, library.Name, library.Type, library.Scope, "JAR"))
			}
		}
	} else {
		for _, p := range projects {
			subproject := domain.NewProject(project.ProjectID, project.CommitHash, p, project.Status, project.BuildTool, project.Data)
			e.db.AddProject(subproject)
			tree := gw.Dependencies(p, true)
			results, err := e.parseGradle(tree)
			if err != nil {
				e.db.UpdateProject(subproject, domain.DependencyTreeFailure)
				e.log.Error("Could not generate dependency tree for " + p)
				return
			}
			for _, proj := range results {
				for _, library := range proj.Libraries {
					e.libraryDB.AddLibrary(domain.NewLibrary(project.ProjectID, commit, library.Name, library.Type, library.Scope, "JAR"))
				}
			}
		}
	}
}
