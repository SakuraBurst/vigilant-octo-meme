package postgres

import (
	"context"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/constants"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/storage"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type DB struct {
	Conn *pgxpool.Pool
	log  *slog.Logger
}

func (d *DB) SaveBanner(ctx context.Context, banner *models.Banner) (int, error) {
	log := d.log.With(slog.String("method", "SaveBanner"))
	log.Info("SaveBanner", slog.Any("banner", banner))
	_, err := d.Conn.Exec(ctx, "insert into feature (id) values ($1) on conflict do nothing", banner.Feature)
	if err != nil {
		log.Error(err.Error())
		return constants.NoValue, errors.Wrap(err, "insert feature failed")
	}
	r := d.Conn.QueryRow(ctx, "insert into banners (feature_id, content, is_active) values ($1,$2,$3) returning id", banner.Feature, banner.Content, banner.IsActive)
	bannerID := 0
	err = r.Scan(&bannerID)
	if err != nil {
		log.Error(err.Error())
		return constants.NoValue, errors.Wrap(err, "insert banner scan failed")
	}
	for _, tag := range banner.Tags {
		_, err = d.Conn.Exec(ctx, "insert into tags (id) values ($1) on conflict do nothing", tag)
		if err != nil {
			log.Error(err.Error())
			return constants.NoValue, errors.Wrap(err, "insert tags failed")
		}

		_, err = d.Conn.Exec(ctx, "insert into banners_tags (tag_id, banner_id) values ($1, $2)", tag, bannerID)
		if err != nil {
			log.Error(err.Error())
			return constants.NoValue, errors.Wrap(err, "insert tags failed")
		}
	}

	return bannerID, nil
}

func (d *DB) UpdateBanner(ctx context.Context, id int, banner *models.Banner) error {
	log := d.log.With(slog.String("method", "UpdateBanner"))
	log.Info("UpdateBanner", slog.Int("id", id), slog.Any("banner", banner))
	_, err := d.Conn.Exec(ctx, "delete from banners_tags where banner_id = $1", id)
	if err != nil {
		log.Error(err.Error())
		return errors.Wrap(err, "delete banner tags failed")
	}
	_, err = d.Conn.Exec(ctx, "UPDATE banners SET feature_id = $1, content = $2, is_active = $3 WHERE id = $4", banner.Feature, banner.Content, banner.IsActive, id)
	if err != nil {
		log.Error(err.Error())
		return errors.Wrap(err, "update banner failed")
	}
	_, err = d.Conn.Exec(ctx, "insert into feature (id) values ($1) on conflict do nothing", banner.Feature)
	if err != nil {
		slog.Any("banner", banner)
		return errors.Wrap(err, "insert feature failed")
	}
	for _, tag := range banner.Tags {
		_, err = d.Conn.Exec(ctx, "insert into tags (id) values ($1) on conflict do nothing", tag)
		if err != nil {
			log.Error(err.Error())
			return errors.Wrap(err, "insert tags failed")
		}

		_, err = d.Conn.Exec(ctx, "insert into banners_tags (tag_id, banner_id) values ($1, $2)", tag, id)
		if err != nil {
			log.Error(err.Error())
			return errors.Wrap(err, "insert tags failed")
		}
	}
	return nil
}

func (d *DB) DeleteBanner(ctx context.Context, bannerID int) error {
	log := d.log.With(slog.String("method", "DeleteBanner"))
	log.Info("DeleteBanner", slog.Int("bannerID", bannerID))
	row := d.Conn.QueryRow(ctx, "select id from banners where id = $1", bannerID)
	err := row.Scan(&bannerID)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.BannerNotFound
		}
		return errors.Wrap(err, "select banner failed")
	}
	_, err = d.Conn.Exec(ctx, "delete from banners_tags where banner_id = $1", bannerID)
	if err != nil {
		log.Error(err.Error())
		return errors.Wrap(err, "delete banner tags failed")
	}
	_, err = d.Conn.Exec(ctx, "delete from banners where id = $1", bannerID)
	if err != nil {
		log.Error(err.Error())
		return errors.Wrap(err, "delete banner failed")
	}
	return nil
}

// GetUserBanner принимает запрос на банер, подразумевается, что будет либо tagID, либо featureID, либо оба
func (d *DB) GetUserBanner(ctx context.Context, bannerRequest *models.BannerRequest) (map[string]interface{}, error) {
	log := d.log.With(slog.String("method", "GetUserBanner"))
	log.Info("GetUserBanner", slog.Any("bannerRequest", bannerRequest))
	var row pgx.Row
	if bannerRequest.TagID != constants.NoValue && bannerRequest.FeatureID != constants.NoValue {
		row = d.Conn.QueryRow(ctx, "select b.content from banners b join banners_tags bt on b.id = bt.banner_id where bt.tag_id = $1 and b.feature_id = $2 and b.is_active = true order by b.id desc limit 1", bannerRequest.TagID, bannerRequest.FeatureID)
	} else if bannerRequest.TagID != constants.NoValue && bannerRequest.FeatureID == constants.NoValue {
		row = d.Conn.QueryRow(ctx, "select b.content from banners b join banners_tags bt on b.id = bt.banner_id where bt.tag_id = $1 and b.is_active = true order by b.id desc limit 1", bannerRequest.TagID)
	} else {
		row = d.Conn.QueryRow(ctx, "select content from banners where feature_id = $1 and is_active = true order by id desc limit 1", bannerRequest.FeatureID)
	}
	banner := map[string]interface{}{}
	err := row.Scan(&banner)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.BannerNotFound
		}
		return nil, errors.Wrap(err, "select banner failed")
	}
	return banner, nil
}

// GetAllBanners принимает запрос на банер, подразумевается, что будет либо tagID, либо featureID, либо оба
func (d *DB) GetAllBanners(ctx context.Context, bannerRequest *models.BannerRequest) ([]models.Banner, error) {
	log := d.log.With(slog.String("method", "GetAllBanners"))
	log.Info("GetAllBanners", slog.Any("bannerRequest", bannerRequest))
	var rows pgx.Rows
	var err error
	if bannerRequest.TagID != constants.NoValue && bannerRequest.FeatureID != constants.NoValue {
		rows, err = d.Conn.Query(ctx, "select b.id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at from banners b join banners_tags bt on b.id = bt.banner_id where bt.tag_id = $1 and b.feature_id = $2 group by b.id order by b.id desc limit $3 offset $4", bannerRequest.TagID, bannerRequest.FeatureID, bannerRequest.Limit, bannerRequest.Offset)
	} else if bannerRequest.TagID != constants.NoValue && bannerRequest.FeatureID == constants.NoValue {
		rows, err = d.Conn.Query(ctx, "select b.id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at from banners b join banners_tags bt on b.id = bt.banner_id where bt.tag_id = $1 order by b.id desc limit $2 offset $3", bannerRequest.TagID, bannerRequest.Limit, bannerRequest.Offset)
	} else {
		rows, err = d.Conn.Query(ctx, "select id, feature_id, content, is_active, created_at, updated_at from banners where feature_id = $1 order by id desc limit $2 offset $3", bannerRequest.FeatureID, bannerRequest.Limit, bannerRequest.Offset)
	}
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.BannerNotFound
		}
		return nil, errors.Wrap(err, "select banner failed")
	}
	banners := make([]models.Banner, 0)
	for rows.Next() {
		banner := models.Banner{}
		err = rows.Scan(&banner.ID, &banner.Feature, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			log.Error(err.Error())
			return nil, errors.Wrap(err, "select banner failed")
		}
		banners = append(banners, banner)
	}
	for i, banner := range banners {
		rows, err = d.Conn.Query(ctx, "select tag_id from banners_tags where banner_id = $1", banner.ID)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		tags := make([]int, 0)
		for rows.Next() {
			var tag int
			err = rows.Scan(&tag)
			if err != nil {
				log.Error(err.Error())
				return nil, errors.Wrap(err, "select banner tags failed")
			}
			tags = append(tags, tag)
		}
		banners[i].Tags = tags
	}
	return banners, nil
}

func (d *DB) Close() error {
	d.Conn.Close()
	return nil
}

func New(cfg *config.Config, log *slog.Logger) (*DB, error) {
	conn, err := initDatabase(cfg, log)
	if err != nil {
		return nil, errors.Wrap(err, "InitDatabase failed")
	}
	return &DB{Conn: conn, log: log}, nil
}

func onConnectScript(conn *pgxpool.Pool, log *slog.Logger) error {
	log = log.With(slog.String("function", "onConnectScript"))
	ctx, cl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cl()
	file, err := os.Open(filepath.Join(".", "config", "init.sql"))
	if err != nil {
		log.Error(err.Error())
		return errors.Wrap(err, "os.Open failed")
	}
	sqlScript, err := io.ReadAll(file)
	if err != nil {
		log.Error(err.Error())
		return errors.Wrap(err, "io.ReadAll failed")
	}
	_, err = conn.Exec(ctx, string(sqlScript))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func initDatabase(cfg *config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	log = log.With(slog.String("function", "initDatabase"))
	ctx, cl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cl()
	pool, err := pgxpool.New(ctx, cfg.StoragePath)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.Wrap(err, "pgxpool.New failed")
	}
	err = onConnectScript(pool, log)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.Wrap(err, "onConnectScript failed")
	}
	return pool, nil
}
