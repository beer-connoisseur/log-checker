package nospecials

import "log/slog"

func TestLogs() {
	// correct messages
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Info("token validated")

	// incorrect nospecials
	slog.Info("server started!")                // want "log message should not contain special characters or emojis"
	slog.Info("server started! 🚀")              // want "log message should not contain special characters or emojis"
	slog.Error("connection failed!!!")          // want "log message should not contain special characters or emojis"
	slog.Warn("warning: something went wrong#") // want "log message should not contain special characters or emojis"
	slog.Info("hello... world!")                // want "log message should not contain special characters or emojis"
	slog.Info("please wait--now?")              // want "log message should not contain special characters or emojis"
}
