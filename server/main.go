package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zhuiye8/Lyss/server/api/agent"
	"github.com/zhuiye8/Lyss/server/api/application"
	"github.com/zhuiye8/Lyss/server/api/auth"
	"github.com/zhuiye8/Lyss/server/api/config"
	"github.com/zhuiye8/Lyss/server/api/conversation"
	"github.com/zhuiye8/Lyss/server/api/model"
	"github.com/zhuiye8/Lyss/server/api/project"
	"github.com/zhuiye8/Lyss/server/models"
	authPkg "github.com/zhuiye8/Lyss/server/pkg/auth"
	"github.com/zhuiye8/Lyss/server/pkg/encryption"
	"github.com/zhuiye8/Lyss/server/pkg/middleware"
	"github.com/zhuiye8/Lyss/server/api/dashboard"
)

func main() {
	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		zap.L().Info("No .env file found")
	}

	// 加载配置
	loadConfig()

	// 设置Gin模式
	if viper.GetString("app.env") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	db, err := initDatabase()
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}

	// 自动迁移模型（仅在开发环境中使用）
	if viper.GetString("app.env") == "development" {
		tx := db.Begin()

		// 禁用外键约束检查（PostgreSQL）
		tx.Exec("SET session_replication_role = 'replica';")

		// 进行迁移
		if err := tx.AutoMigrate(
			&models.User{},
			&models.Config{},
			&models.ModelConfig{},
			&models.Model{},
			&models.KnowledgeBase{},
			&models.Project{},
			&models.Application{},
			&models.Agent{},
			&models.AgentKnowledgeBase{},
			&models.Conversation{},
			&models.Message{},
			&models.Log{},
			&models.SystemMetric{},
		); err != nil {
			tx.Rollback()
			zap.L().Fatal("Failed to migrate database", zap.Error(err))
		}

		// 重新启用外键约束检查
		tx.Exec("SET session_replication_role = 'origin';")

		// 提交事务
		tx.Commit()
	}

	// 初始化JWT管理器
	jwtManager := authPkg.NewJWTManager(authPkg.Config{
		SecretKey:     viper.GetString("jwt.secret"),
		TokenExpiry:   viper.GetDuration("jwt.expiry"),
		RefreshExpiry: viper.GetDuration("jwt.refresh_expiry"),
		Issuer:        viper.GetString("app.name"),
	})

	// 初始化加密服务
	encryptionService, err := encryption.NewService(viper.GetString("encryption.secret"))
	if err != nil {
		zap.L().Fatal("Failed to initialize encryption service", zap.Error(err))
	}

	// 初始化认证中间件
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// 初始化认证服务
	authService := auth.NewService(db, jwtManager)
	authHandler := auth.NewHandler(authService, authMiddleware)

	// 初始化项目服务
	projectService := project.NewService(db)
	projectHandler := project.NewHandler(projectService, authMiddleware)

	// 初始化应用服务
	applicationService := application.NewService(db)
	applicationHandler := application.NewHandler(applicationService, authMiddleware)

	// 初始化配置服务
	configService := config.NewService(db)
	configHandler := config.NewHandler(configService, authMiddleware)

	// 初始化模型服务
	modelService := model.NewService(db, encryptionService)
	modelHandler := model.NewHandler(modelService, authMiddleware)

	// 初始化智能体服务
	agentService := agent.NewService(db)
	agentHandler := agent.NewHandler(agentService, authMiddleware)

	// 初始化对话服务
	conversationService := conversation.NewService(db)
	conversationHandler := conversation.NewHandler(conversationService, authMiddleware)

	// 初始化仪表盘服务
	dashboardService := dashboard.NewService(db)
	dashboardHandler := dashboard.NewHandler(dashboardService, authMiddleware)

	// 创建路由
	r := gin.Default()

	// 添加请求ID中间件
	r.Use(middleware.RequestID())
	
	// 添加错误处理中间件
	r.Use(middleware.ErrorHandler(logger))
	
	// 允许跨域
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		// 检查数据库连接
		if err := db.Raw("SELECT 1").Error; err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"message": "数据库连接失败",
				"error": err.Error(),
				"time": time.Now().Format(time.RFC3339),
			})
			return
		}

		// 返回服务信息
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"version": viper.GetString("app.version"),
			"environment": viper.GetString("app.env"),
			"time": time.Now().Format(time.RFC3339),
		})
	})

	// API路由
	api := r.Group("/api/v1")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "智能体构建平台API",
				"version": viper.GetString("app.version"),
			})
		})

		// 注册各种处理器路由
		authHandler.RegisterRoutes(api)
		projectHandler.RegisterRoutes(api)
		applicationHandler.RegisterRoutes(api)
		configHandler.RegisterRoutes(api)
		modelHandler.RegisterRoutes(api)
		
		// 注册新增的处理器路由
		agentHandler.RegisterRoutes(api)
		conversationHandler.RegisterRoutes(api)
		dashboardHandler.RegisterRoutes(api)
	}

	// 启动服务
	port := viper.GetString("app.port")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	zap.L().Info(fmt.Sprintf("Server started on port %s", port))

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server exited")
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 默认配置
	viper.SetDefault("app.name", "agent-platform")
	viper.SetDefault("app.version", "0.1.0")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("jwt.expiry", "24h")
	viper.SetDefault("jwt.refresh_expiry", "168h") // 7天
	viper.SetDefault("jwt.secret", "your-secret-key-change-me")
	viper.SetDefault("encryption.secret", "your-encryption-key-must-be-32-chars")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			zap.L().Info("No config file found, using defaults")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}
}

func initDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		viper.GetString("database.host"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
		viper.GetString("database.port"),
	)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	config := &gorm.Config{
		Logger: gormLogger,
	}

	return gorm.Open(postgres.Open(dsn), config)
} 
