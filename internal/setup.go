package internal

import (
	"context"
	"log"

	"rainbowcoloringbooks/internal/db"
	userRepo "rainbowcoloringbooks/internal/repository/user"
	userService "rainbowcoloringbooks/internal/service/user"

	_ "github.com/lib/pq"
)

func SetupApp(configFileName string) (context.Context, *db.PostgresDatabase, userService.UserService) {
    ctx := context.Background()
    config, err := LoadConfig(configFileName)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    postgresDB := &db.PostgresDatabase{
        User:     config.DBUser,
        Password: config.DBPassword,
        DBName:   config.DBName,
        SSLMode:  "disable",
    }

    err = postgresDB.Connect(ctx)

    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    userRepo := userRepo.NewPostgresUserRepository(postgresDB.DB)
    userService := userService.NewUserService(userRepo)

    return ctx, postgresDB, userService
}