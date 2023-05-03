package main

import (
	"log"

	dotenv "github.com/joho/godotenv"
	"github.com/sudo-nick16/env"
	"github.com/sudo-nick16/smark/galactus/types"
)

func setupConfig() *types.Config {
	err := dotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &types.Config{
		Port:       env.GetEnv("PORT", ":42069"),
		DbUrl:      env.GetEnv("DB_URL", "mongodb://root:shorty1@127.0.0.1:27017/?serverSelectionTimeoutMS=2000"),
		AccessKey:  env.GetEnv("ACCESS_KEY", "neioneio"),
		RefreshKey: env.GetEnv("REFRESH_KEY", "arstarst"),
		GoogleConfig: types.GoogleConfig{
			ClientId:     env.GetEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: env.GetEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectUrl:  env.GetEnv("GOOGLE_REDIRECT_URI", ""),
		},
	}
}