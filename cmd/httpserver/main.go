package main

import (
	"flag"
	"net/http"

	"github.com/eXoterr/FLProject/internal/config"
	"github.com/eXoterr/FLProject/internal/logger"
	"github.com/eXoterr/FLProject/internal/routing"
	"github.com/eXoterr/FLProject/internal/routing/middlewares"
	"github.com/eXoterr/FLProject/internal/storage/store"
	"github.com/go-chi/chi/v5"
)

var (
	configPath string
	doSync     bool
)

// Задание флага для указания пути файла конфигурации
func init() {
	flag.StringVar(&configPath, "config", "config/config.yaml", "config.yaml path")
	flag.BoolVar(&doSync, "sync", false, "syncs db struct with defined models")
}

func main() {
	// Считывание переданных флагов и их значений
	flag.Parse()

	// Считывание файла конфигурации по переданному или стандартному пути
	conf := config.MustLoad(configPath)

	// Настройка модуля журналирования
	log := logger.MustSetupLogger(conf.Logger.LogLevel, conf.Logger.Format)
	log.Debug("logger is ready")

	// Создание подключения к базе данных
	store := store.MustSetup(conf.Database.PostgreSQL, doSync)

	// Настройка маршрутизатора запросов
	router := chi.NewRouter()
	reqLogger := logger.RequestLogger(log)       // "Request Logger" middleware for router
	cors := middlewares.SetupCORS(conf.API.CORS) // CORS middleware

	// Подключение к маршрутизатору основных и промежуточных
	// обработчиков
	routing.SetupMiddleware(router, reqLogger, cors)
	routing.SetupHandlers(router, store, log, conf)

	// Настройка HTTP сервера через полученную конфигурацию
	srv := &http.Server{
		Addr:         conf.API.ListenAddr,
		Handler:      router,
		ReadTimeout:  conf.API.Timeout,
		WriteTimeout: conf.API.Timeout,
		IdleTimeout:  conf.API.IdleTimeout,
	}

	//Запуск HTTP сервера
	err := srv.ListenAndServe()
	if err != nil {
		log.Error("failed to start server")
	}
}
