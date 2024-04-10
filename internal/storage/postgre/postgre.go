package postgre

import (
	"context"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
	"path/filepath"
	"time"
)

type DB struct {
	Conn *pgxpool.Pool
}

func (d *DB) SaveBanner(banner *models.Banner) error {
	r := d.Conn.QueryRow(context.TODO(), "insert into banner (feature_id, content, is_active) values ($1,$2,$3)", banner.Feature, banner.Content, banner.IsActive)
	bannerID := 0
	err := r.Scan(&bannerID)
	if err != nil {
		return errors.Wrap(err, "insert banner scan failed")
	}
	_, err = d.Conn.Exec(context.TODO(), "insert into feature (id) values ($1) on conflict do nothing", banner.Feature)
	if err != nil {
		return errors.Wrap(err, "insert feature failed")
	}
	for _, tag := range banner.Tags {
		_, err = d.Conn.Exec(context.TODO(), "insert into tags (id) values ($1) on conflict do nothing", tag)
		if err != nil {
			return errors.Wrap(err, "insert tags failed")
		}

		_, err = d.Conn.Exec(context.TODO(), "insert into banner_tags (tag_id, banner_id) values ($1, $2) on conflict do nothing", tag, bannerID)
		if err != nil {
			return errors.Wrap(err, "insert tags failed")
		}
	}

	return nil
}

func (d *DB) DeleteBanner(banner *models.Banner) error {

	return nil
}

func (d *DB) GetBanner() error {
	return nil

}

func (d *DB) GetAllBanners() error {

	return nil
}

func (d *DB) Close() error {
	d.Conn.Close()
	return nil
}

func NewDB(cfg config.Config) (*DB, error) {
	conn, err := initDatabase(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "InitDatabase failed")
	}
	return &DB{Conn: conn}, nil
}

func onConnectScript(conn *pgxpool.Pool) error {
	ctx, cl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cl()
	file, err := os.Open(filepath.Join(".", "config", "init.sql"))
	if err != nil {
		return errors.Wrap(err, "os.Open failed")
	}
	sqlScript, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "io.ReadAll failed")
	}
	_, err = conn.Exec(ctx, string(sqlScript))
	if err != nil {
		return err
	}
	return nil
}

func initDatabase(cfg config.Config) (*pgxpool.Pool, error) {
	ctx, cl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cl()
	pool, err := pgxpool.New(ctx, cfg.StoragePath)
	if err != nil {

		return nil, errors.Wrap(err, "pgxpool.New failed")
	}
	err = onConnectScript(pool)
	if err != nil {
		return nil, errors.Wrap(err, "onConnectScript failed")
	}
	return pool, nil
}
