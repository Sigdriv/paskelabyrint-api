package dbstrucs

import "time"

type Team struct {
	ID                     int       `db:"id"`
	Name                   string    `db:"name"`
	Email                  string    `db:"email"`
	CountParticipants      int       `db:"count_participants"`
	YoungestParticipantAge *int      `db:"youngest_participant_age"`
	OldestParticipantAge   *int      `db:"oldest_participant_age"`
	TeamName               string    `db:"team_name"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
	CreatedByID            string    `db:"created_by_id"`
}
