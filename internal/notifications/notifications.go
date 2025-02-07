package notification

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func SendNotificationHandler(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	var raw json.RawMessage
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		logger.Error("Failed first unmarshal: %v", err)
		return "", err
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(raw, &payloadMap); err != nil {
		logger.Error("Failed second unmarshal: %v", err)
		return "", err
	}

	subject := payloadMap["subject"].(string)
	content := map[string]interface{}{
		"message": subject,
	}

	code := 1
	senderID := ""
	persistent := true

	if err := nk.NotificationSend(ctx, userID, subject, content, code, senderID, persistent); err != nil {
		logger.Error("Failed to send notification: %v", err)
		return "", err
	}

	return "success", nil
}
