package main

import (
	"StoryGoAPI/cmd"
	"StoryGoAPI/config"
	"StoryGoAPI/docs"
	_ "StoryGoAPI/docs"
	"flag"
)

var (
	cfg     *config.Config
	cfgPath string
)

func init() {
	flag.StringVar(&cfgPath, "c", "config.yml", "path to configuration file (shorthand)")
	flag.StringVar(&cfgPath, "config", "config.yml", "path to configuration file")
	flag.Parse()
	var err error
	cfg, err = config.ReadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	docs.SwaggerInfo.Host = cfg.Server.Docs.Host
	docs.SwaggerInfo.BasePath = cfg.Server.Docs.BasePath
	docs.SwaggerInfo.Version = cfg.Server.Docs.Version
	docs.SwaggerInfo.Title = cfg.Server.Docs.Title
	docs.SwaggerInfo.Description = cfg.Server.Docs.Description

	if err := cmd.Run(cfg); err != nil {
		panic(err)
	}
}
