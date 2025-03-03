package user

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mmianl/api/v2/cmd/build"
	api "github.com/mmianl/api/v2/internal/api"
	userV1 "github.com/mmianl/api/v2/internal/api/user/v1"
	"github.com/mmianl/api/v2/internal/config"
	"github.com/mmianl/api/v2/pkg/logging"
	"github.com/mmianl/api/v2/pkg/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configFile string

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Runs the user api.",
		Long: `Runs the user api.
Requirements:
- postgres`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.InitLogger(config.Global().Log)
			log.Debugf("Using log config: %+v", config.Global().Log)

			metrics.Setup("user_api", build.Version, build.Commit)
			shutdown := make(chan os.Signal, 1)
			signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

			httpServer := api.NewHttpServer(config.Global(), shutdown)
			userV1.NewAPI(config.Global(), httpServer.Router)

			return httpServer.Serve()
		},
	}

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "path to user api config file")

	return cmd
}

func initConfig() {
	err := config.Gather(configFile)
	if err != nil {
		log.Fatalf("Unable to gather config from %s: %s", configFile, err.Error())
	}
}
