package handler

import (
	"net/http"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/mapping"
	"github.com/Sigdriv/paskelabyrint-api/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func postTeam(c *gin.Context) {
	var newTeam model.Team

	err := c.Bind(&newTeam)
	if err != nil {
		log.Error("Failed to bind JSON >>", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	dbTeam := mapping.TeamToDBTeam(newTeam)

	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Failed to connect to database >>", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database connection error"})
		return
	}
	defer conn.Close()

	_, err = conn.Exec(c, `INSERT INTO teams`+
		`(name, email, count_participants, youngest_participant_age, oldest_participant_age, team_name, created_by_id)`+
		`VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		dbTeam.Name,
		dbTeam.Email,
		dbTeam.CountParticipants,
		dbTeam.YoungestParticipantAge,
		dbTeam.OldestParticipantAge,
		dbTeam.TeamName,
		dbTeam.CreatedByID,
	)
	if err != nil {
		log.Error("Failed to insert team into database >>", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database insert error"})
		return
	}

	log.Infof("New team added: %+v", newTeam)
	c.IndentedJSON(http.StatusCreated, newTeam)
}
