package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gongfu/cmd/wire"
	"gongfu/internal/config"
	"gongfu/internal/http/middlewares"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config/production.toml", "use config file")
	flag.Parse()

	logrus.Printf("use config file: %s\n", configFile)
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("read config failed: " + err.Error())
	}
	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		panic("unmarshal config failed: " + err.Error())
	}
	initEnv(conf.Env)
	if err := initLogger(&conf.Logger); err != nil {
		panic(err)
	}

	route, err := wire.NewApp(&conf)
	if err != nil {
		logrus.WithError(err).Panic("wire failed")
	}

	r := gin.New()
	r.Use(middlewares.JSONLogMiddleware(), gin.Recovery())

	route.Route(r)
	if err := r.Run("0.0.0.0:9003"); err != nil {
		logrus.WithError(err).Panic("run server failed")
	}
}

func initEnv(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if env == "test" {
		gin.SetMode(gin.TestMode)
	}
}

func initLogger(conf *config.Logger) error {
	lv, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lv)
	if conf.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	return nil
}
