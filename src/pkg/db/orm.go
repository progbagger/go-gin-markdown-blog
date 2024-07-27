package db

import (
	"html/template"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model

	Title string
	Text  template.HTML

	Timestamp time.Time
}
