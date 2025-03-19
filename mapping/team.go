package mapping

import (
	"github.com/Sigdriv/paskelabyrint-api/dbstrucs"
	"github.com/Sigdriv/paskelabyrint-api/model"
)

func TeamToDBTeam(team model.Team) dbstrucs.Team {
	return dbstrucs.Team{
		ID:                     team.ID,
		Name:                   team.TeamName,
		Email:                  team.Email,
		CountParticipants:      team.CountParticipants,
		YoungestParticipantAge: team.YoungestParticipantAge,
		OldestParticipantAge:   team.OldestParticipantAge,
		TeamName:               team.TeamName,
		CreatedByID:            team.CreatedByID,
	}
}
