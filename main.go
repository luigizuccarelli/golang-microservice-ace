package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	logger  simple.Logger
	config  Config
	redisdb *redis.Client
)

func startHttpServer(cfg Config) *http.Server {

	config = cfg

	logger.Debug(fmt.Sprintf("Config in startServer %v ", config))
	srv := &http.Server{Addr: ":" + config.Port}

	r := mux.NewRouter()
	r.HandleFunc("/sys/info/isalive", IsAlive).Methods("GET")
	r.HandleFunc("/login", MiddlewareLogin).Methods("POST")
	r.HandleFunc("/alldata", MiddlewareData).Methods("POST")
	http.Handle("/", r)

	// connect to redis
	redisdb = redis.NewClient(&redis.Options{
		Addr:         config.RedisDB.Host + ":" + config.RedisDB.Port,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		Password:     "",
		DB:           0,
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Httpserver: ListenAndServe() error: " + err.Error())
		}
	}()

	return srv
}

func main() {
	// read the config
	config, _ := Init("config.json")
	logger.Level = config.Level
	srv := startHttpServer(config)
	logger.Info("Starting server on port " + config.Port)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan

	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
	logger.Info("Server shutdown successfully")
	os.Exit(code)
}
