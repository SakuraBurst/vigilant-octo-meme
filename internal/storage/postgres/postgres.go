package postgres

import (
	"context"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/constants"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
	"path/filepath"
	"strings"
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

		_, err = d.Conn.Exec(context.TODO(), "insert into banner_tags (tag_id, banner_id) values ($1, $2)", tag, bannerID)
		if err != nil {
			return errors.Wrap(err, "insert tags failed")
		}
	}

	return nil
}

func (d *DB) UpdateBanner(id int, banner *models.Banner) error {
	_, err := d.Conn.Exec(context.TODO(), "delete from banner_tags where banner_id = $1", id)
	if err != nil {
		return errors.Wrap(err, "delete banner tags failed")
	}
	_, err = d.Conn.Exec(context.TODO(), "UPDATE banner SET feature_id = $1, content = $2, is_active = $3 WHERE id = $4", banner.Feature, banner.Content, banner.IsActive, id)
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

		_, err = d.Conn.Exec(context.TODO(), "insert into banner_tags (tag_id, banner_id) values ($1, $2)", tag, id)
		if err != nil {
			return errors.Wrap(err, "insert tags failed")
		}
	}
	return nil
}

func (d *DB) DeleteBanner(bannerId int) error {
	_, err := d.Conn.Exec(context.TODO(), "delete from banner_tags where banner_id = $1", bannerId)
	if err != nil {
		return errors.Wrap(err, "delete banner tags failed")
	}
	_, err = d.Conn.Exec(context.TODO(), "delete from banner where id = $1", bannerId)
	if err != nil {
		return errors.Wrap(err, "delete banner failed")

	}
	return nil
}

func (d *DB) GetUserBanner(bannerRequest *models.BannerRequest) (*models.Banner, error) {
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("select * from banner b ")
	if bannerRequest.TagId != constants.NoValue {
		sqlBuilder.WriteString("join banner_tags bt on b.id = bt.banner_id where bt.tag_id = $1 ")
	}
	if bannerRequest.FeatureId != constants.NoValue && bannerRequest.TagId != constants.NoValue {
		sqlBuilder.WriteString("and b.feature_id = $2 ")
	}
	if bannerRequest.FeatureId != constants.NoValue && bannerRequest.TagId == constants.NoValue {
		sqlBuilder.WriteString("where b.feature_id = $1 ")
	}
	sqlBuilder.WriteString("and b.is_active = true order by b.id desc limit 1")
	row := d.Conn.QueryRow(context.TODO(), sqlBuilder.String(), bannerRequest.FeatureId, bannerRequest.TagId)
	banner := models.Banner{}
	err := row.Scan(&banner.ID, &banner.Feature, &banner.Content, &banner.IsActive)
	if err != nil {
		return nil, errors.Wrap(err, "select banner failed")
	}
	return &banner, nil
}

func (d *DB) GetAllBanners(bannerRequest *models.BannerRequest) ([]models.Banner, error) {
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("select * from banner b ")
	if bannerRequest.TagId != constants.NoValue {
		sqlBuilder.WriteString("join banner_tags bt on b.id = bt.banner_id where bt.tag_id = $1 ")
	}
	if bannerRequest.FeatureId != constants.NoValue && bannerRequest.TagId != constants.NoValue {
		sqlBuilder.WriteString("and b.feature_id = $2 ")
	}
	if bannerRequest.FeatureId != constants.NoValue && bannerRequest.TagId == constants.NoValue {
		sqlBuilder.WriteString("where b.feature_id = $1 ")
	}
	sqlBuilder.WriteString("order by b.id desc limit $3 offset $4")
	rows, err := d.Conn.Query(context.TODO(), sqlBuilder.String(), bannerRequest.FeatureId, bannerRequest.TagId, bannerRequest.Limit, bannerRequest.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "select banner failed")
	}
	banners := make([]models.Banner, 0)
	for rows.Next() {
		banner := models.Banner{}
		err = rows.Scan(&banner.ID, &banner.Feature, &banner.Content, &banner.IsActive)
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
