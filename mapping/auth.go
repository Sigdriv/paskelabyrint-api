package mapping

import (
	"github.com/Sigdriv/paskelabyrint-api/dbstrucs"
	"github.com/Sigdriv/paskelabyrint-api/model"
)

func SessionToSBSession(session model.Session) dbstrucs.Session {
	return dbstrucs.Session{
		ID:        session.ID,
		Email:     session.Email,
		CreatedAt: session.CreatedAt,
		Expiry:    session.Expiry,
	}
}
