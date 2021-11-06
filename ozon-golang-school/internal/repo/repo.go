package repo

import (
	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ozonva/ova-link-api/internal/link"
)

type Repo interface {
	AddEntities(entities []link.Link) error
	ListEntities(limit uint64, offset uint64) ([]link.Link, error)
	DescribeEntity(entityId uint64) (*link.Link, error)
	DeleteEntity(entityId uint64) error
}

type LinkRepo struct {
	db *sqlx.DB
}

func NewLinkRepo(db *sqlx.DB) *LinkRepo {
	return &LinkRepo{
		db: db,
	}
}

func (lp *LinkRepo) AddEntities(entities []link.Link) error {
	sqlBuilder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("links").
		Columns("user_id", "url", "description", "tags")

	for _, entity := range entities {
		sqlBuilder = sqlBuilder.Values(entity.UserID, entity.Url, entity.Description, entity.Tags)
	}

	_, err := sqlBuilder.RunWith(lp.db).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (lp *LinkRepo) ListEntities(limit uint64, offset uint64) ([]link.Link, error) {
	sqlBuilder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "user_id", "url", "description", "tags", "created_at").
		From("links").
		Limit(limit).Offset(offset)

	sql, params, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result := make([]link.Link, 0, 0)
	err = lp.db.Select(&result, sql, params...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (lp *LinkRepo) DescribeEntity(entityId uint64) (*link.Link, error) {
	sqlBuilder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "user_id", "url", "description", "tags", "created_at").
		From("links").
		Where("id = ?")

	sql, _, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result := &link.Link{}
	err = lp.db.Get(result, sql, entityId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (lp *LinkRepo) DeleteEntity(entityId uint64) error {
	sqlBuilder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Delete("links").
		Where("id = ?")

	sql, _, err := sqlBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = lp.db.Exec(sql, entityId)
	if err != nil {
		return err
	}

	return nil
}
