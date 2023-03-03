package persistence

import (
	"context"
	"errors"
	"github.com/devararishivian/antrekuy/internal/domain/entity"
	"github.com/devararishivian/antrekuy/internal/domain/repository"
	"github.com/devararishivian/antrekuy/internal/infrastructure"
)

type UserRepositoryImpl struct {
	db *infrastructure.Database
}

func NewUserRepository(db *infrastructure.Database) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (u *UserRepositoryImpl) Store(user *entity.User) (*entity.User, error) {
	var insertUserStmt = `INSERT INTO "user"(id, name, email, password, role_id) 
							VALUES ($1, $2, $3, $4, $5) 
							RETURNING id, name, email, role_id, created_at, updated_at`

	existingUser, err := u.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if existingUser.ID != "" {
		return nil, errors.New("user with given email already exists")
	}

	result := new(entity.User)
	err = u.db.Conn.QueryRow(
		context.Background(),
		insertUserStmt,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Role.ID,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Role.ID,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *UserRepositoryImpl) FindByEmail(email string) (result entity.User, err error) {
	var selectUserStmt = `SELECT id, name, email, password, created_at, updated_at FROM "user" WHERE email = $1`

	err = u.db.Conn.QueryRow(context.Background(), selectUserStmt, email).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if err.Error() != "no rows in result set" {
			return result, err
		}

		return result, nil
	}

	return result, nil
}
