package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/iantal/dta/internal/domain"
	"github.com/jinzhu/gorm"
)

// ProjectDB defines the CRUD operations for storing projects in the db
type ProjectDB struct {
	log hclog.Logger
	db  *gorm.DB
}

// NewProjectDB returns a ProjectDB object for handling CRUD operations
func NewProjectDB(log hclog.Logger, db *gorm.DB) *ProjectDB {
	db.AutoMigrate(&domain.Project{})
	return &ProjectDB{
		log: log,
		db:  db,
	}
}

// AddProject adds a project to the db
func (p *ProjectDB) AddProject(project *domain.Project) {
	p.db.Create(&project)
	return
}

// UpdateProject updates the analysis status of the project 
func (p *ProjectDB) UpdateProject(project *domain.Project, status domain.Status) {
	p.db.Model(&project).Update("status", status.String())
}

// UpdateProjectBuildToolAndStatus updates the build tool and status of project
func (p *ProjectDB) UpdateProjectBuildToolAndStatus(project *domain.Project, status domain.Status, buildTool string) {
	p.db.Model(&project).Updates(domain.Project{Status: status.String(), BuildTool: buildTool})
}

// UpdateProjectName updates the name of a project
func (p *ProjectDB) UpdateProjectName(project *domain.Project, name string) {
	p.db.Model(&project).Update("name", name)
}

// GetProjects returns all existing projects in the db
func (p *ProjectDB) GetProjects() ([]*domain.Project, error) {
	var projects []*domain.Project
	p.db.Find(&projects)
	return projects, nil
}

// GetProjectByID returns the project with the given id
func (p *ProjectDB) GetProjectByID(id string) *domain.Project {
	project := &domain.Project{}
	uid, err := uuid.Parse(id)
	if err != nil {
		p.log.Error("Project with projectId {} was not found")
		return nil
	}
	p.db.Find(&project, "project_id = ?", uid)
	return project
}

// GetProjectByIDAndCommit returns the project with the given id and commit
func (p *ProjectDB) GetProjectByIDAndCommit(id, commit string) *domain.Project {
	project := &domain.Project{}
	uid, err := uuid.Parse(id)
	if err != nil {
		p.log.Error("Project with projectId {} was not found")
		return nil
	}
	p.db.Where("project_id = ? AND commit_hash = ?", uid, commit).Find(&project)
	return project
}
