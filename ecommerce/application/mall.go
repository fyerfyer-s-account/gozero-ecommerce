package main

import (
	"flag"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/handler"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/middleware"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/mall-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	server.Use(middleware.Cors)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
