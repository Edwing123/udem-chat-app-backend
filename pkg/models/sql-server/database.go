package sqlserver

import (
	"context"
	"database/sql"

	"github.com/Edwing123/udem-chat-app/pkg/models"
)

var (
	rootCtx = context.Background()
)

func New(db *sql.DB) models.Database {
	userManager := &UserManager{
		db: db,
	}

	return models.Database{
		UserManager: userManager,
	}
}
