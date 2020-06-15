package cmd

import (
	"context"
	"fmt"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"tf2bdd/core"
	"tf2bdd/steamid"
	"time"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tf2bdd",
	Short: "Backend services for tf2_bot_detector",
	Long:  `tf2bdd provides HTTP and discord services for use with tf2_bot_detector`,
	Run: func(cmd *cobra.Command, args []string) {
		steamKey := os.Getenv("STEAM_TOKEN")
		if steamKey == "" || len(steamKey) != 32 {
			log.Fatalf("Invalid steam token: %s", steamKey)
		}
		steamid.SetKey(steamKey)
		token := os.Getenv("BOT_TOKEN")
		if token == "" || len(token) != 59 {
			log.Fatalf("Invalid bot token: %s", token)
		}
		ctx := context.Background()
		app, err := core.NewApp(ctx, "./db.sqlite", []string{
			"2SVNWuXfv4VewYCuvMXEC3G4CsCY3Ur4Szhihekw4uiK9trfMfxS9QveFuVBQ8a8j",
		})
		if err != nil {
			log.Fatalf("Failed to init app: %v", err)
		}
		ml, err := core.DownloadMasterList()
		if err != nil {
			log.Warnf("Failed to download master list from GH: %s", err)
		}
		app.LoadMasterIDS(ml)
		opts := core.DefaultHTTPOpts()
		opts.Handler = core.NewRouter(app)
		srv := core.NewHTTPServer(opts)
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Errorf("Listener error: %s", err)
			}
		}()
		dg, err := core.NewBot(app, token)
		if err != nil {
			log.Fatalf("Could not connect to discord: %s", err)
		}
		//log.Infof("Add bot linK: %s", core.AddUrl())
		core.Wait(ctx, func(ctx context.Context) error {
			if err := dg.Close(); err != nil {
				log.Errorf("Failed to properly shutdown discord client: %s", err)
			}
			c, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
			defer cancel()
			if err := srv.Shutdown(c); err != nil {
				log.Errorf("Failed to cleanly shutdown http service: %s", err)
			}
			return nil
		})
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tf2bdd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tf2bdd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tf2bdd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
