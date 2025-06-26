// models/document.go
// This file contains the models for the document resource
// It defines the Document struct and the DocumentModel struct
// It also contains the methods for the DocumentModel struct
// It uses the GORM library to interact with the database

package models

import (
	"time"
	"github.com/atheeralattar/pbl-week2/internal/kafka"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Document struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required"`
	Content   []byte    `json:"content" binding:"required"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DocumentModel handles all database operations for documents
type DocumentModel struct {
	DB *gorm.DB
}

// NewDocumentModel creates a new document model
func NewDocumentModel(db *gorm.DB) *DocumentModel {
	return &DocumentModel{DB: db}
}

// Create adds a new document to the database
func (m *DocumentModel) Create(doc *Document) error {
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	err := m.DB.Create(doc).Error
	if err != nil {
		return err // don't publish Kafka event if DB write fails
	}
	err = kafka.PublishDocEvent("created", fmt.Sprint(doc.ID), doc.Title, string(doc.Content))
	if err != nil {
		log.Println("[KAFKA] Failed to publish event:", err)
	}
	return nil

}

// FindAll retrieves all documents
func (m *DocumentModel) FindAll() ([]Document, error) {
	var documents []Document
	err := m.DB.Find(&documents).Error
	return documents, err
}

// FindByID retrieves a document by its ID
func (m *DocumentModel) FindByID(id string) (Document, error) {
	var document Document
	err := m.DB.First(&document, id).Error
	return document, err
}

// Update modifies an existing document
func (m *DocumentModel) Update(doc *Document) error {
	doc.UpdatedAt = time.Now()
	err := m.DB.Save(doc).Error
	if err != nil {
		return err
	}
	err = kafka.PublishDocEvent("updated", fmt.Sprint(doc.ID), doc.Title, string(doc.Content))
	if err != nil {
		log.Println("[KAFKA] Failed to publish event:", err)
	}
	return nil

}

// Delete removes a document
func (m *DocumentModel) Delete(id string) error {
	var doc Document
	if err := m.DB.First(&doc, id).Error; err != nil {
		return err
	}
	err := m.DB.Delete(&Document{}, id).Error
	if err != nil {
		return err
	}
	err = kafka.PublishDocEvent("deleted", fmt.Sprint(doc.ID), doc.Title, string(doc.Content))
	if err != nil {
		log.Println("[KAFKA] Failed to publish event:", err)
	}
	return nil

}
