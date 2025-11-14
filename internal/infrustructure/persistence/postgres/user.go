package postgres

import (
	"database/sql"
	"errors"
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"

	"github.com/Masterminds/squirrel"
)

type UserDataBase struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func NewUserDataBase(db *sql.DB) *UserDataBase {
	return &UserDataBase{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (u *UserDataBase) Add(user *models.User) error {
	query, args, err := u.sb.
		Insert("users").
		Columns("name", "email", "team_name", "is_active").
		Values(user.Name, user.Email, user.TeamName, user.IsActive).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return err
	}

	err = u.db.QueryRow(query, args...).Scan(&user.ID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			return repositories.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (u *UserDataBase) GetByID(id int) (*models.User, error) {
	query, args, err := u.sb.
		Select("id", "name", "email", "team_name", "is_active").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	user := &models.User{}
	err = u.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Name, &user.Email, &user.TeamName, &user.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrUserNotFoundInPersistence
		}
		return nil, err
	}
	return user, nil
}

func (u *UserDataBase) GetByEmail(email string) (*models.User, error) {
	query, args, err := u.sb.
		Select("id", "name", "email", "team_name", "is_active").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		return nil, err
	}

	user := &models.User{}
	err = u.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Name, &user.Email, &user.TeamName, &user.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrUserWithThatEmailNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *UserDataBase) GetAll() ([]*models.User, error) {
	query, args, err := u.sb.
		Select("id", "name", "email", "team_name", "is_active").
		From("users").
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.TeamName, &user.IsActive,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserDataBase) GetWithFilters(teamName string, isActive bool) ([]*models.User, error) {
	builder := u.sb.
		Select("id", "name", "email", "team_name", "is_active").
		From("users")

	if teamName != "" {
		builder = builder.Where(squirrel.Eq{"team_name": teamName})
	}

	builder = builder.Where(squirrel.Eq{"is_active": isActive})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.TeamName, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserDataBase) GetActiveUsers() ([]*models.User, error) {
	query, args, err := u.sb.
		Select("id", "name", "email", "team_name", "is_active").
		From("users").
		Where(squirrel.Eq{"is_active": true}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.TeamName, &user.IsActive,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserDataBase) Update(user *models.User) error {
	query, args, err := u.sb.
		Update("users").
		Set("name", user.Name).
		Set("email", user.Email).
		Set("team_name", user.TeamName).
		Set("is_active", user.IsActive).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()

	if err != nil {
		return err
	}

	result, err := u.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repositories.ErrUserNotFoundInPersistence
	}

	return nil
}

func (u *UserDataBase) Deactivate(userID int) error {
	query, args, err := u.sb.
		Update("users").
		Set("is_active", false).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		return err
	}

	result, err := u.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repositories.ErrUserNotFoundInPersistence
	}

	return nil
}
