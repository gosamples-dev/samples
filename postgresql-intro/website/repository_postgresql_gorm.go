package website

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type gormWebsite struct {
	ID   int64  `gorm:"primary_key"`
	Name string `gorm:"uniqueIndex;not null"`
	URL  string `gorm:"not null"`
	Rank int64  `gorm:"not null"`
}

func (gormWebsite) TableName() string {
	return "websites"
}

type PostgreSQLGORMRepository struct {
	db *gorm.DB
}

func NewPostgreSQLGORMRepository(db *gorm.DB) *PostgreSQLGORMRepository {
	return &PostgreSQLGORMRepository{
		db: db,
	}
}

func (r *PostgreSQLGORMRepository) Migrate(ctx context.Context) error {
	m := &gormWebsite{}
	return r.db.WithContext(ctx).AutoMigrate(&m)
}

func (r *PostgreSQLGORMRepository) Create(ctx context.Context, website Website) (*Website, error) {
	gormWebsite := gormWebsite{
		Name: website.Name,
		URL:  website.URL,
		Rank: website.Rank,
	}

	if err := r.db.WithContext(ctx).Create(&gormWebsite).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	result := Website(gormWebsite)

	return &result, nil
}

func (r *PostgreSQLGORMRepository) All(ctx context.Context) ([]Website, error) {
	var gormWebsites []gormWebsite
	if err := r.db.WithContext(ctx).Find(&gormWebsites).Error; err != nil {
		return nil, err
	}

	var result []Website
	for _, gw := range gormWebsites {
		result = append(result, Website(gw))
	}
	return result, nil
}

func (r *PostgreSQLGORMRepository) GetByName(ctx context.Context, name string) (*Website, error) {
	var gormWebsite gormWebsite
	if err := r.db.WithContext(ctx).Where("name = ?", name).Find(&gormWebsite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	website := Website(gormWebsite)
	return &website, nil
}

func (r *PostgreSQLGORMRepository) Update(ctx context.Context, id int64, updated Website) (*Website, error) {
	gormWebsite := Website(updated)
	updateRes := r.db.WithContext(ctx).Where("id = ?", id).Save(&gormWebsite)
	if err := updateRes.Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected := updateRes.RowsAffected
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *PostgreSQLGORMRepository) Delete(ctx context.Context, id int64) error {
	deleteRes := r.db.WithContext(ctx).Delete(&gormWebsite{}, id)
	if err := deleteRes.Error; err != nil {
		return err
	}

	if deleteRes.RowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
