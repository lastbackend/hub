//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/registry/pkg/api"
	"github.com/lastbackend/registry/pkg/builder"
	"github.com/lastbackend/registry/pkg/controller"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config string
	debug  int

	// CLI - default command line endpoint
	CLI = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  "",

		// parse the config if one is provided, or use the defaults. Set the backend
		// driver to be used
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			// if --config is passed, attempt to parse the config file
			if config != "" {

				// get the filepath
				abs, err := filepath.Abs(config)
				if err != nil {
					log.Errorf("Error reading filepath: %s \n", err)
				}

				// get the config name
				base := filepath.Base(abs)

				// get the path
				path := filepath.Dir(abs)

				//
				viper.SetConfigName(strings.Split(base, ".")[0])
				viper.AddConfigPath(path)

				// Find and read the config file; Handle errors reading the config file
				if err := viper.ReadInConfig(); err != nil {
					log.Fatalf("Failed to read config file: %s\n", err)
				}
			}
		},

		// either run hoarder as a server, or run it as a CLI depending on what flags
		// are provided
		Run: func(cmd *cobra.Command, args []string) {
			log.New(viper.GetInt("verbose"))

			var (
				done    = make(chan bool, 1)
				apps    = make(chan bool)
				wait    = 0
				daemons = map[string]func() bool{
					"api":        api.Daemon,
					"controller": controller.Daemon,
					"builder":    builder.Daemon,
				}
			)

			components := []string{"api"}

			if len(args) > 0 {
				components = args
			}

			for _, app := range components {
				go func(app string) {
					if _, ok := daemons[app]; ok {
						wait++
						apps <- daemons[app]()
					}
				}(app)
			}

			go func() {
				for {
					select {
					case <-apps:
						wait--
						if wait == 0 {
							done <- true
							return
						}
					}
				}
			}()

			<-done

		},
	}
)

func init() {

	// set config defaults
	viper.SetDefault("garbage-collect", false)

	// local flags;
	CLI.Flags().StringVarP(&config, "config", "c", "", "/path/to/config.yml")
	CLI.Flags().IntVarP(&debug, "verbose", "v", 0, "verbose level")

	_ = viper.BindPFlag("verbose", CLI.Flags().Lookup("verbose"))

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

}

func main() {
	if err := CLI.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
