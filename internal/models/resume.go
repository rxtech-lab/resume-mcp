package models

import (
	"time"

	"gorm.io/gorm"
)

type Resume struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Photo       string    `json:"photo"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Contacts        []Contact        `gorm:"foreignKey:ResumeID" json:"contacts,omitempty"`
	WorkExperiences []WorkExperience `gorm:"foreignKey:ResumeID" json:"work_experiences,omitempty"`
	Educations      []Education      `gorm:"foreignKey:ResumeID" json:"educations,omitempty"`
	OtherExperiences []OtherExperience `gorm:"foreignKey:ResumeID" json:"other_experiences,omitempty"`
	Templates       []Template       `gorm:"foreignKey:ResumeID;constraint:OnDelete:CASCADE" json:"templates,omitempty"`
}

type Contact struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ResumeID uint   `gorm:"not null" json:"resume_id"`
	Key      string `gorm:"not null" json:"key"`
	Value    string `gorm:"not null" json:"value"`
	Resume   Resume `gorm:"foreignKey:ResumeID" json:"-"`
}

type WorkExperience struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ResumeID  uint      `gorm:"not null" json:"resume_id"`
	Company   string    `gorm:"not null" json:"company"`
	JobTitle  string    `gorm:"not null" json:"job_title"`
	Type      string    `gorm:"default:fulltime" json:"type"` // fulltime, parttime, internship
	StartDate time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Resume    Resume    `gorm:"foreignKey:ResumeID" json:"-"`

	FeatureMaps []FeatureMap `gorm:"foreignKey:ExperienceID;constraint:OnDelete:CASCADE" json:"feature_maps,omitempty"`
}

type Education struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	ResumeID   uint       `gorm:"not null" json:"resume_id"`
	SchoolName string     `gorm:"not null" json:"school_name"`
	Type       string     `gorm:"default:fulltime" json:"type"` // fulltime, parttime, internship
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Resume     Resume     `gorm:"foreignKey:ResumeID" json:"-"`

	FeatureMaps []FeatureMap `gorm:"foreignKey:ExperienceID;constraint:OnDelete:CASCADE" json:"feature_maps,omitempty"`
}

type OtherExperience struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ResumeID uint   `gorm:"not null" json:"resume_id"`
	Category string `gorm:"not null" json:"category"`
	Resume   Resume `gorm:"foreignKey:ResumeID" json:"-"`

	FeatureMaps []FeatureMap `gorm:"foreignKey:ExperienceID;constraint:OnDelete:CASCADE" json:"feature_maps,omitempty"`
}

type FeatureMap struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	ExperienceID uint   `gorm:"not null" json:"experience_id"`
	Key          string `gorm:"not null" json:"key"`
	Value        string `gorm:"type:text" json:"value"`
}

type PreviewSession struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	ResumeID  uint      `gorm:"not null" json:"resume_id"`
	Template  string    `gorm:"type:text;not null" json:"template"`
	CSS       string    `gorm:"type:text" json:"css"`
	CreatedAt time.Time `json:"created_at"`
	Resume    Resume    `gorm:"foreignKey:ResumeID" json:"resume"`
}

type Template struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ResumeID     uint      `gorm:"not null" json:"resume_id"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	TemplateData string    `gorm:"type:text;not null" json:"template_data"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Resume       Resume    `gorm:"foreignKey:ResumeID" json:"-"`
}