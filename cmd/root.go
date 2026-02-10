package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "./sentinel",
	Short: "⌀ Sentinel - The vigilant guard for code authenticity.",
	Long: `⌀ Sentinel - The vigilant guard for code authenticity.

SUPPORTED LANGUAGES (early support):
  Python, Java, JavaScript, TypeScript, Go, Rust, C/C++, Ruby, PHP, C#, Kotlin, Swift

QUICK START:
  # Scan current directory
    ./sentinel scan --path .

  # Scan with treshold specified
    ./sentinel scan --path ./your-target-dir --threshold 0.8

  # CI/CD integration that cancel build if an AI-slop is detected
    ./sentinel scan --git-diff main --fail-on-detection`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sentinel.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output with detailed logging")
	rootCmd.PersistentFlags().Bool("json", false, "output results in JSON format for CI/CD integration")

	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding verbose flag: %v\n", err)
	}
	if err := viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding json flag: %v\n", err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sentinel")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && viper.GetBool("verbose") {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
