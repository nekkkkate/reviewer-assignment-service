package persistence

import (
	"database/sql"
	"errors"
	"regexp"
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"reviewer-assignment-service/internal/infrastructure/persistence/postgres"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamDataBase_GetByID(t *testing.T) {
	t.Run("successful get by id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)

		expectedTeam := &models.Team{
			ID:   1,
			Name: "backend",
			Members: map[int]*models.TeamMember{
				1: {UserID: 1, Username: "John", IsActive: true},
				2: {UserID: 2, Username: "Jane", IsActive: true},
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name FROM teams WHERE id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "backend"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.name, u.is_active FROM users u JOIN teams tm ON u.team_name = tm.name WHERE tm.id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active"}).
				AddRow(1, "John", true).
				AddRow(2, "Jane", true))

		team, err := teamDB.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedTeam.ID, team.ID)
		assert.Equal(t, expectedTeam.Name, team.Name)
		assert.Len(t, team.Members, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("team not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name FROM teams WHERE id = $1`)).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		team, err := teamDB.GetByID(999)
		assert.Nil(t, team)
		assert.ErrorIs(t, err, repositories.ErrTeamNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("team with no members", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name FROM teams WHERE id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "backend"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.name, u.is_active FROM users u JOIN teams tm ON u.team_name = tm.name WHERE tm.id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active"}))

		team, err := teamDB.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, team.ID)
		assert.Equal(t, "backend", team.Name)
		assert.Empty(t, team.Members)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeamDataBase_GetByName(t *testing.T) {
	t.Run("successful get by name", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name FROM teams WHERE name = $1`)).
			WithArgs("backend").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "backend"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.name, u.is_active FROM users u JOIN teams tm ON u.team_name = tm.name WHERE tm.name = $1`)).
			WithArgs("backend").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_active"}).
				AddRow(1, "John", true))

		team, err := teamDB.GetByName("backend")
		assert.NoError(t, err)
		assert.Equal(t, 1, team.ID)
		assert.Equal(t, "backend", team.Name)
		assert.Len(t, team.Members, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("team name not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name FROM teams WHERE name = $1`)).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		team, err := teamDB.GetByName("nonexistent")
		assert.Nil(t, team)
		assert.ErrorIs(t, err, repositories.ErrTeamNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeamDataBase_GetAll(t *testing.T) {
	t.Run("empty teams result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name FROM teams`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		teams, err := teamDB.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, teams)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestTeamDataBase_Update(t *testing.T) {
	t.Run("successful update team name only", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)
		team := &models.Team{
			ID:   1,
			Name: "backend-updated",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE teams SET name = $1 WHERE id = $2`)).
			WithArgs("backend-updated", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err = teamDB.Update(team)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("team not found for update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)
		team := &models.Team{
			ID:   999,
			Name: "nonexistent",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE teams SET name = $1 WHERE id = $2`)).
			WithArgs("nonexistent", 999).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectRollback()

		err = teamDB.Update(team)
		assert.ErrorIs(t, err, repositories.ErrTeamNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate team name error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)
		team := &models.Team{
			ID:   1,
			Name: "existing-name",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE teams SET name = $1 WHERE id = $2`)).
			WithArgs("existing-name", 1).
			WillReturnError(errors.New("pq: duplicate key value violates unique constraint"))
		mock.ExpectRollback()

		err = teamDB.Update(team)
		assert.ErrorIs(t, err, repositories.ErrTeamAlreadyExists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeamDataBase_AddUserToTeam(t *testing.T) {
	t.Run("team not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 1 FROM teams WHERE id = $1`)).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		err = teamDB.AddUserToTeam(999, 1)
		assert.ErrorIs(t, err, repositories.ErrTeamNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		teamDB := postgres.NewTeamDataBase(db)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 1 FROM teams WHERE id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 1 FROM users WHERE id = $1`)).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		err = teamDB.AddUserToTeam(1, 999)
		assert.ErrorIs(t, err, repositories.ErrUserNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
