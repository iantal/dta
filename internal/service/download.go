package service

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func (e *Explorer) downloadRepository(projectID, commit string) error {
	resp, err := http.DefaultClient.Get("http://" + e.rmHost + "/api/v1/projects/" + projectID + "/" + commit + "/download")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected error code 200 got %d", resp.StatusCode)
	}

	e.log.Info("Content-Dispozition", "file", resp.Header.Get("Content-Disposition"))

	e.save(projectID, commit, resp.Body)
	resp.Body.Close()

	return nil
}

func (e *Explorer) save(projectID, commit string, r io.ReadCloser) {
	e.log.WithFields(logrus.Fields{
		"projectID": projectID,
		"commit":    commit,
	}).Info("Save project to storage")

	bp := commit + ".bundle"
	fp := filepath.Join(projectID, "bundle", bp)
	err := e.store.Save(fp, r)

	if err != nil {
		e.log.WithFields(logrus.Fields{
			"projectID": projectID,
			"commit":    commit,
			"error":     err,
		}).Error("Unable to save file")
	}
}
