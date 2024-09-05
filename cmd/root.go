/*
Copyright © 2024 Mariano Zunino <marianoz@posteo.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/marianozunino/crapis/internal/aof"
	"github.com/marianozunino/crapis/internal/command"
	"github.com/marianozunino/crapis/internal/logger"
	"github.com/marianozunino/crapis/internal/server"
	"github.com/marianozunino/crapis/internal/store"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var port string
var bind string
var passiveEvictionEnabled bool
var debug bool

var evictionIntervalMs int64
var evictionTimeoutMs int64

var aofEnabled bool
var aofPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crapis",
	Short: "Spawns a redis like server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		logger.ConfigureLogger(debug)
		var dbFile *aof.Aof
		var err error

		if aofEnabled {
			spew.Dump(aofPath)
			if dbFile, err = aof.NewAof(aofPath); err != nil {
				log.Error().Msgf("Error creating AOF: %v", err)
				os.Exit(1)
			}
		}
		db := store.NewStore(
			store.WithPassiveEviction(passiveEvictionEnabled),
			store.WithEvictionInterval(time.Duration(evictionIntervalMs)*time.Millisecond),
			store.WithEvictionTimeout(time.Duration(evictionTimeoutMs)*time.Millisecond),
		)
		executor := command.NewExecutor(db)
		config := server.NewConfig(
			server.WithPort(port),
			server.WithBind(bind),
			server.WithCommandExecutor(executor),
			server.WithAof(dbFile),
		)
		server.NewServer(config).Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&port, "port", "p", "6379", "Port to listen on")
	rootCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "Bind address")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolVarP(&passiveEvictionEnabled, "passive-eviction", "e", true, "Enable passive eviction")

	rootCmd.Flags().Int64VarP(&evictionIntervalMs, "eviction-interval-ms", "i", 250, "Eviction interval in milliseconds")
	rootCmd.Flags().Int64VarP(&evictionTimeoutMs, "eviction-timeout-ms", "t", 10, "Eviction timeout in milliseconds, must be at at most half of eviction-interval-ms")

	rootCmd.Flags().BoolVarP(&aofEnabled, "aof-enabled", "a", false, "Enable AOF")
	rootCmd.Flags().StringVarP(&aofPath, "aof", "f", "database.aof", "Path to AOF file")
}
