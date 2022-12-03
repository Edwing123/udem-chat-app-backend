package sqlserver

import (
	"context"
	"database/sql"

	"github.com/Edwing123/udem-chat-app/pkg/models"
	"golang.org/x/exp/slog"
)

var (
	rootCtx = context.Background()
)

func New(db *sql.DB, logger *slog.Logger) models.Database {
	userManager := &UserManager{
		db:     db,
		logger: logger,
	}

	return models.Database{
		UserManager: userManager,
	}
}
