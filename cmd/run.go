package cmd

import (
	"github.com/gitlayzer/gee_loadbalancer/cmd/config"
	"github.com/gitlayzer/gee_loadbalancer/cmd/middleware"
	"github.com/gitlayzer/gee_loadbalancer/cmd/proxy"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strconv"
)

var cfg string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the load balancer",
	Long:  `run the load balancer and start listening on the specified port`,

	Run: RunCmd,
}

func RunCmd(cmd *cobra.Command, args []string) {
	c, err := config.NewReadConfig(cfg)
	if err != nil {
		log.Fatalf("read config file error: %v", err)
	}

	err = c.Validation()
	if err != nil {
		log.Fatalf("config file validation error: %v", err)
	}

	// router 创建一个新的路由
	router := mux.NewRouter()
	for _, location := range c.Locations {
		httpProxy, err := proxy.NewHTTPProxy(location.Servers, location.BalanceMode)
		if err != nil {
			log.Fatalf("create http proxy error: %v", err)
		}

		if c.TCPHealthCheck {
			httpProxy.HealthCheck(c.TCPHealthCheckInterval)
		}

		router.Handle(location.Prefix, httpProxy)
	}

	if c.MaxAllowed > 0 {
		router.Use(middleware.MaxAllowedMiddleware(c.MaxAllowed))
	}

	svr := http.Server{
		Addr:    ":" + strconv.Itoa(c.ListenPort),
		Handler: router,
	}

	c.Print()

	if c.ListenSchema == "http" {
		err = svr.ListenAndServe()
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	} else if c.ListenSchema == "https" {
		err = svr.ListenAndServeTLS(c.SSLCertificate, c.SSLCertificateKey)
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	}
}

func init() {
	runCmd.Flags().StringVarP(&cfg, "config", "c", "config.yaml", "config file path")
	rootCmd.AddCommand(runCmd)
}
