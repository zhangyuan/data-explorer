package serve

import (
	"data-explorer/pkg/dataexplorer/conf"
	"data-explorer/pkg/dataexplorer/server"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		} else {
			if err := godotenv.Load(); err != nil {
				log.Fatal("Error loading .env file")
			}
		}

		connectionsConf, err := conf.LoadConnection(connectionsPath)
		if err != nil {
			log.Fatal(err)
		}

		server, err := server.NewServer(connectionsConf)
		if err != nil {
			log.Fatal(err)
		}

		if err := server.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

var connectionsPath string

func init() {
	ServeCmd.Flags().StringVarP(&connectionsPath, "connections", "c", "connections.yaml", "Path to the connection conf file")
	_ = ServeCmd.MarkFlagRequired("connections")
}
