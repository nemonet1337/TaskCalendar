package main

import (
	"log"
	"os"

	"task-calendar-backend/internal/config"
	"task-calendar-backend/internal/database"
	"task-calendar-backend/internal/handlers"
	"task-calendar-backend/internal/middleware"
	"task-calendar-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 設定読み込み
	cfg := config.Load()

	// データベース接続
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	// マイグレーション実行
	if err := database.Migrate(db); err != nil {
		log.Fatal("マイグレーションに失敗しました:", err)
	}

	// サービス初期化
	authService := services.NewAuthService(db, cfg.JWTSecret)
	userService := services.NewUserService(db)
	teamService := services.NewTeamService(db)
	taskService := services.NewTaskService(db)
	eventService := services.NewEventService(db)

	// Cronサービス開始
	cronService := services.NewCronService(eventService)
	cronService.Start()
	defer cronService.Stop()

	// Ginルーター設定
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS設定
	r.Use(middleware.CORS())

	// ハンドラー初期化
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	taskHandler := handlers.NewTaskHandler(taskService)
	eventHandler := handlers.NewEventHandler(eventService)

	// ルート設定
	api := r.Group("/api")
	{
		// 認証不要ルート
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 認証必要ルート
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// ユーザー管理
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
			}

			// チーム管理
			teams := protected.Group("/teams")
			{
				teams.GET("", teamHandler.GetTeams)
				teams.POST("", teamHandler.CreateTeam)
				teams.GET("/:id", teamHandler.GetTeam)
				teams.PUT("/:id", teamHandler.UpdateTeam)
				teams.DELETE("/:id", teamHandler.DeleteTeam)
				teams.POST("/:id/members", teamHandler.AddMember)
				teams.DELETE("/:id/members/:userId", teamHandler.RemoveMember)
			}

			// タスク管理
			tasks := protected.Group("/tasks")
			{
				tasks.GET("", taskHandler.GetTasks)
				tasks.POST("", taskHandler.CreateTask)
				tasks.GET("/:id", taskHandler.GetTask)
				tasks.PUT("/:id", taskHandler.UpdateTask)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
				tasks.POST("/:id/comments", taskHandler.AddComment)
			}

			// イベント管理
			events := protected.Group("/events")
			{
				events.GET("", eventHandler.GetEvents)
				events.POST("", eventHandler.CreateEvent)
				events.GET("/:id", eventHandler.GetEvent)
				events.PUT("/:id", eventHandler.UpdateEvent)
				events.DELETE("/:id", eventHandler.DeleteEvent)
			}
		}
	}

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	log.Printf("🚀 サーバーがポート %s で開始されました", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
