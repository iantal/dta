package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	btdprotos "github.com/iantal/btd/protos/btd"
	"github.com/iantal/dta/internal/domain"
	"github.com/iantal/dta/internal/files"
	"github.com/iantal/dta/internal/repository"
	gpprotos "github.com/iantal/dta/protos/gradle-parser"
)

// Explorer is a service that handles the analysis flows
type Explorer struct {
	log          hclog.Logger
	db           *repository.ProjectDB
	libraryDB    *repository.LibraryDB
	basePath     string
	btdClient    btdprotos.UsedBuildToolsClient
	gradleClient gpprotos.GradleParseServiceClient
	rmHost       string
	store        files.Storage
}

// NewExplorer creates an Explorer
func NewExplorer(l hclog.Logger, db *repository.ProjectDB, ldb *repository.LibraryDB, basePath string, btdClient btdprotos.UsedBuildToolsClient, gradleClient gpprotos.GradleParseServiceClient, rmHost string, store files.Storage) *Explorer {
	return &Explorer{l, db, ldb, basePath, btdClient, gradleClient, rmHost, store}
}

// Explore performs the analysis steps for a given commit of a project
func (e *Explorer) Explore(projectID, commit string) error {
	projectPath := filepath.Join(e.store.FullPath(projectID), "bundle")

	var project *domain.Project
	project = e.db.GetProjectByIDAndCommit(projectID, commit)
	if project != nil && project.Status == domain.DependencyTreeSuccess.String() {
		// send request to gradle parser
		return nil
	}

	project = domain.NewProject(uuid.MustParse(projectID), commit, "root", domain.DownloadSuccess.String(), "", "")

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		err := e.downloadRepository(projectID, commit)
		if err != nil {
			e.log.Error("Could not download bundled repository", "projectID", projectID, "commit", commit, "err", err)
			project.Status = domain.DownloadFailure.String()
			e.db.AddProject(project)
			return fmt.Errorf("Could not download bundled repository for %s", projectID)
		}
	}

	e.db.AddProject(project)
	bp := commit + ".bundle"
	srcPath := e.store.FullPath(filepath.Join(projectID, "bundle", bp))
	destPath := e.store.FullPath(filepath.Join(projectID, "unbundle"))

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err := e.store.Unbundle(srcPath, destPath)
		if err != nil {
			e.log.Error("Could not unbundle repository", "projectID", projectID, "commit", commit, "err", err)
			e.db.UpdateProject(project, domain.UnbundleFailure)
			return fmt.Errorf("Could not unbundle repository for %s", projectID)
		}
	}
	e.db.UpdateProject(project, domain.UnbundleSuccess)

	go e.detectBuildTools(destPath, project)

	return nil
}
