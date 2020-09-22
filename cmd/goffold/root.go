package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "goffold",
		Short: "goffold is an opinionated kickstarter for your Go application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Root command invoced: %#v\n", args)
		},
	}
)

func main() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")

	newCmd.PersistentFlags().StringP("name", "n", "", "the name of the go module")
	newCmd.PersistentFlags().StringP("docker", "d", "", "the name of the docker image without the tag")

	rootCmd.AddCommand(newCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
