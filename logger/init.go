package sl

import (
	"io"
	"log"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func ExecLog(path string) *slog.Logger {

	var (
		level  slog.Level
		writer io.Writer
	)

	if path == "" {
		path = getConfig()
	}

	level = slog.LevelInfo
	writer = io.MultiWriter(
		os.Stdout,
		createFileWriter(path),
	)

	return slog.New(
		slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level}),
	)
}

func createFileWriter(path string) io.Writer {
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}

	return &lumberjack.Logger{
		Filename:   path + "/info.log",
		MaxSize:    25,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}
}

func getConfig() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
