package persistence

import (
	"database/sql"
	"regexp"
	"reviewer-assignment-service/internal/domain/repositories"
	"reviewer-assignment-service/internal/infrustructure/persistence/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullRequestDataBase_GetByID(t *testing.T) {
	t.Run("pr not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		prDB := postgres.NewPullRequestDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.title, p.status, p.created_at, p.merged_at, u.id, u.name, u.email, u.team_name, u.is_active FROM prs p JOIN users u ON p.author_id = u.id WHERE p.id = $1`)).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		pr, err := prDB.GetByID(999)
		assert.Nil(t, pr)
		assert.ErrorIs(t, err, repositories.ErrPullRequestNotFoundInPersistence)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("pr without reviewers", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		prDB := postgres.NewPullRequestDataBase(db)
		createdAt := time.Now()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.title, p.status, p.created_at, p.merged_at, u.id, u.name, u.email, u.team_name, u.is_active FROM prs p JOIN users u ON p.author_id = u.id WHERE p.id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"p.id", "p.title", "p.status", "p.created_at", "p.merged_at", "u.id", "u.name", "u.email", "u.team_name", "u.is_active"}).
				AddRow(1, "Test PR", "open", createdAt, nil, 1, "User 1", "user1@test.com", "Team A", true))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.name, u.email, u.team_name, u.is_active FROM assigned_reviewers ar JOIN users u ON ar.user_id = u.id WHERE ar.pr_id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"u.id", "u.name", "u.email", "u.team_name", "u.is_active"}))

		pr, err := prDB.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, pr.ID)
		assert.Empty(t, pr.Reviewers)
		assert.True(t, pr.MergedAt.IsZero())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPullRequestDataBase_GetAll(t *testing.T) {
	t.Run("successful get all", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		prDB := postgres.NewPullRequestDataBase(db)
		createdAt1 := time.Now()
		createdAt2 := time.Now().Add(-time.Hour)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.title, p.status, p.created_at, p.merged_at, u.id, u.name, u.email, u.team_name, u.is_active FROM prs p JOIN users u ON p.author_id = u.id ORDER BY p.created_at DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{"p.id", "p.title", "p.status", "p.created_at", "p.merged_at", "u.id", "u.name", "u.email", "u.team_name", "u.is_active"}).
				AddRow(1, "PR 1", "open", createdAt1, nil, 1, "User 1", "user1@test.com", "Team A", true).
				AddRow(2, "PR 2", "merged", createdAt2, createdAt2.Add(time.Hour), 2, "User 2", "user2@test.com", "Team B", true))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT ar.pr_id, u.id, u.name, u.email, u.team_name, u.is_active FROM assigned_reviewers ar JOIN users u ON ar.user_id = u.id WHERE ar.pr_id IN ($1,$2)`)).
			WithArgs(1, 2).
			WillReturnRows(sqlmock.NewRows([]string{"ar.pr_id", "u.id", "u.name", "u.email", "u.team_name", "u.is_active"}).
				AddRow(1, 3, "Reviewer 1", "reviewer1@test.com", "Team A", true).
				AddRow(2, 4, "Reviewer 2", "reviewer2@test.com", "Team B", true))

		prs, err := prDB.GetAll()
		assert.NoError(t, err)
		assert.Len(t, prs, 2)
		assert.Equal(t, 1, prs[0].ID)
		assert.Equal(t, 2, prs[1].ID)
		assert.Len(t, prs[0].Reviewers, 1)
		assert.Len(t, prs[1].Reviewers, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		prDB := postgres.NewPullRequestDataBase(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.title, p.status, p.created_at, p.merged_at, u.id, u.name, u.email, u.team_name, u.is_active FROM prs p JOIN users u ON p.author_id = u.id ORDER BY p.created_at DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{"p.id", "p.title", "p.status", "p.created_at", "p.merged_at", "u.id", "u.name", "u.email", "u.team_name", "u.is_active"}))

		prs, err := prDB.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, prs)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPullRequestDataBase_GetByAuthorID(t *testing.T) {
	t.Run("successful get by author id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		prDB := postgres.NewPullRequestDataBase(db)
		createdAt := time.Now()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.title, p.status, p.created_at, p.merged_at, u.id, u.name, u.email, u.team_name, u.is_active FROM prs p JOIN users u ON p.author_id = u.id WHERE p.author_id = $1 ORDER BY p.created_at DESC`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"p.id", "p.title", "p.status", "p.created_at", "p.merged_at", "u.id", "u.name", "u.email", "u.team_name", "u.is_active"}).
				AddRow(1, "Author PR", "open", createdAt, nil, 1, "User 1", "user1@test.com", "Team A", true))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT ar.pr_id, u.id, u.name, u.email, u.team_name, u.is_active FROM assigned_reviewers ar JOIN users u ON ar.user_id = u.id WHERE ar.pr_id IN ($1)`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"ar.pr_id", "u.id", "u.name", "u.email", "u.team_name", "u.is_active"}))

		prs, err := prDB.GetByAuthorID(1)
		assert.NoError(t, err)
		assert.Len(t, prs, 1)
		assert.Equal(t, 1, prs[0].Author.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPullRequestDataBase_GetByReviewerID(t *testing.T) {
	t.Run("successful get by reviewer id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		prDB := postgres.NewPullRequestDataBase(db)
		createdAt := time.Now()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.title, p.status, p.created_at, p.merged_at, u.id, u.name, u.email, u.team_name, u.is_active FROM prs p JOIN users u ON p.author_id = u.id JOIN assigned_reviewers ar ON p.id = ar.pr_id WHERE ar.user_id = $1 ORDER BY p.created_at DESC`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"p.id", "p.title", "p.status", "p.created_at", "p.merged_at", "u.id", "u.name", "u.email", "u.team_name", "u.is_active"}).
				AddRow(1, "Reviewed PR", "open", createdAt, nil, 1, "User 1", "user1@test.com", "Team A", true))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT ar.pr_id, u.id, u.name, u.email, u.team_name, u.is_active FROM assigned_reviewers ar JOIN users u ON ar.user_id = u.id WHERE ar.pr_id IN ($1)`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"ar.pr_id", "u.id", "u.name", "u.email", "u.team_name", "u.is_active"}).
				AddRow(1, 2, "Reviewer", "reviewer@test.com", "Team A", true))

		prs, err := prDB.GetByReviewerID(2)
		assert.NoError(t, err)
		assert.Len(t, prs, 1)
		assert.Len(t, prs[0].Reviewers, 1)
		assert.Equal(t, 2, prs[0].Reviewers[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
