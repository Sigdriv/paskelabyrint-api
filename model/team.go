package model

type Team struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name" binding:"required"`
	Email                  string `json:"email" binding:"required"`
	CountParticipants      int    `json:"countParticipants" binding:"required"`
	YoungestParticipantAge *int   `json:"youngestParticipantAge"`
	OldestParticipantAge   *int   `json:"oldestParticipantAge"`
	TeamName               string `json:"teamName" binding:"required"`
	CreatedAt              string `json:"createdAt"`
	UpdatedAt              string `json:"updatedAt"`
	CreatedByID            string `json:"createdByID"`
}
