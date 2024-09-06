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
	"fmt"
	"os"
	"time"

	"github.com/marianozunino/crapis/internal/aof"
	"github.com/marianozunino/crapis/internal/command"
	"github.com/marianozunino/crapis/internal/logger"
	"github.com/marianozunino/crapis/internal/server"
	"github.com/marianozunino/crapis/internal/store"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type config struct {
	port                   string
	bind                   string
	debug                  bool
	passiveEvictionEnabled bool
	evictionInterval       time.Duration
	evictionTimeout        time.Duration
	aofEnabled             bool
	aofPath                string
}

func newRootCmd() *cobra.Command {
	cfg := &config{}

	rootCmd := &cobra.Command{
		Use:   "crapis",
		Short: "Spawns a Redis-like server",
		Long:  `CRAPIs is a lightweight, Redis-like server implementation.`,
		Run: func(cmd *cobra.Command, args []string) {
			run(cfg)
		},
	}

	rootCmd.Flags().StringVarP(&cfg.port, "port", "p", "6379", "Port to listen on")
	rootCmd.Flags().StringVarP(&cfg.bind, "bind", "b", "0.0.0.0", "Bind address")
	rootCmd.Flags().BoolVarP(&cfg.debug, "debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolVarP(&cfg.passiveEvictionEnabled, "passive-eviction", "e", true, "Enable passive eviction")
	rootCmd.Flags().DurationVarP(&cfg.evictionInterval, "eviction-interval", "i", 250*time.Millisecond, "Eviction interval")
	rootCmd.Flags().DurationVarP(&cfg.evictionTimeout, "eviction-timeout", "t", 10*time.Millisecond, "Eviction timeout (must be at most half of eviction-interval)")
	rootCmd.Flags().BoolVarP(&cfg.aofEnabled, "aof-enabled", "a", true, "Enable AOF")
	rootCmd.Flags().StringVarP(&cfg.aofPath, "aof", "f", "database.aof", "Path to AOF file")

	return rootCmd
}

func Execute() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg *config) error {
	if err := validateConfig(cfg); err != nil {
		return err
	}

	logger.ConfigureLogger(cfg.debug)

	db := setupStore(cfg)
	executor := command.NewExecutor(db)
	srv := setupServer(cfg, executor)

	log.Info().Msg("Starting server...")
	return srv.Run()
}

func validateConfig(cfg *config) error {
	if cfg.evictionTimeout > cfg.evictionInterval/2 {
		return fmt.Errorf("eviction timeout must be at most half of eviction interval")
	}
	return nil
}

func setupStore(cfg *config) *store.Store {
	return store.NewStore(
		store.WithPassiveEviction(cfg.passiveEvictionEnabled),
		store.WithEvictionInterval(cfg.evictionInterval),
		store.WithEvictionTimeout(cfg.evictionTimeout),
	)
}

func setupServer(cfg *config, executor command.Executor) *server.Server {
	var dbFile *aof.Aof
	var err error

	if cfg.aofEnabled {
		dbFile, err = aof.NewAof(cfg.aofPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating AOF")
		}
	}

	return server.NewServer(
		server.WithPort(cfg.port),
		server.WithBind(cfg.bind),
		server.WithCommandExecutor(executor),
		server.WithAof(dbFile),
	)
}
