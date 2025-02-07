package account

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func AccountDeleteHandler(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	if err := nk.AccountDeleteId(ctx, userID, true); err != nil {
		logger.Error("Unable to delete account: %v", err)
		return "error", err
	}
	return "success", nil
}
