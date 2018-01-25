// Copyright © 2017 Fuxin Hao <haofxpro@gmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/FX-HAO/crypto-market-overwatch/collector"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Host     string
	Port     int
	Interval int
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "crypto_market_overwatch",
	Short: "Tracking crypto market cap",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		http.Handle("/metrics", promhttp.Handler())
		collector.NewCollector(Interval).Start()
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			s := <-c
			log.Infof("Receive signal: %v", s)

			switch {
			default:
				os.Exit(-1)
			case s == syscall.SIGHUP:
				os.Exit(1)
			case s == syscall.SIGINT:
				os.Exit(2)
			case s == syscall.SIGQUIT:
				os.Exit(3)
			case s == syscall.SIGTERM:
				os.Exit(0xf)
			}
		}()

		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", Host, Port), nil))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().BoolP("debug", "d", true, "Debug mode")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().StringVarP(&Host, "host", "H", "0.0.0.0", "Host")
	RootCmd.Flags().IntVarP(&Port, "port", "p", 80, "Port")
	RootCmd.Flags().IntVarP(&Interval, "interval", "i", 30, "Interval")
}
