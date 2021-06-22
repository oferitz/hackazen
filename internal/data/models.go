package data

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("no rows in result set")
	ErrEditConflict   = errors.New("edit conflict")
	pgErr             *pgconn.PgError
)

type Models struct {
	//Issues      IssueModel
	//Permissions PermissionModel
	Tokens TokenModel
	Users  UserModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		//Issues:      IssueModel{DB: db},
		//Permissions: PermissionModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}

func IsUniqueConstraintError(err error) bool {
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
