package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	btdprotos "github.com/iantal/btd/protos/btd"
	"github.com/iantal/dta/internal/files"
)

type Explorer struct {
	log       hclog.Logger
	basePath string
	btdClient btdprotos.UsedBuildToolsClient
	rmHost    string
	store     files.Storage
}

func NewExplorer(l hclog.Logger, basePath string, btdClient btdprotos.UsedBuildToolsClient, rmHost string, store files.Storage) *Explorer {
	return &Explorer{l, basePath, btdClient, rmHost, store}
}

func (e *Explorer) Explore(projectID, commit string) error {
	projectPath := filepath.Join(e.store.FullPath(projectID), "bundle")

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		err := e.downloadRepository(projectID, commit)
		if err != nil {
			e.log.Error("Could not download bundled repository", "projectID", projectID, "commit", commit, "err", err)
			return fmt.Errorf("Could not download bundled repository for %s", projectID)
		}
	}

	bp := projectID + ".bundle"
	srcPath := e.store.FullPath(filepath.Join(projectID, "bundle", bp))
	destPath := e.store.FullPath(filepath.Join(projectID, "unbundle"))

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err := e.store.Unbundle(srcPath, destPath)
		if err != nil {
			e.log.Error("Could not unbundle repository", "projectID", projectID, "commit", commit, "err", err)
			return fmt.Errorf("Could not unbundle repository for %s", projectID)
		}
	}

	return nil
}
