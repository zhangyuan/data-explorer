package serve

import (
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

		server, err := server.NewServer()
		if err != nil {
			log.Fatal(err)
		}

		if err := server.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
}
