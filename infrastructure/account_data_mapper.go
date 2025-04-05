package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/freerware/morph"
	"github.com/freerware/tutor/domain"
	"github.com/freerware/work/v4/unit"
	"github.com/gofrs/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Errors that are potentially thrown during data mapper interactions.
var (

	// ErrInvalidType represents an error that indicates a unexpected type
	// was provided to the data mapper.
	ErrInvalidType = errors.New("infrastructure: invalid type provided to data mapper")
)

type AccountDataMapperParameters struct {
	fx.In

	DB     *sql.DB `name:"rwDB"`
	Logger *zap.Logger
}

type AccountDataMapper struct {
	db           *sql.DB
	logger       *zap.Logger
	accountTable morph.Table
	postsTable   morph.Table
}

func NewAccountDataMapper(parameters AccountDataMapperParameters) AccountDataMapper {
	opts := []morph.ReflectOption{
		morph.WithPrimaryKeyColumn("UUID"),
		morph.WithInferredTableName(morph.ScreamingSnakeCaseStrategy, false),
		morph.WithInferredColumnNames(morph.ScreamingSnakeCaseStrategy),
		morph.WithInferredTableAlias(morph.UpperCaseStrategy, 1),
		morph.WithColumnNameMapping("Username", "PRIMARY_CREDENTIAL"),
		morph.WithoutMethods("HasPost", "Posts", "AddPost", "AddPosts"),
	}
	at := morph.Must(morph.Reflect(domain.Account{}, opts...))

	opts = []morph.ReflectOption{
		morph.WithPrimaryKeyColumn("UUID"),
		morph.WithInferredTableName(morph.ScreamingSnakeCaseStrategy, false),
		morph.WithInferredColumnNames(morph.ScreamingSnakeCaseStrategy),
		morph.WithInferredTableAlias(morph.UpperCaseStrategy, 1),
		morph.WithColumnNameMapping("Likes", "LIKE_COUNT"),
		morph.WithColumnNameMapping("IsDraft", "DRAFT"),
		morph.WithoutMethods("Publish", "IncLikes"),
	}
	pt := morph.Must(morph.Reflect(domain.Post{}, opts...))

	return AccountDataMapper{
		db:           parameters.DB,
		logger:       parameters.Logger,
		accountTable: at,
		postsTable:   pt,
	}
}

func (dm *AccountDataMapper) FindPosts(ctx context.Context, mCtx unit.MapperContext, accountUUID uuid.UUID) ([]domain.Post, error) {
	query := "SELECT " + strings.Join(dm.postsTable.ColumnNames(), ", ") + " FROM " + dm.postsTable.Name() + " WHERE " + morph.Must(dm.postsTable.ColumnName("AuthorUUID")) + " = ?;"
	stmt, err := mCtx.Tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, accountUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []domain.Post{}
	for rows.Next() {
		var params domain.PostParameters
		err = rows.Scan(
			&params.AuthorUUID,
			&params.Content,
			&params.CreatedAt,
			&params.DeletedAt,
			&params.Draft,
			&params.Likes,
			&params.Title,
			&params.UpdatedAt,
			&params.UUID,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, domain.ReconstitutePost(params))
	}

	return posts, nil
}

func (dm *AccountDataMapper) Find(ctx context.Context, mCtx unit.MapperContext, uuid uuid.UUID) (domain.Account, error) {
	sql, err := dm.accountTable.SelectQuery()
	if err != nil {
		return domain.Account{}, err
	}

	stmt, err := mCtx.Tx.Prepare(sql)
	if err != nil {
		return domain.Account{}, err
	}

	rows, err := stmt.QueryContext(ctx, uuid)
	if err != nil {
		stmt.Close()
		return domain.Account{}, err
	}

	var params domain.AccountParameters
	if rows.Next() {
		err = rows.Scan(
			&params.CreatedAt,
			&params.DeletedAt,
			&params.GivenName,
			&params.Username,
			&params.Surname,
			&params.UpdatedAt,
			&params.UUID,
		)
		if err != nil {
			rows.Close()
			return domain.Account{}, err
		}
	}
	rows.Close()
	stmt.Close()

	posts, err := dm.FindPosts(ctx, mCtx, params.UUID)
	if err != nil {
		return domain.Account{}, err
	}
	params.Posts = posts
	return domain.ReconstituteAccount(params), nil
}

func (dm *AccountDataMapper) Insert(ctx context.Context, mCtx unit.MapperContext, accounts ...any) error {
	for _, account := range accounts {
		sql, args, err := dm.accountTable.InsertQueryWithArgs(account)
		if err != nil {
			return err
		}

		stmt, err := mCtx.Tx.Prepare(sql)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, args...)
		if err != nil {
			return err
		}

		acc := account.(domain.Account)
		for _, post := range acc.Posts() {
			sql, args, err := dm.postsTable.InsertQueryWithArgs(post)
			if err != nil {
				return err
			}

			stmt, err := mCtx.Tx.Prepare(sql)
			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.ExecContext(ctx, args...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (dm *AccountDataMapper) Update(ctx context.Context, mCtx unit.MapperContext, accounts ...any) error {
	for _, account := range accounts {
		sql, args, err := dm.accountTable.UpdateQueryWithArgs(account)
		if err != nil {
			return err
		}

		stmt, err := mCtx.Tx.Prepare(sql)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, args...)
		if err != nil {
			return err
		}

		acc := account.(domain.Account)
		before, err := dm.Find(ctx, mCtx, acc.UUID())
		if err != nil {
			return err
		}

		for _, post := range acc.Posts() {
			// update.
			if before.HasPost(post) {
				sql, args, err := dm.postsTable.UpdateQueryWithArgs(post)
				if err != nil {
					return err
				}

				stmt, err := mCtx.Tx.Prepare(sql)
				if err != nil {
					return err
				}
				defer stmt.Close()

				_, err = stmt.ExecContext(ctx, args...)
				if err != nil {
					return err
				}
				continue
			}

			// insert.
			if !before.HasPost(post) {
				sql, args, err := dm.postsTable.InsertQueryWithArgs(post)
				if err != nil {
					return err
				}

				stmt, err := mCtx.Tx.Prepare(sql)
				if err != nil {
					return err
				}
				defer stmt.Close()

				_, err = stmt.ExecContext(ctx, args...)
				if err != nil {
					return err
				}
				continue
			}
		}

		for _, post := range before.Posts() {
			// delete.
			if !acc.HasPost(post) {
				sql, args, err := dm.postsTable.DeleteQueryWithArgs(post)
				if err != nil {
					return err
				}

				stmt, err := mCtx.Tx.Prepare(sql)
				if err != nil {
					return err
				}
				defer stmt.Close()

				_, err = stmt.ExecContext(ctx, args...)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (dm *AccountDataMapper) Delete(ctx context.Context, mCtx unit.MapperContext, accounts ...any) error {
	for _, account := range accounts {
		sql, args, err := dm.accountTable.DeleteQueryWithArgs(account)
		if err != nil {
			return err
		}

		stmt, err := mCtx.Tx.Prepare(sql)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, args...)
		if err != nil {
			return err
		}

		acc := account.(domain.Account)
		for _, post := range acc.Posts() {
			sql, args, err := dm.postsTable.DeleteQueryWithArgs(post)
			if err != nil {
				return err
			}

			stmt, err := mCtx.Tx.Prepare(sql)
			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.ExecContext(ctx, args...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
