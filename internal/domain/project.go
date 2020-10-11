package domain

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Project defines data related to a project repository
type Project struct {
	gorm.Model `json:"-"`
	ProjectID  uuid.UUID `gorm:"type:uuid;primary_key;" json:"projectId"`
	CommitHash string    `gorm:"primary_key" json:"commit,omitempty"`
	Name       string    `json:"name,omitempty"`
	Status     string    `json:"status,omitempty"`
	BuildTool  string    `json:"-"`
	Data       string    `json:"-"`
}

// NewProject creates an instance of Project
func NewProject(id uuid.UUID, commit, name, status, buildTool, data string) *Project {
	return &Project{
		ProjectID:  id,
		CommitHash: commit,
		Name:       name,
		Status:     status,
		BuildTool:  buildTool,
		Data:       data,
	}
}
