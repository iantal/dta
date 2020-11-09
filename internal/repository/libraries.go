package repository

import (
	"github.com/google/uuid"
	"github.com/iantal/dta/internal/domain"
	"github.com/iantal/dta/internal/utils"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// LibraryDB defines the CRUD operations for storing projects in the db
type LibraryDB struct {
	log *utils.StandardLogger
	db  *gorm.DB
}

// NewLibraryDB returns a LibraryDB object for handling CRUD operations
func NewLibraryDB(log *utils.StandardLogger, db *gorm.DB) *LibraryDB {
	db.AutoMigrate(&domain.Library{})
	return &LibraryDB{
		log: log,
		db:  db,
	}
}

// AddLibrary adds a library to the db
func (l *LibraryDB) AddLibrary(library *domain.Library) {
	l.db.Create(&library)
	return
}

// GetLibraryByIDAndCommit returns all libraries for the given id and commit
func (l *LibraryDB) GetLibraryByIDAndCommit(id, commit string) []*domain.Library {
	var libraries []*domain.Library
	uid, err := uuid.Parse(id)
	if err != nil {
		l.log.WithFields(logrus.Fields{
			"projectID": id,
			"commit":    commit,
		}).Error("No libraries were found")
		return nil
	}
	l.db.Where("project_id = ? AND commit_hash = ?", uid, commit).Find(&libraries)
	return libraries
}
