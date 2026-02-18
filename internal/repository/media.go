package repository

import (
	"context"
	"database/sql"
	"golang-yt/domain"

	"github.com/doug-martin/goqu/v9"
)

type mediaRepository struct {
	db *goqu.Database
}

func Newmedia(con *sql.DB) domain.MediaRepository {
	return &mediaRepository{
		db: goqu.New("postgres", con),
	}
}

// FindById implements [domain.MediaRepository].
func (m *mediaRepository) FindById(ctx context.Context, id string) (media domain.Media,err error) {
	dataset := m.db.From("media").Where(goqu.Ex{
		"id" : id,
	})

	_, err = dataset.ScanStructContext(ctx, &media)
	return
}

// FindByIds implements [domain.MediaRepository].
func (m *mediaRepository) FindByIds(ctx context.Context, ids []string) (media []domain.Media, err error) {
	dataset := m.db.From("media").Where(goqu.C("id").In(ids))

	err = dataset.ScanStructsContext(ctx, &media)
	return
}

// Save implements [domain.MediaRepository].
func (m *mediaRepository) Save(ctx context.Context, media *domain.Media) error {
	_, err := m.db.Insert("media").Rows(media).Executor().ExecContext(ctx)
	return err
}