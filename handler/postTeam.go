package handler

import (
	"net/http"

	"github.com/Sigdriv/paskelabyrint-api/db"

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

	user, err := getUser(c)
	if err != nil {
		log.Error("Failed to get user >> ", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Failed to connect to database >>", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database connection error"})
		return
	}
	defer conn.Close()

	rows, err := conn.Query(c, `select team_name from teams`)
	if err != nil {
		log.Error("Failed to query database >>", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database query error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var teamName string
		err := rows.Scan(&teamName)
		if err != nil {
			log.Error("Failed to scan row >>", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database scan error"})
			return
		}

		if teamName == newTeam.TeamName {
			log.Error("Team name already exists >> ", newTeam.TeamName)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Team name already exists"})
			return
		}
	}

	_, err = conn.Exec(c, `INSERT INTO teams`+
		`(name, email, count_participants, youngest_participant_age, oldest_participant_age, team_name, created_by_id)`+
		`VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		newTeam.Name,
		newTeam.Email,
		newTeam.CountParticipants,
		newTeam.YoungestParticipantAge,
		newTeam.OldestParticipantAge,
		newTeam.TeamName,
		user.ID,
	)
	if err != nil {
		log.Error("Failed to insert team into database >>", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database insert error"})
		return
	}

	log.Infof("New team added: %+v", newTeam)
	c.IndentedJSON(http.StatusCreated, newTeam)
}
