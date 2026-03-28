/*
Setup logging for us
*/

package main 

import (
	"log/slog"
	"os"
)

func SetupLoggingToFile() (*slog.Logger, func()) {
	fname, err := logsFilePath()
	if err != nil {
		panic("Can't set up logging: " + err.Error())
	}
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Error opening file:"  + err.Error())
	}

	closeFunc := func() { file.Close() }

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)

	return logger, closeFunc
}
