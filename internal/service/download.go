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

	e.save(projectID, resp.Body)
	resp.Body.Close()

	return nil
}

func (a *Explorer) save(projectID string, r io.ReadCloser) error {
	a.log.Info("Save project - storage", "projectID", projectID)

	bp := projectID + ".bundle"
	fp := filepath.Join(projectID, "bundle", bp)
	err := a.store.Save(fp, r)

	if err != nil {
		a.log.Error("Unable to save file", "error", err)
		return fmt.Errorf("Unable to save file %s", err)
	}

	return nil
}
