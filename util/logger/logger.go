package logger

import (
	"log/slog"
	"os"

	"github.com/MatusOllah/slogcolor"
)

func Init() {
	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions)))
}
