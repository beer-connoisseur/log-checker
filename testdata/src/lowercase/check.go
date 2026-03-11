package lowercase

import "log/slog"

func TestLogs() {
	// correct messages
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Info("token validated")

	// incorrect lowercase
	slog.Info("Starting server")        // want "log message should start with a lowercase letter"
	slog.Error("Failed to connect")     // want "log message should start with a lowercase letter"
	slog.Error("  \nFailed to connect") // want "log message should start with a lowercase letter"
}
