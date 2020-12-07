package service

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	btdprotos "github.com/iantal/btd/protos/btd"
	"github.com/iantal/dta/internal/domain"
	"github.com/iantal/dta/internal/files"
	"github.com/iantal/dta/internal/repository"
	"github.com/iantal/dta/internal/utils"
	gpprotos "github.com/iantal/dta/protos/gradle-parser"
	mcdprotos "github.com/iantal/mcd/protos/mcd"
	"github.com/sirupsen/logrus"
)

// Explorer is a service that handles the analysis flows
type Explorer struct {
	log          *utils.StandardLogger
	db           *repository.ProjectDB
	libraryDB    *repository.LibraryDB
	basePath     string
	btdClient    btdprotos.UsedBuildToolsClient
	gradleClient gpprotos.GradleParseServiceClient
	mcdClient    mcdprotos.DownloaderClient
	rmHost       string
	store        files.Storage
}

// NewExplorer creates an Explorer
func NewExplorer(l *utils.StandardLogger, db *repository.ProjectDB, ldb *repository.LibraryDB, basePath string, btdClient btdprotos.UsedBuildToolsClient, gradleClient gpprotos.GradleParseServiceClient, mcdClient mcdprotos.DownloaderClient, rmHost string, store files.Storage) *Explorer {
	return &Explorer{l, db, ldb, basePath, btdClient, gradleClient, mcdClient, rmHost, store}
}

// Explore performs the analysis steps for a given commit of a project
func (e *Explorer) Explore(projectID, commit string) {
	projectPath := filepath.Join(e.store.FullPath(projectID), "bundle")

	var project *domain.Project
	project = e.db.GetProjectByIDAndCommit(projectID, commit)
	if project != nil && project.Status == domain.DependencyTreeSuccess.String() {
		// send libraries for pID and commit to mcd for further analysis
		e.log.WithFields(logrus.Fields{
			"projectID": projectID,
			"commit":    commit,
		}).Info("Already anayzed. Sending data to mcd")
		e.sendMCDRequest(project)
		return
	}

	project = domain.NewProject(uuid.MustParse(projectID), commit, "root", domain.DownloadSuccess.String(), "", "")

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		err := e.downloadRepository(projectID, commit)
		if err != nil {
			e.log.WithFields(logrus.Fields{
				"projectID": projectID,
				"commit":    commit,
				"error":     err,
			}).Error("Could not download bundled repository")
			project.Status = domain.DownloadFailure.String()
			e.db.AddProject(project)
			return
		}
	}

	e.db.AddProject(project)
	bp := commit + ".bundle"
	srcPath := e.store.FullPath(filepath.Join(projectID, commit, "bundle", bp))
	destPath := e.store.FullPath(filepath.Join(projectID, commit, "unbundle"))

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err := e.store.Unbundle(srcPath, destPath)
		if err != nil {
			e.log.WithFields(logrus.Fields{
				"projectID": projectID,
				"commit":    commit,
				"error":     err,
			}).Error("Could not unbundle repository")
			e.db.UpdateProject(project, domain.UnbundleFailure)
			return
		}
	}
	e.db.UpdateProject(project, domain.UnbundleSuccess)

	go e.detectBuildTools(destPath, project)
}
