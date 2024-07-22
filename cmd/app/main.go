package main

import(
	"log/slog"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	
}

func newLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(
		os.Stdout, 
		&slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}

func initEnv() {
	err := godotenv.Load("app.env")
	if err != nil {
		panic(err)
	}
}