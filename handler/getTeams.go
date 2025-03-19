package handler

import (
	"net/http"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/dbstrucs"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func getTeams(c *gin.Context) {
	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Failed to connect to database >>", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database connection error"})
		return
	}

	defer conn.Close()

	rows, err := conn.Query(c, `select id, name, email, count_participants, youngest_participant_age, oldest_participant_age, team_name, created_by_id from teams`)
	if err != nil {
		log.Error("Failed to query teams >>", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database query error"})
		return
	}
	defer rows.Close()

	var teams []dbstrucs.Team
	for rows.Next() {
		var team dbstrucs.Team
		err := rows.Scan(&team.ID, &team.Name, &team.Email, &team.CountParticipants, &team.YoungestParticipantAge, &team.OldestParticipantAge, &team.TeamName, &team.CreatedByID)
		if err != nil {
			log.Error("Failed to scan team >>", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database scan error"})
			return
		}
		teams = append(teams, team)
	}
	if err := rows.Err(); err != nil {
		log.Error("Failed to iterate over teams >>", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database iteration error"})
		return
	}

	if len(teams) == 0 {
		log.Warn("No teams found")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No teams found"})
		return
	}

	log.Info("Successfully retrieved teams")
	c.IndentedJSON(http.StatusOK, teams)
}
