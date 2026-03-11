package english

import "log/slog"

func TestLogs() {
	// correct messages
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Info("token validated")

	// incorrect english
	slog.Info("запуск сервера")      // want "log message should be in english"
	slog.Error("ошибка подключения") // want "log message should be in english"
}
