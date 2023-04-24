package model

import "github.com/google/uuid"

type Course struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Students    []string `json:"students"`
}

type CourseInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Students    []string `json:"students"`
}

type CourseUpdateInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Students    []string `json:"students"`
}

func NewCourse(courseInput *CourseInput) (*Course, error) {
	course := &Course{
		ID:          generateID(),
		Name:        courseInput.Name,
		Description: courseInput.Description,
		Students:    courseInput.Students,
	}

	return course, nil
}

// Update обновляет информацию о курсе на основе входных параметров
func (c *Course) Update(courseInput *CourseUpdateInput) error {
	if courseInput.Name != "" {
		c.Name = courseInput.Name
	}

	if courseInput.Description != "" {
		c.Description = courseInput.Description
	}

	if courseInput.Students != nil {
		c.Students = courseInput.Students
	}

	return nil
}

// generateID генерирует уникальный идентификатор
func generateID() string {
	u := uuid.New()
	return u.String()
}
