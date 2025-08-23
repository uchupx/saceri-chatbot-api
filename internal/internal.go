// Package internal provides a factory for initializing and managing application components.
package internal

import (
	"github.com/uchupx/saceri-chatbot-api/internal/api/handlers"
	"github.com/uchupx/saceri-chatbot-api/internal/api/middlewares"
	"github.com/uchupx/saceri-chatbot-api/internal/config"
	"github.com/uchupx/saceri-chatbot-api/internal/database"
	"github.com/uchupx/saceri-chatbot-api/internal/repository"
	"github.com/uchupx/saceri-chatbot-api/internal/repository/mongodb"
	"github.com/uchupx/saceri-chatbot-api/internal/repository/redis"
	"github.com/uchupx/saceri-chatbot-api/internal/service"
	"github.com/uchupx/saceri-chatbot-api/pkg/apilog"
	"github.com/uchupx/saceri-chatbot-api/pkg/grpc/client"
)

type Factory struct {
	userHandler    *handlers.UserHandler
	authHandler    *handlers.AuthHandler
	settingHandler *handlers.SettingHandler
	handler        *handlers.Handler

	authService *client.AuthClient

	middleware *middlewares.Middleware

	dbConn  *database.MongoDB
	dbRedis *database.Cache

	userRepo    repository.UserRepoInterface
	settingRepo repository.SettingRepoInterface

	cacheRepo        *redis.CacheRepo
	settingRepoCache repository.SettingRepoInterface

	log *apilog.ApiLog

	userService    *service.UserService
	settingService *service.SettingService
}

func (f *Factory) GetUserHandler() *handlers.UserHandler {
	if f.userHandler != nil {
		return f.userHandler
	}

	f.userHandler = &handlers.UserHandler{
		Handler:     *f.Handler(),
		UserService: f.GetUserService(),
	}

	return f.userHandler
}

func (f *Factory) Handler() *handlers.Handler {

	if f.handler != nil {
		return f.handler
	}

	f.handler = handlers.NewHandler(f.GetLog())

	return f.handler
}

func (f *Factory) GetAuthHandler() *handlers.AuthHandler {
	if f.authHandler != nil {
		return f.authHandler
	}

	f.authHandler = &handlers.AuthHandler{
		Handler:     *f.Handler(),
		AuthClient:  f.AuthClient(),
		UserService: f.GetUserService(),
	}

	return f.authHandler
}

func (f *Factory) GetSettingHandler() *handlers.SettingHandler {
	if f.settingHandler != nil {
		return f.settingHandler
	}

	f.settingHandler = &handlers.SettingHandler{
		Handler:        *f.Handler(),
		SettingService: f.GetSettingService(),
	}

	return f.settingHandler
}
func (f *Factory) GetDBConnection() *database.MongoDB {
	if f.dbConn != nil {
		return f.dbConn
	}

	conf := config.GetConfig()

	var err error

	f.dbConn, err = database.NewMongoDB(conf.Database.URL)
	if err != nil {
		panic(err)
	}

	return f.dbConn
}

func (f *Factory) AuthClient() *client.AuthClient {
	if f.authService != nil {
		return f.authService
	}

	auth, err := client.NewAuthClient(config.GetConfig().Service.AuthServiceUrl)
	if err != nil {
		panic(err)
	}

	f.authService = auth

	return f.authService
}

func (f *Factory) GetMiddleware() *middlewares.Middleware {
	if f.middleware != nil {
		return f.middleware
	}

	middleware := middlewares.Middleware{
		AuthClient: f.AuthClient(),
	}

	f.middleware = &middleware

	return f.middleware
}

func (f *Factory) GetUserRepo() repository.UserRepoInterface {
	if f.userRepo != nil {
		return f.userRepo
	}

	db := f.GetDBConnection()
	userRepo := mongodb.NewUserRepoMongodb(db.Client)

	f.userRepo = userRepo

	return f.userRepo
}

func (f *Factory) GetSettingRepo() repository.SettingRepoInterface {
	if f.settingRepo != nil {
		return f.settingRepo
	}

	db := f.GetDBConnection()
	settingRepo := mongodb.NewSettingRepoMongodb(db.Client)

	f.settingRepo = settingRepo

	return f.settingRepo
}

func (f *Factory) GetUserService() *service.UserService {
	if f.userService != nil {
		return f.userService
	}

	userService := &service.UserService{
		UserRepo: f.GetUserRepo(),
	}

	f.userService = userService

	return f.userService
}

func (f *Factory) GetSettingService() *service.SettingService {
	if f.settingService != nil {
		return f.settingService
	}

	settingService := &service.SettingService{
		SettingRepo: f.GetSettingRepo(),
	}

	f.settingService = settingService

	return f.settingService
}

func (f *Factory) GetSettingCacheRepo() repository.SettingRepoInterface {
	if f.settingRepoCache != nil {
		return f.settingRepoCache
	}

	settingCacheRepo := redis.NewSettingCacheRepo(
		f.GetCacheRepo(),
		f.GetSettingRepo(),
	)

	f.settingRepoCache = settingCacheRepo

	return f.settingRepoCache
}

func (f *Factory) GetCacheRepo() *redis.CacheRepo {
	if f.cacheRepo != nil {
		return f.cacheRepo
	}

	cacheRepo := redis.NewCacheRepo(
		f.GetCache(),
		f.GetLog(),
	)

	f.cacheRepo = cacheRepo

	return f.cacheRepo
}

func (f *Factory) GetCache() *database.Cache {
	if f.dbRedis != nil {
		return f.dbRedis
	}

	conf := config.GetConfig()

	cache, err := database.GetConnection(database.RedisConfig{
		Host:         conf.Redis.Host,
		Password:     conf.Redis.Password,
		Database:     conf.Redis.Database,
		PoolSize:     conf.Redis.PoolSize,
		MinIdleConns: conf.Redis.MinIdleConns,
		// IdleTimeout: conf.Redis.IdleTimeout,

	})
	if err != nil {
		panic(err)
	}

	f.dbRedis = cache

	return f.dbRedis
}

func (f *Factory) GetLog() *apilog.ApiLog {
	if f.log != nil {
		return f.log
	}

	conf := config.GetConfig()

	log := apilog.NewApiLog(apilog.Params{
		ServiceName: conf.App.Name,
		Level:       4,
		Version:     conf.App.Version,
	})

	f.log = log

	return f.log
}
