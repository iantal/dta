package domain

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Library represents a third party library from the dependency tree
type Library struct {
	gorm.Model `json:"-"`
	ProjectID  uuid.UUID `gorm:"type:uuid;primary_key;" json:"projectId"`
	CommitHash string    `gorm:"primary_key" json:"commit,omitempty"`
	Name       string    `json:"name,omitempty"`
	Type       string    `json:"type,omitempty"`
	Scope      string    `json:"scope,omitempty"`
	Packaging  string    `json:"packaging,omitempty"`
}

// NewLibrary creates a Library
func NewLibrary(projectID uuid.UUID, commitHash, name, libraryType, scope, packaging string) *Library {
	return &Library{
		ProjectID:  projectID,
		CommitHash: commitHash,
		Name:       name,
		Type:       libraryType,
		Scope:      scope,
		Packaging: packaging,
	}
}
