package persistence

import (
	"context"
	"errors"
	"github.com/devararishivian/antrekuy/internal/domain/entity"
	"github.com/devararishivian/antrekuy/internal/domain/repository"
	"github.com/devararishivian/antrekuy/internal/infrastructure"
	"github.com/jackc/pgx/v5"
	"log"
)

type UserRepositoryImpl struct {
	db *infrastructure.Database
}

func NewUserRepository(db *infrastructure.Database) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (u *UserRepositoryImpl) Store(user *entity.User, roleID int) (*entity.User, error) {
	const (
		insertUserStmt string = `INSERT INTO "user"(id, name, email, password) 
						VALUES ($1, $2, $3, $4) 
						RETURNING id, name, email, created_at, updated_at`
		insertUserRole string = `INSERT INTO user_has_role(user_id, role_id) VALUES ($1, $2)`
	)

	existingUser, err := u.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if existingUser.ID != "" {
		return nil, errors.New("user with given email already exists")
	}

	res := new(entity.User)

	tx, err := u.db.Conn.Begin(context.Background())
	if err != nil {
		return res, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println(err)
		}
	}(tx, context.Background())

	err = tx.QueryRow(
		context.Background(),
		insertUserStmt,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Email,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return res, err
	}

	_, err = tx.Exec(context.Background(), insertUserRole, res.ID, roleID)
	if err != nil {
		return res, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u *UserRepositoryImpl) FindByEmail(email string) (res entity.User, err error) {
	const sql string = `SELECT id, name, email, created_at, updated_at FROM "user" WHERE email = $1`

	err = u.db.Conn.QueryRow(context.Background(), sql, email).Scan(
		&res.ID,
		&res.Name,
		&res.Email,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		if err.Error() != "no rows in result set" {
			return res, err
		}

		return res, nil
	}

	return res, nil
}
