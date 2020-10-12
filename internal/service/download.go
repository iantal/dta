package service

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
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

func (e *Explorer) save(projectID, commit string, r io.ReadCloser) error {
	e.log.Info("Save project - storage", "projectID", projectID)

	bp := commit + ".bundle"
	fp := filepath.Join(projectID, "bundle", bp)
	err := e.store.Save(fp, r)

	if err != nil {
		e.log.Error("Unable to save file", "error", err)
		return fmt.Errorf("Unable to save file %s", err)
	}

	return nil
}
