package persistence

import (
	"database/sql"
	"errors"
	"regexp"
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"reviewer-assignment-service/internal/infrustructure/persistence/postgres"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDataBase_GetByID(t *testing.T) {
	t.Run("successful get by id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		expectedUser := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "team_name", "is_active"}).
				AddRow(1, "John Doe", "john@example.com", "backend", true))

		user, err := userDB.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE id = $1`)).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		user, err := userDB.GetByID(999)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, repositories.ErrUserNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE id = $1`)).
			WithArgs(1).
			WillReturnError(errors.New("connection failed"))

		user, err := userDB.GetByID(1)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserDataBase_GetByEmail(t *testing.T) {
	t.Run("successful get by email", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		expectedUser := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE email = $1`)).
			WithArgs("john@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "team_name", "is_active"}).
				AddRow(1, "John Doe", "john@example.com", "backend", true))

		user, err := userDB.GetByEmail("john@example.com")
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("email not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE email = $1`)).
			WithArgs("nonexistent@example.com").
			WillReturnError(sql.ErrNoRows)

		user, err := userDB.GetByEmail("nonexistent@example.com")
		assert.Nil(t, user)
		assert.ErrorIs(t, err, repositories.ErrUserWithThatEmailNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserDataBase_GetAll(t *testing.T) {
	t.Run("successful get all users", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		expectedUsers := []*models.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", TeamName: "backend", IsActive: true},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", TeamName: "frontend", IsActive: true},
			{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", TeamName: "backend", IsActive: false},
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "team_name", "is_active"}).
				AddRow(1, "John Doe", "john@example.com", "backend", true).
				AddRow(2, "Jane Smith", "jane@example.com", "frontend", true).
				AddRow(3, "Bob Johnson", "bob@example.com", "backend", false))

		users, err := userDB.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "team_name", "is_active"}))

		users, err := userDB.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserDataBase_GetActiveUsers(t *testing.T) {
	t.Run("successful get active users", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		expectedUsers := []*models.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", TeamName: "backend", IsActive: true},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", TeamName: "frontend", IsActive: true},
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE is_active = $1`)).
			WithArgs(true).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "team_name", "is_active"}).
				AddRow(1, "John Doe", "john@example.com", "backend", true).
				AddRow(2, "Jane Smith", "jane@example.com", "frontend", true))

		users, err := userDB.GetActiveUsers()
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserDataBase_GetWithFilters(t *testing.T) {
	t.Run("filter by team name and inactive users", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		expectedUsers := []*models.User{
			{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", TeamName: "backend", IsActive: false},
		}

		// Теперь ожидаем два условия: team_name AND is_active
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE team_name = $1 AND is_active = $2`)).
			WithArgs("backend", false).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "team_name", "is_active"}).
				AddRow(3, "Bob Johnson", "bob@example.com", "backend", false))

		users, err := userDB.GetWithFilters("backend", false)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("filter by active status only", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		expectedUsers := []*models.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", TeamName: "backend", IsActive: true},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", TeamName: "frontend", IsActive: true},
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE is_active = $1`)).
			WithArgs(true).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "team_name", "is_active"}).
				AddRow(1, "John Doe", "john@example.com", "backend", true).
				AddRow(2, "Jane Smith", "jane@example.com", "frontend", true))

		users, err := userDB.GetWithFilters("", true)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("filter by team and active status", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		expectedUsers := []*models.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", TeamName: "backend", IsActive: true},
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, team_name, is_active FROM users WHERE team_name = $1 AND is_active = $2`)).
			WithArgs("backend", true).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "team_name", "is_active"}).
				AddRow(1, "John Doe", "john@example.com", "backend", true))

		users, err := userDB.GetWithFilters("backend", true)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserDataBase_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		user := &models.User{
			ID:       1,
			Name:     "John Doe Updated",
			Email:    "john.updated@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET name = $1, email = $2, team_name = $3, is_active = $4 WHERE id = $5`)).
			WithArgs("John Doe Updated", "john.updated@example.com", "backend", true, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = userDB.Update(user)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found for update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		user := &models.User{
			ID:       999,
			Name:     "Nonexistent User",
			Email:    "nonexistent@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET name = $1, email = $2, team_name = $3, is_active = $4 WHERE id = $5`)).
			WithArgs("Nonexistent User", "nonexistent@example.com", "backend", true, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = userDB.Update(user)
		assert.ErrorIs(t, err, repositories.ErrUserNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserDataBase_Deactivate(t *testing.T) {
	t.Run("successful deactivation", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET is_active = $1 WHERE id = $2`)).
			WithArgs(false, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = userDB.Deactivate(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found for deactivation", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		userDB := postgres.NewUserDataBase(db)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET is_active = $1 WHERE id = $2`)).
			WithArgs(false, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = userDB.Deactivate(999)
		assert.ErrorIs(t, err, repositories.ErrUserNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
