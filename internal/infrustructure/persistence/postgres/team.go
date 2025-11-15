package postgres

import (
	"database/sql"
	"errors"
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"strings"

	"github.com/Masterminds/squirrel"
)

type TeamDataBase struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func NewTeamDataBase(db *sql.DB) *TeamDataBase {
	return &TeamDataBase{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (t *TeamDataBase) Add(team *models.Team) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			panic(err)
		}
	}(tx)

	query, args, err := t.sb.
		Insert("teams").
		Columns("name").
		Values(team.Name).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = tx.QueryRow(query, args...).Scan(&team.ID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			return repositories.ErrTeamAlreadyExists
		}
		return err
	}

	if len(team.Members) > 0 {
		for _, member := range team.Members {
			memberQuery, memberArgs, err := t.sb.
				Insert("team_members").
				Columns("team_id", "user_id").
				Values(team.ID, member.UserID).
				ToSql()
			if err != nil {
				return err
			}
			_, err = tx.Exec(memberQuery, memberArgs...)
			if err != nil {
				if err.Error() == "pq: duplicate key value violates unique constraint" {
					return models.ErrMemberAlreadyInTeam
				}
				if strings.Contains(err.Error(), "violates foreign key constraint") {
					return repositories.ErrTeamNotFoundInPersistence
				}
				return err
			}
		}
	}

	return tx.Commit()
}

func (t *TeamDataBase) GetByID(id int) (*models.Team, error) {
	query, args, err := t.sb.
		Select("id", "name").
		From("teams").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}
	team := &models.Team{Members: make(map[int]*models.TeamMember)}
	err = t.db.QueryRow(query, args...).Scan(&team.ID, &team.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrTeamNotFoundInPersistence
		}
		return nil, err
	}
	membersQuery, membersArgs, err := t.sb.
		Select("u.id", "u.name", "u.is_active").
		From("users u").
		Join("teams tm ON u.team_name = tm.name").
		Where(squirrel.Eq{"tm.id": team.ID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	membersRows, err := t.db.Query(membersQuery, membersArgs...)
	if err != nil {
		return nil, err
	}
	defer membersRows.Close()

	for membersRows.Next() {
		var userID int
		var username string
		var isActive bool

		err := membersRows.Scan(&userID, &username, &isActive)
		if err != nil {
			return nil, err
		}

		team.Members[userID] = models.NewTeamMember(userID, username, isActive)
	}

	if err = membersRows.Err(); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *TeamDataBase) GetByName(name string) (*models.Team, error) {
	query, args, err := t.sb.
		Select("id", "name").
		From("teams").
		Where(squirrel.Eq{"name": name}).
		ToSql()

	if err != nil {
		return nil, err
	}
	team := &models.Team{Members: make(map[int]*models.TeamMember)}
	err = t.db.QueryRow(query, args...).Scan(&team.ID, &team.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrTeamNotFoundInPersistence
		}
		return nil, err
	}
	membersQuery, membersArgs, err := t.sb.
		Select("u.id", "u.name", "u.is_active").
		From("users u").
		Join("teams tm ON u.team_name = tm.name").
		Where(squirrel.Eq{"tm.name": name}).
		ToSql()

	if err != nil {
		return nil, err
	}

	membersRows, err := t.db.Query(membersQuery, membersArgs...)
	if err != nil {
		return nil, err
	}
	defer membersRows.Close()

	for membersRows.Next() {
		var userID int
		var username string
		var isActive bool

		err := membersRows.Scan(&userID, &username, &isActive)
		if err != nil {
			return nil, err
		}

		team.Members[userID] = models.NewTeamMember(userID, username, isActive)
	}

	if err = membersRows.Err(); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *TeamDataBase) GetAll() ([]*models.Team, error) {
	teamsQuery, teamsArgs, err := t.sb.
		Select("id", "name").
		From("teams").
		ToSql()

	if err != nil {
		return nil, err
	}

	teamsRows, err := t.db.Query(teamsQuery, teamsArgs...)
	if err != nil {
		return nil, err
	}
	defer teamsRows.Close()

	var teams []*models.Team
	teamsByID := make(map[int]*models.Team)

	for teamsRows.Next() {
		team := &models.Team{
			Members: make(map[int]*models.TeamMember),
		}
		err := teamsRows.Scan(&team.ID, &team.Name)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
		teamsByID[team.ID] = team
	}

	if err = teamsRows.Err(); err != nil {
		return nil, err
	}

	if len(teams) == 0 {
		return teams, nil
	}

	teamIDs := make([]int, 0, len(teams))
	for _, team := range teams {
		teamIDs = append(teamIDs, team.ID)
	}

	membersQuery, membersArgs, err := t.sb.
		Select("tm.team_id", "u.id", "u.name", "u.is_active").
		From("team_members tm").
		Join("users u ON tm.user_id = u.id").
		Where(squirrel.Eq{"tm.team_id": teamIDs}).
		ToSql()

	if err != nil {
		return nil, err
	}

	membersRows, err := t.db.Query(membersQuery, membersArgs...)
	if err != nil {
		return nil, err
	}
	defer membersRows.Close()

	for membersRows.Next() {
		var teamID int
		var userID int
		var username string
		var isActive bool

		err := membersRows.Scan(&teamID, &userID, &username, &isActive)
		if err != nil {
			return nil, err
		}

		if team, exists := teamsByID[teamID]; exists {
			team.Members[userID] = models.NewTeamMember(userID, username, isActive)
		}
	}

	if err = membersRows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (t *TeamDataBase) Update(team *models.Team) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, args, err := t.sb.
		Update("teams").
		Set("name", team.Name).
		Where(squirrel.Eq{"id": team.ID}).
		ToSql()

	if err != nil {
		return err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			return repositories.ErrTeamAlreadyExists
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repositories.ErrTeamNotFoundInPersistence
	}

	if team.Members != nil {
		deleteQuery, deleteArgs, err := t.sb.
			Delete("team_members").
			Where(squirrel.Eq{"team_id": team.ID}).
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.Exec(deleteQuery, deleteArgs...)
		if err != nil {
			return err
		}

		if len(team.Members) > 0 {
			for _, member := range team.Members {
				memberQuery, memberArgs, err := t.sb.
					Insert("team_members").
					Columns("team_id", "user_id").
					Values(team.ID, member.UserID).
					ToSql()
				if err != nil {
					return err
				}

				_, err = tx.Exec(memberQuery, memberArgs...)
				if err != nil {
					if err.Error() == "pq: duplicate key value violates unique constraint" {
						return models.ErrMemberAlreadyInTeam
					}
					if strings.Contains(err.Error(), "violates foreign key constraint") {
						return repositories.ErrTeamNotFoundInPersistence
					}
					return err
				}
			}
		}
	}

	return tx.Commit()
}

func (t *TeamDataBase) AddUserToTeam(teamID, userID int) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	teamExistsQuery, teamExistsArgs, err := t.sb.
		Select("1").
		From("teams").
		Where(squirrel.Eq{"id": teamID}).
		ToSql()

	if err != nil {
		return err
	}

	var teamExists int
	err = tx.QueryRow(teamExistsQuery, teamExistsArgs...).Scan(&teamExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repositories.ErrTeamNotFoundInPersistence
		}
		return err
	}

	userExistsQuery, userExistsArgs, err := t.sb.
		Select("1").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		return err
	}

	var userExists int
	err = tx.QueryRow(userExistsQuery, userExistsArgs...).Scan(&userExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repositories.ErrUserNotFoundInPersistence
		}
		return err
	}

	insertQuery, insertArgs, err := t.sb.
		Insert("team_members").
		Columns("team_id", "user_id").
		Values(teamID, userID).
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(insertQuery, insertArgs...)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			return models.ErrMemberAlreadyInTeam
		}
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return repositories.ErrTeamNotFoundInPersistence
		}
		return err
	}

	return tx.Commit()
}

func (t *TeamDataBase) RemoveUserFromTeam(teamID, userID int) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	membershipQuery, membershipArgs, err := t.sb.
		Select("1").
		From("team_members").
		Where(squirrel.And{
			squirrel.Eq{"team_id": teamID},
			squirrel.Eq{"user_id": userID},
		}).
		ToSql()

	if err != nil {
		return err
	}

	var membershipExists int
	err = tx.QueryRow(membershipQuery, membershipArgs...).Scan(&membershipExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrMemberNotInTeam
		}
		return err
	}

	deleteQuery, deleteArgs, err := t.sb.
		Delete("team_members").
		Where(squirrel.And{
			squirrel.Eq{"team_id": teamID},
			squirrel.Eq{"user_id": userID},
		}).
		ToSql()

	if err != nil {
		return err
	}

	result, err := tx.Exec(deleteQuery, deleteArgs...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrMemberNotInTeam
	}

	return tx.Commit()
}
