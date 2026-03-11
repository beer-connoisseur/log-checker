package nosensitive

import "log/slog"

type User struct {
	Password string
}

func TestLogs() {
	// correct messages
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Info("token validated")
	slog.Info("token")
	slog.Info("token: " + "abc123")

	safe := "safe"
	slog.Info("value: " + safe)

	// incorrect nosensitive
	password := "secret"
	slog.Info("user password: " + password) // want "log message may contain sensitive data"

	apiKey := "xyz"
	slog.Info("apikey=" + apiKey) // want "log message may contain sensitive data"

	user := User{Password: "secret"}
	slog.Info("user password: " + user.Password) // want "log message may contain sensitive data"

	ptr := &user
	slog.Info("user password: " + ptr.Password)    // want "log message may contain sensitive data"
	slog.Info("user password: " + (*ptr).Password) // want "log message may contain sensitive data"

	key := "mykey"
	slog.Info("key=" + key) // want "log message may contain sensitive data"
}
