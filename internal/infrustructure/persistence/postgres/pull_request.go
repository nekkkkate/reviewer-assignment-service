package postgres

import (
	"database/sql"
	"errors"
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"strings"

	"github.com/Masterminds/squirrel"
)

type PullRequestDataBase struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func NewPullRequestDataBase(db *sql.DB) *PullRequestDataBase {
	return &PullRequestDataBase{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (p *PullRequestDataBase) Add(pr *models.PullRequest) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	teamQuery, teamArgs, err := p.sb.
		Select("team_id").
		From("users").
		Where(squirrel.Eq{"id": pr.Author.ID}).
		ToSql()

	if err != nil {
		return err
	}

	var teamID int
	err = tx.QueryRow(teamQuery, teamArgs...).Scan(&teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repositories.ErrUserNotFoundInPersistence
		}
		return err
	}

	var mergedAt interface{}
	if !pr.MergedAt.IsZero() {
		mergedAt = pr.MergedAt
	} else {
		mergedAt = nil
	}

	prQuery, prArgs, err := p.sb.
		Insert("prs").
		Columns("title", "author_id", "team_id", "status", "created_at", "merged_at").
		Values(pr.Name, pr.Author.ID, teamID, string(pr.Status), pr.CreatedAt, mergedAt).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return err
	}

	err = tx.QueryRow(prQuery, prArgs...).Scan(&pr.ID)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			if strings.Contains(err.Error(), "author_id") {
				return repositories.ErrUserNotFoundInPersistence
			}
			if strings.Contains(err.Error(), "team_id") {
				return repositories.ErrTeamNotFoundInPersistence
			}
		}
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			return repositories.ErrPullRequestAlreadyExists
		}
		return err
	}

	if len(pr.Reviewers) > 0 {
		for _, reviewer := range pr.Reviewers {
			reviewerQuery, reviewerArgs, err := p.sb.
				Insert("assigned_reviewers").
				Columns("pr_id", "user_id").
				Values(pr.ID, reviewer.ID).
				ToSql()

			if err != nil {
				return err
			}

			_, err = tx.Exec(reviewerQuery, reviewerArgs...)
			if err != nil {
				if err.Error() == "pq: duplicate key value violates unique constraint" {
					return models.ErrReviewerAlreadyAssigned
				}
				if strings.Contains(err.Error(), "violates foreign key constraint") {
					return repositories.ErrUserNotFoundInPersistence
				}
				return err
			}
		}
	}

	return tx.Commit()
}

func (p *PullRequestDataBase) GetByID(id int) (*models.PullRequest, error) {
	prQuery, prArgs, err := p.sb.
		Select("p.id", "p.title", "p.status", "p.created_at", "p.merged_at",
			"u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("prs p").
		Join("users u ON p.author_id = u.id").
		Where(squirrel.Eq{"p.id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	pr := &models.PullRequest{}
	var status string
	var mergedAt sql.NullTime
	author := &models.User{}

	row := p.db.QueryRow(prQuery, prArgs...)
	err = row.Scan(
		&pr.ID, &pr.Name, &status, &pr.CreatedAt, &mergedAt,
		&author.ID, &author.Name, &author.Email, &author.TeamName, &author.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrPullRequestNotFoundInPersistence
		}
		return nil, err
	}

	pr.Author = author
	pr.Status = models.PRStatus(status)
	if mergedAt.Valid {
		pr.MergedAt = mergedAt.Time
	}

	reviewersQuery, reviewersArgs, err := p.sb.
		Select("u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("assigned_reviewers ar").
		Join("users u ON ar.user_id = u.id").
		Where(squirrel.Eq{"ar.pr_id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	reviewersRows, err := p.db.Query(reviewersQuery, reviewersArgs...)
	if err != nil {
		return nil, err
	}
	defer reviewersRows.Close()

	pr.Reviewers = make([]*models.User, 0)
	for reviewersRows.Next() {
		reviewer := &models.User{}
		err := reviewersRows.Scan(
			&reviewer.ID, &reviewer.Name, &reviewer.Email, &reviewer.TeamName, &reviewer.IsActive,
		)
		if err != nil {
			return nil, err
		}
		pr.Reviewers = append(pr.Reviewers, reviewer)
	}

	if err = reviewersRows.Err(); err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *PullRequestDataBase) GetAll() ([]*models.PullRequest, error) {
	prsQuery, prsArgs, err := p.sb.
		Select("p.id", "p.title", "p.status", "p.created_at", "p.merged_at",
			"u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("prs p").
		Join("users u ON p.author_id = u.id").
		OrderBy("p.created_at DESC").
		ToSql()

	if err != nil {
		return nil, err
	}

	prsRows, err := p.db.Query(prsQuery, prsArgs...)
	if err != nil {
		return nil, err
	}
	defer prsRows.Close()

	var prs []*models.PullRequest
	prsByID := make(map[int]*models.PullRequest)

	for prsRows.Next() {
		pr := &models.PullRequest{}
		var status string
		var mergedAt sql.NullTime
		author := &models.User{}

		err := prsRows.Scan(
			&pr.ID, &pr.Name, &status, &pr.CreatedAt, &mergedAt,
			&author.ID, &author.Name, &author.Email, &author.TeamName, &author.IsActive,
		)
		if err != nil {
			return nil, err
		}

		pr.Author = author
		pr.Status = models.PRStatus(status)
		if mergedAt.Valid {
			pr.MergedAt = mergedAt.Time
		}
		pr.Reviewers = make([]*models.User, 0)

		prs = append(prs, pr)
		prsByID[pr.ID] = pr
	}

	if err = prsRows.Err(); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return prs, nil
	}

	prIDs := make([]int, 0, len(prs))
	for _, pr := range prs {
		prIDs = append(prIDs, pr.ID)
	}

	reviewersQuery, reviewersArgs, err := p.sb.
		Select("ar.pr_id", "u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("assigned_reviewers ar").
		Join("users u ON ar.user_id = u.id").
		Where(squirrel.Eq{"ar.pr_id": prIDs}).
		ToSql()

	if err != nil {
		return nil, err
	}

	reviewersRows, err := p.db.Query(reviewersQuery, reviewersArgs...)
	if err != nil {
		return nil, err
	}
	defer reviewersRows.Close()

	for reviewersRows.Next() {
		var prID int
		reviewer := &models.User{}

		err := reviewersRows.Scan(
			&prID,
			&reviewer.ID, &reviewer.Name, &reviewer.Email, &reviewer.TeamName, &reviewer.IsActive,
		)
		if err != nil {
			return nil, err
		}

		if pr, exists := prsByID[prID]; exists {
			pr.Reviewers = append(pr.Reviewers, reviewer)
		}
	}

	if err = reviewersRows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

func (p *PullRequestDataBase) GetByStatus(status models.PRStatus) ([]*models.PullRequest, error) {
	prsQuery, prsArgs, err := p.sb.
		Select("p.id", "p.title", "p.status", "p.created_at", "p.merged_at",
			"u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("prs p").
		Join("users u ON p.author_id = u.id").
		Where(squirrel.Eq{"p.status": string(status)}).
		OrderBy("p.created_at DESC").
		ToSql()

	if err != nil {
		return nil, err
	}

	prsRows, err := p.db.Query(prsQuery, prsArgs...)
	if err != nil {
		return nil, err
	}
	defer prsRows.Close()

	var prs []*models.PullRequest
	prsByID := make(map[int]*models.PullRequest)

	for prsRows.Next() {
		pr := &models.PullRequest{}
		var statusStr string
		var mergedAt sql.NullTime
		author := &models.User{}

		err := prsRows.Scan(
			&pr.ID, &pr.Name, &statusStr, &pr.CreatedAt, &mergedAt,
			&author.ID, &author.Name, &author.Email, &author.TeamName, &author.IsActive,
		)
		if err != nil {
			return nil, err
		}

		pr.Author = author
		pr.Status = models.PRStatus(statusStr)
		if mergedAt.Valid {
			pr.MergedAt = mergedAt.Time
		}
		pr.Reviewers = make([]*models.User, 0)

		prs = append(prs, pr)
		prsByID[pr.ID] = pr
	}

	if err = prsRows.Err(); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return prs, nil
	}

	prIDs := make([]int, 0, len(prs))
	for _, pr := range prs {
		prIDs = append(prIDs, pr.ID)
	}

	reviewersQuery, reviewersArgs, err := p.sb.
		Select("ar.pr_id", "u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("assigned_reviewers ar").
		Join("users u ON ar.user_id = u.id").
		Where(squirrel.Eq{"ar.pr_id": prIDs}).
		ToSql()

	if err != nil {
		return nil, err
	}

	reviewersRows, err := p.db.Query(reviewersQuery, reviewersArgs...)
	if err != nil {
		return nil, err
	}
	defer reviewersRows.Close()

	for reviewersRows.Next() {
		var prID int
		reviewer := &models.User{}

		err := reviewersRows.Scan(
			&prID,
			&reviewer.ID, &reviewer.Name, &reviewer.Email, &reviewer.TeamName, &reviewer.IsActive,
		)
		if err != nil {
			return nil, err
		}

		if pr, exists := prsByID[prID]; exists {
			pr.Reviewers = append(pr.Reviewers, reviewer)
		}
	}

	if err = reviewersRows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

func (p *PullRequestDataBase) GetByAuthorID(authorID int) ([]*models.PullRequest, error) {
	prsQuery, prsArgs, err := p.sb.
		Select("p.id", "p.title", "p.status", "p.created_at", "p.merged_at",
			"u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("prs p").
		Join("users u ON p.author_id = u.id").
		Where(squirrel.Eq{"p.author_id": authorID}).
		OrderBy("p.created_at DESC").
		ToSql()

	if err != nil {
		return nil, err
	}

	prsRows, err := p.db.Query(prsQuery, prsArgs...)
	if err != nil {
		return nil, err
	}
	defer prsRows.Close()

	var prs []*models.PullRequest
	prsByID := make(map[int]*models.PullRequest)

	for prsRows.Next() {
		pr := &models.PullRequest{}
		var status string
		var mergedAt sql.NullTime
		author := &models.User{}

		err := prsRows.Scan(
			&pr.ID, &pr.Name, &status, &pr.CreatedAt, &mergedAt,
			&author.ID, &author.Name, &author.Email, &author.TeamName, &author.IsActive,
		)
		if err != nil {
			return nil, err
		}

		pr.Author = author
		pr.Status = models.PRStatus(status)
		if mergedAt.Valid {
			pr.MergedAt = mergedAt.Time
		}
		pr.Reviewers = make([]*models.User, 0)

		prs = append(prs, pr)
		prsByID[pr.ID] = pr
	}

	if err = prsRows.Err(); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return prs, nil
	}

	prIDs := make([]int, 0, len(prs))
	for _, pr := range prs {
		prIDs = append(prIDs, pr.ID)
	}

	reviewersQuery, reviewersArgs, err := p.sb.
		Select("ar.pr_id", "u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("assigned_reviewers ar").
		Join("users u ON ar.user_id = u.id").
		Where(squirrel.Eq{"ar.pr_id": prIDs}).
		ToSql()

	if err != nil {
		return nil, err
	}

	reviewersRows, err := p.db.Query(reviewersQuery, reviewersArgs...)
	if err != nil {
		return nil, err
	}
	defer reviewersRows.Close()

	for reviewersRows.Next() {
		var prID int
		reviewer := &models.User{}

		err := reviewersRows.Scan(
			&prID,
			&reviewer.ID, &reviewer.Name, &reviewer.Email, &reviewer.TeamName, &reviewer.IsActive,
		)
		if err != nil {
			return nil, err
		}

		if pr, exists := prsByID[prID]; exists {
			pr.Reviewers = append(pr.Reviewers, reviewer)
		}
	}

	if err = reviewersRows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

func (p *PullRequestDataBase) GetByReviewerID(reviewerID int) ([]*models.PullRequest, error) {
	prsQuery, prsArgs, err := p.sb.
		Select("p.id", "p.title", "p.status", "p.created_at", "p.merged_at",
			"u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("prs p").
		Join("users u ON p.author_id = u.id").
		Join("assigned_reviewers ar ON p.id = ar.pr_id").
		Where(squirrel.Eq{"ar.user_id": reviewerID}).
		OrderBy("p.created_at DESC").
		ToSql()

	if err != nil {
		return nil, err
	}

	prsRows, err := p.db.Query(prsQuery, prsArgs...)
	if err != nil {
		return nil, err
	}
	defer prsRows.Close()

	var prs []*models.PullRequest
	prsByID := make(map[int]*models.PullRequest)

	for prsRows.Next() {
		pr := &models.PullRequest{}
		var status string
		var mergedAt sql.NullTime
		author := &models.User{}

		err := prsRows.Scan(
			&pr.ID, &pr.Name, &status, &pr.CreatedAt, &mergedAt,
			&author.ID, &author.Name, &author.Email, &author.TeamName, &author.IsActive,
		)
		if err != nil {
			return nil, err
		}

		pr.Author = author
		pr.Status = models.PRStatus(status)
		if mergedAt.Valid {
			pr.MergedAt = mergedAt.Time
		}
		pr.Reviewers = make([]*models.User, 0)

		prs = append(prs, pr)
		prsByID[pr.ID] = pr
	}

	if err = prsRows.Err(); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return prs, nil
	}

	prIDs := make([]int, 0, len(prs))
	for _, pr := range prs {
		prIDs = append(prIDs, pr.ID)
	}

	reviewersQuery, reviewersArgs, err := p.sb.
		Select("ar.pr_id", "u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("assigned_reviewers ar").
		Join("users u ON ar.user_id = u.id").
		Where(squirrel.Eq{"ar.pr_id": prIDs}).
		ToSql()

	if err != nil {
		return nil, err
	}

	reviewersRows, err := p.db.Query(reviewersQuery, reviewersArgs...)
	if err != nil {
		return nil, err
	}
	defer reviewersRows.Close()

	for reviewersRows.Next() {
		var prID int
		reviewer := &models.User{}

		err := reviewersRows.Scan(
			&prID,
			&reviewer.ID, &reviewer.Name, &reviewer.Email, &reviewer.TeamName, &reviewer.IsActive,
		)
		if err != nil {
			return nil, err
		}

		if pr, exists := prsByID[prID]; exists {
			pr.Reviewers = append(pr.Reviewers, reviewer)
		}
	}

	if err = reviewersRows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

func (p *PullRequestDataBase) Update(pr *models.PullRequest) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var mergedAt interface{}
	if !pr.MergedAt.IsZero() {
		mergedAt = pr.MergedAt
	} else {
		mergedAt = nil
	}

	updateQuery, updateArgs, err := p.sb.
		Update("prs").
		Set("title", pr.Name).
		Set("status", string(pr.Status)).
		Set("merged_at", mergedAt).
		Where(squirrel.Eq{"id": pr.ID}).
		ToSql()

	if err != nil {
		return err
	}

	result, err := tx.Exec(updateQuery, updateArgs...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repositories.ErrPullRequestNotFoundInPersistence
	}

	if pr.Reviewers != nil {
		deleteQuery, deleteArgs, err := p.sb.
			Delete("assigned_reviewers").
			Where(squirrel.Eq{"pr_id": pr.ID}).
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.Exec(deleteQuery, deleteArgs...)
		if err != nil {
			return err
		}

		if len(pr.Reviewers) > 0 {
			for _, reviewer := range pr.Reviewers {
				reviewerQuery, reviewerArgs, err := p.sb.
					Insert("assigned_reviewers").
					Columns("pr_id", "user_id").
					Values(pr.ID, reviewer.ID).
					ToSql()

				if err != nil {
					return err
				}

				_, err = tx.Exec(reviewerQuery, reviewerArgs...)
				if err != nil {
					if err.Error() == "pq: duplicate key value violates unique constraint" {
						return models.ErrReviewerAlreadyAssigned
					}
					if strings.Contains(err.Error(), "violates foreign key constraint") {
						return repositories.ErrUserNotFoundInPersistence
					}
					return err
				}
			}
		}
	}

	return tx.Commit()
}

func (p *PullRequestDataBase) FindPossibleReviewers(author *models.User) ([]*models.User, error) {
	reviewersQuery, reviewersArgs, err := p.sb.
		Select("u.id", "u.name", "u.email", "u.team_name", "u.is_active").
		From("team_members m").
		Join("users u ON m.user_id = u.id").
		Where(squirrel.And{
			squirrel.Eq{"u.team_name": author.TeamName},
			squirrel.Eq{"u.is_active": true},
			squirrel.NotEq{"u.id": author.ID},
		}).
		ToSql()

	if err != nil {
		return nil, err
	}

	reviewersRows, err := p.db.Query(reviewersQuery, reviewersArgs...)
	if err != nil {
		return nil, err
	}
	defer reviewersRows.Close()

	var reviewers []*models.User

	for reviewersRows.Next() {
		reviewer := &models.User{}
		err := reviewersRows.Scan(&reviewer.ID, &reviewer.Name, &reviewer.Email, &reviewer.TeamName, &reviewer.IsActive)
		if err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewer)
	}

	if err = reviewersRows.Err(); err != nil {
		return nil, err
	}

	return reviewers, nil
}
