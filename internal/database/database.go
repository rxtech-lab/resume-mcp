package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rxtech-lab/resume-mcp/internal/models"
	"gorm.io/driver/postgres"
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

// NewPostgresDatabase creates a new database connection to a PostgreSQL database
func NewPostgresDatabase(postgresURL string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{
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

func (d *Database) CreateResume(resume *models.Resume, userID *string) error {
	if userID != nil {
		resume.UserID = *userID
	}
	return d.DB.Create(resume).Error
}

func (d *Database) GetResumeByName(name string, userID *string) (*models.Resume, error) {
	var resume models.Resume
	query := d.DB.Preload("Contacts").
		Preload("WorkExperiences.FeatureMaps").
		Preload("Educations.FeatureMaps").
		Preload("OtherExperiences.FeatureMaps").
		Where("name = ?", name)
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.First(&resume).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (d *Database) GetResumeByID(id uint, userID *string) (*models.Resume, error) {
	var resume models.Resume
	query := d.DB.Preload("Contacts").
		Preload("WorkExperiences.FeatureMaps").
		Preload("Educations.FeatureMaps").
		Preload("OtherExperiences.FeatureMaps")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.First(&resume, id).Error
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (d *Database) ListResumes(userID *string) ([]models.Resume, error) {
	var resumes []models.Resume
	query := d.DB.Select("id, name, created_at, updated_at")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.Find(&resumes).Error
	return resumes, err
}

func (d *Database) UpdateResume(resume *models.Resume, userID *string) error {
	if userID != nil {
		resume.UserID = *userID
	}
	return d.DB.Save(resume).Error
}

func (d *Database) DeleteResume(id uint, userID *string) error {
	query := d.DB
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	return query.Delete(&models.Resume{}, id).Error
}

func (d *Database) AddContact(contact *models.Contact, userID *string) error {
	if userID != nil {
		contact.UserID = *userID
	}
	return d.DB.Create(contact).Error
}

func (d *Database) UpdateContact(contact *models.Contact, userID *string) error {
	if userID != nil {
		contact.UserID = *userID
	}
	return d.DB.Save(contact).Error
}

func (d *Database) AddWorkExperience(experience *models.WorkExperience, userID *string) error {
	if userID != nil {
		experience.UserID = *userID
	}
	return d.DB.Create(experience).Error
}

func (d *Database) AddEducation(education *models.Education, userID *string) error {
	if userID != nil {
		education.UserID = *userID
	}
	return d.DB.Create(education).Error
}

func (d *Database) AddOtherExperience(experience *models.OtherExperience, userID *string) error {
	if userID != nil {
		experience.UserID = *userID
	}
	return d.DB.Create(experience).Error
}

func (d *Database) AddFeatureMap(featureMap *models.FeatureMap, userID *string) error {
	if userID != nil {
		featureMap.UserID = *userID
	}
	return d.DB.Create(featureMap).Error
}

func (d *Database) UpdateFeatureMap(featureMap *models.FeatureMap, userID *string) error {
	if userID != nil {
		featureMap.UserID = *userID
	}
	return d.DB.Save(featureMap).Error
}

func (d *Database) GetFeatureMapByID(id uint, userID *string) (*models.FeatureMap, error) {
	var featureMap models.FeatureMap
	query := d.DB
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.First(&featureMap, id).Error
	if err != nil {
		return nil, err
	}
	return &featureMap, nil
}

func (d *Database) DeleteFeatureMap(id uint, userID *string) error {
	query := d.DB
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	return query.Delete(&models.FeatureMap{}, id).Error
}

func (d *Database) GeneratePreview(resumeID uint, template string, css string, userID *string) (string, error) {
	sessionID := uuid.New().String()
	session := &models.PreviewSession{
		ID:       sessionID,
		ResumeID: resumeID,
		Template: template,
		CSS:      css,
	}
	if userID != nil {
		session.UserID = *userID
	}

	err := d.DB.Create(session).Error
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (d *Database) CreatePreviewSession(session *models.PreviewSession, userID *string) error {
	if userID != nil {
		session.UserID = *userID
	}
	return d.DB.Create(session).Error
}

func (d *Database) GetPreviewSession(sessionID string, userID *string) (*models.PreviewSession, error) {
	var session models.PreviewSession
	query := d.DB.Preload("Resume.Contacts").
		Preload("Resume.WorkExperiences.FeatureMaps").
		Preload("Resume.Educations.FeatureMaps").
		Preload("Resume.OtherExperiences.FeatureMaps").
		Where("id = ?", sessionID)
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (d *Database) UpdatePreviewSessionCSS(sessionID string, css string, userID *string) error {
	query := d.DB.Model(&models.PreviewSession{}).
		Where("id = ?", sessionID)
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	return query.Update("css", css).Error
}

// Template CRUD operations
func (d *Database) CreateTemplate(template *models.Template, userID *string) error {
	if userID != nil {
		template.UserID = *userID
	}
	return d.DB.Create(template).Error
}

func (d *Database) GetTemplateByID(id uint, userID *string) (*models.Template, error) {
	var template models.Template
	query := d.DB
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (d *Database) ListTemplatesByResumeID(resumeID uint, userID *string) ([]models.Template, error) {
	var templates []models.Template
	query := d.DB.Where("resume_id = ?", resumeID)
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.Find(&templates).Error
	return templates, err
}

func (d *Database) UpdateTemplate(template *models.Template, userID *string) error {
	if userID != nil {
		template.UserID = *userID
	}
	return d.DB.Save(template).Error
}

func (d *Database) DeleteTemplate(id uint, userID *string) error {
	query := d.DB
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	return query.Delete(&models.Template{}, id).Error
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
