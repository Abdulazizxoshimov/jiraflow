package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type favoriteRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewFavoriteRepo(p *pg.Postgres) repository.FavoriteRepository {
	return &favoriteRepo{db: p.DB, builder: p.Builder}
}

func (r *favoriteRepo) Add(ctx context.Context, fav *entity.Favorite) error {
	if fav.ID == "" {
		fav.ID = uuid.NewString()
	}
	fav.CreatedAt = time.Now().UTC()
	_, err := r.db.Exec(ctx,
		`INSERT INTO user_favorites(id, user_id, entity_type, entity_id, created_at)
		 VALUES($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING`,
		fav.ID, fav.UserID, fav.EntityType, fav.EntityID, fav.CreatedAt,
	)
	return err
}

func (r *favoriteRepo) Remove(ctx context.Context, userID, entityType, entityID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM user_favorites WHERE user_id=$1 AND entity_type=$2 AND entity_id=$3`,
		userID, entityType, entityID,
	)
	return err
}

func (r *favoriteRepo) List(ctx context.Context, userID string, filter *entity.FavoriteFilter) ([]*entity.Favorite, int, error) {
	where := sq.And{sq.Eq{"user_id": userID}}
	if filter.EntityType != "" {
		where = append(where, sq.Eq{"entity_type": filter.EntityType})
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("user_favorites").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("favoriteRepo.List count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select("id", "user_id", "entity_type", "entity_id", "created_at").
		From("user_favorites").Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("favoriteRepo.List query: %w", err)
	}
	defer rows.Close()

	var favs []*entity.Favorite
	for rows.Next() {
		f := &entity.Favorite{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.EntityType, &f.EntityID, &f.CreatedAt); err != nil {
			return nil, 0, err
		}
		favs = append(favs, f)
	}
	return favs, total, rows.Err()
}

func (r *favoriteRepo) Exists(ctx context.Context, userID, entityType, entityID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM user_favorites WHERE user_id=$1 AND entity_type=$2 AND entity_id=$3)`,
		userID, entityType, entityID,
	).Scan(&exists)
	return exists, err
}
