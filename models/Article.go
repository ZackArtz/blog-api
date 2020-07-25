package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Article struct {
	ID        uint32    `gorm:"primary_key;auto_increment;" json:"id"`
	Slug      string    `gorm:"size:255;not null;unique" json:"slug"`
	SubHeader string    `gorm:"type:longtext;not null;" json:"sub_header"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Body      string    `gorm:"type:longtext;not null" json:"body"`
	Category  string    `gorm:"size:255;not null" json:"category"`
	Tags      string    `json:"tags"`
	AuthorID  uint32    `gorm:"not null;" json:"author_id,omitempty"`
	CreatedAt time.Time `gorm:"not null;" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;" json:"updated_at"`
}

func (a *Article) Prepare() {
	a.ID = 0
	a.Body = html.EscapeString(strings.TrimSpace(a.Body))
	a.Slug = html.EscapeString(strings.Replace(strings.TrimSpace(strings.ToLower(strings.Replace(a.Title, " ", "-", -1))), ".", "", -1))
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
}

func (a *Article) Validate() error {
	if a.Title == "" {
		return errors.New("required title")
	}
	if a.Body == "" {
		return errors.New("required body")
	}
	if a.AuthorID == 0 {
		return errors.New("required author")
	}
	return nil
}

func (a *Article) CreateArticle(db *gorm.DB) (*Article, error) {
	err := a.Validate()
	if err != nil {
		return &Article{}, err
	}
	a.Prepare()
	err = db.Debug().Model(&Article{}).Create(&a).Error
	if err != nil {
		return &Article{}, err
	}
	return a, nil
}

func (a Article) GetAllArticles(db *gorm.DB) ([]Article, error) {
	var articles []Article
	err := db.Debug().Model(&Article{}).Limit(100).Find(&articles).Error
	if err != nil {
		return []Article{}, err
	}
	return articles, err
}

func (a Article) GetArticleByID(db *gorm.DB, id uint32) (*Article, error) {
	err := db.Debug().Model(&Article{}).Where("id = ?", id).Limit(1).Take(&a).Error
	if err != nil {
		return &Article{}, err
	}
	return &a, nil
}

func (a Article) GetArticleBySlug(db *gorm.DB, slug string) (*Article, error) {
	err := db.Debug().Model(&Article{}).Where("slug = ?", slug).Limit(1).Take(&a).Error
	if err != nil {
		return &Article{}, err
	}
	return &a, nil
}

func (a Article) UpdateArticle(db *gorm.DB, aid uint32) (Article, error) {
	a.Prepare()
	err := db.Debug().Model(&Article{}).Where("id = ?", aid).Updates(Article{
		Title:     a.Title,
		Body:      a.Body,
		Category:  a.Category,
		Tags:      a.Tags,
		AuthorID:  a.AuthorID,
		UpdatedAt: time.Now(),
	}).Error
	if err != nil {
		return Article{}, err
	}
	return a, nil
}

func (a *Article) DeleteArticle(db *gorm.DB, aid uint32) (int64, error) {
	db = db.Debug().Model(&Article{}).Where(" id = ?", aid).Take(&Article{}).Delete(&Article{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("article not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
