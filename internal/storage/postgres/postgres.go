package postgres

import (
	"context"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/constants"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
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
	_, err := d.Conn.Exec(context.TODO(), "insert into feature (id) values ($1) on conflict do nothing", banner.Feature)
	if err != nil {
		return errors.Wrap(err, "insert feature failed")
	}
	r := d.Conn.QueryRow(context.TODO(), "insert into banners (feature_id, tags, content, is_active) values ($1,$2,$3,$4) returning id", banner.Feature, banner.Tags, banner.Content, banner.IsActive)
	bannerID := 0
	err = r.Scan(&bannerID)
	if err != nil {
		return errors.Wrap(err, "insert banner scan failed")
	}
	for _, tag := range banner.Tags {
		_, err = d.Conn.Exec(context.TODO(), "insert into tags (id) values ($1) on conflict do nothing", tag)
		if err != nil {
			return errors.Wrap(err, "insert tags failed")
		}

		_, err = d.Conn.Exec(context.TODO(), "insert into banners_tags (tag_id, banner_id) values ($1, $2)", tag, bannerID)
		if err != nil {
			return errors.Wrap(err, "insert tags failed")
		}
	}

	return nil
}

func (d *DB) UpdateBanner(id int, banner *models.Banner) error {
	_, err := d.Conn.Exec(context.TODO(), "delete from banners_tags where banner_id = $1", id)
	if err != nil {
		return errors.Wrap(err, "delete banner tags failed")
	}
	_, err = d.Conn.Exec(context.TODO(), "UPDATE banners SET feature_id = $1, tags = $2, content = $3, is_active = $4 WHERE id = $4", banner.Feature, banner.Tags, banner.Content, banner.IsActive, id)
	if err != nil {
		return errors.Wrap(err, "update banner failed")
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

		_, err = d.Conn.Exec(context.TODO(), "insert into banners_tags (tag_id, banner_id) values ($1, $2)", tag, id)
		if err != nil {
			return errors.Wrap(err, "insert tags failed")
		}
	}
	return nil
}

func (d *DB) DeleteBanner(bannerID int) error {
	_, err := d.Conn.Exec(context.TODO(), "delete from banners_tags where banner_id = $1", bannerID)
	if err != nil {
		return errors.Wrap(err, "delete banner tags failed")
	}
	_, err = d.Conn.Exec(context.TODO(), "delete from banners where id = $1", bannerID)
	if err != nil {
		return errors.Wrap(err, "delete banner failed")

	}
	return nil
}

func (d *DB) GetUserBanner(bannerRequest *models.BannerRequest) (*models.Banner, error) {
	var row pgx.Row
	if bannerRequest.TagID != constants.NoValue && bannerRequest.FeatureID != constants.NoValue {
		row = d.Conn.QueryRow(context.TODO(), "select b.id, b.tags, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at from banners b join banners_tags bt on b.id = bt.banner_id where bt.tag_id = $1 and b.feature_id = $2 and b.is_active = true order by b.id desc limit 1", bannerRequest.TagID, bannerRequest.FeatureID)
	} else if bannerRequest.TagID != constants.NoValue && bannerRequest.FeatureID == constants.NoValue {
		row = d.Conn.QueryRow(context.TODO(), "select b.id, b.tags, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at from banners b join banners_tags bt on b.id = bt.banner_id where bt.tag_id = $1 and b.is_active = true order by b.id desc limit 1", bannerRequest.TagID)
	} else {
		row = d.Conn.QueryRow(context.TODO(), "select id, tags, feature_id, content, is_active, created_at, updated_at from banners where feature_id = $1 and is_active = true order by id desc limit 1", bannerRequest.FeatureID)
	}
	banner := models.Banner{}
	err := row.Scan(&banner.ID, &banner.Tags, &banner.Feature, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "select banner failed")
	}
	return &banner, nil
}

// подразумевается, что будет либо tagID, либо featureID, либо оба
func (d *DB) GetAllBanners(bannerRequest *models.BannerRequest) ([]models.Banner, error) {
	var rows pgx.Rows
	var err error
	if bannerRequest.TagID != constants.NoValue && bannerRequest.FeatureID != constants.NoValue {
		rows, err = d.Conn.Query(context.TODO(), "select b.id, b.tags, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at from banners b join banners_tags bt on b.id = bt.banner_id where bt.tag_id = $1 and b.feature_id = $2 order by b.id desc limit $3 offset $4", bannerRequest.TagID, bannerRequest.FeatureID, bannerRequest.Limit, bannerRequest.Offset)
	} else if bannerRequest.TagID != constants.NoValue && bannerRequest.FeatureID == constants.NoValue {
		rows, err = d.Conn.Query(context.TODO(), "select b.id, b.tags, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at from banners b join banners_tags bt on b.id = bt.banner_id where bt.tag_id = $1 order by b.id desc limit $2 offset $3", bannerRequest.TagID, bannerRequest.Limit, bannerRequest.Offset)
	} else {
		rows, err = d.Conn.Query(context.TODO(), "select id, tags, feature_id, content, is_active, created_at, updated_at from banners where feature_id = $1 order by id desc limit $2 offset $3", bannerRequest.FeatureID, bannerRequest.Limit, bannerRequest.Offset)
	}
	if err != nil {
		return nil, errors.Wrap(err, "select banner failed")
	}
	banners := make([]models.Banner, 0)
	for rows.Next() {
		banner := models.Banner{}
		err = rows.Scan(&banner.ID, &banner.Tags, &banner.Feature, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "select banner failed")
		}
		banners = append(banners, banner)
	}
	return banners, nil
}

func (d *DB) Close() error {
	d.Conn.Close()
	return nil
}

func New(cfg *config.Config) (*DB, error) {
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

func initDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
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
