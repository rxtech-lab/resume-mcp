package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rxtech-lab/resume-mcp/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: nil, // Disable GORM logging to prevent color output
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	database := &Database{DB: db}
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}

func (d *Database) migrate() error {
	return d.DB.AutoMigrate(
		&models.Resume{},
		&models.Contact{},
		&models.WorkExperience{},
		&models.Education{},
		&models.OtherExperience{},
		&models.FeatureMap{},
		&models.PreviewSession{},
		&models.Template{},
	)
}

func (d *Database) CreateResume(resume *models.Resume) error {
	return d.DB.Create(resume).Error
}

func (d *Database) GetResumeByName(name string) (*models.Resume, error) {
	var resume models.Resume
	err := d.DB.Preload("Contacts").
		Preload("WorkExperiences.FeatureMaps").
		Preload("Educations.FeatureMaps").
		Preload("OtherExperiences.FeatureMaps").
		Where("name = ?", name).
		First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (d *Database) GetResumeByID(id uint) (*models.Resume, error) {
	var resume models.Resume
	err := d.DB.Preload("Contacts").
		Preload("WorkExperiences.FeatureMaps").
		Preload("Educations.FeatureMaps").
		Preload("OtherExperiences.FeatureMaps").
		First(&resume, id).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (d *Database) ListResumes() ([]models.Resume, error) {
	var resumes []models.Resume
	err := d.DB.Select("id, name, created_at, updated_at").Find(&resumes).Error
	return resumes, err
}

func (d *Database) UpdateResume(resume *models.Resume) error {
	return d.DB.Save(resume).Error
}

func (d *Database) DeleteResume(id uint) error {
	return d.DB.Delete(&models.Resume{}, id).Error
}

func (d *Database) AddContact(contact *models.Contact) error {
	return d.DB.Create(contact).Error
}

func (d *Database) UpdateContact(contact *models.Contact) error {
	return d.DB.Save(contact).Error
}

func (d *Database) AddWorkExperience(experience *models.WorkExperience) error {
	return d.DB.Create(experience).Error
}

func (d *Database) AddEducation(education *models.Education) error {
	return d.DB.Create(education).Error
}

func (d *Database) AddOtherExperience(experience *models.OtherExperience) error {
	return d.DB.Create(experience).Error
}

func (d *Database) AddFeatureMap(featureMap *models.FeatureMap) error {
	return d.DB.Create(featureMap).Error
}

func (d *Database) UpdateFeatureMap(featureMap *models.FeatureMap) error {
	return d.DB.Save(featureMap).Error
}

func (d *Database) GetFeatureMapByID(id uint) (*models.FeatureMap, error) {
	var featureMap models.FeatureMap
	err := d.DB.First(&featureMap, id).Error
	if err != nil {
		return nil, err
	}
	return &featureMap, nil
}

func (d *Database) DeleteFeatureMap(id uint) error {
	return d.DB.Delete(&models.FeatureMap{}, id).Error
}

func (d *Database) GeneratePreview(resumeID uint, template string, css string) (string, error) {
	sessionID := uuid.New().String()
	session := &models.PreviewSession{
		ID:       sessionID,
		ResumeID: resumeID,
		Template: template,
		CSS:      css,
	}
	
	err := d.DB.Create(session).Error
	if err != nil {
		return "", err
	}
	
	return sessionID, nil
}

func (d *Database) CreatePreviewSession(session *models.PreviewSession) error {
	return d.DB.Create(session).Error
}

func (d *Database) GetPreviewSession(sessionID string) (*models.PreviewSession, error) {
	var session models.PreviewSession
	err := d.DB.Preload("Resume.Contacts").
		Preload("Resume.WorkExperiences.FeatureMaps").
		Preload("Resume.Educations.FeatureMaps").
		Preload("Resume.OtherExperiences.FeatureMaps").
		Where("id = ?", sessionID).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (d *Database) UpdatePreviewSessionCSS(sessionID string, css string) error {
	return d.DB.Model(&models.PreviewSession{}).
		Where("id = ?", sessionID).
		Update("css", css).Error
}

// Template CRUD operations
func (d *Database) CreateTemplate(template *models.Template) error {
	return d.DB.Create(template).Error
}

func (d *Database) GetTemplateByID(id uint) (*models.Template, error) {
	var template models.Template
	err := d.DB.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (d *Database) ListTemplatesByResumeID(resumeID uint) ([]models.Template, error) {
	var templates []models.Template
	err := d.DB.Where("resume_id = ?", resumeID).Find(&templates).Error
	return templates, err
}

func (d *Database) UpdateTemplate(template *models.Template) error {
	return d.DB.Save(template).Error
}

func (d *Database) DeleteTemplate(id uint) error {
	return d.DB.Delete(&models.Template{}, id).Error
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}