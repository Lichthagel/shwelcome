package cmd

import (
	"fmt"
	"os"

	"github.com/Lichthagel/shwelcome/image"
	"github.com/charmbracelet/lipgloss"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var cfgPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shwelcome",
	Short: "Prints an opinionated welcome message",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		imgPath := viper.GetString("image.path")
		imgWidth := viper.GetUint("image.width")
		imgHeight := viper.GetUint("image.height")

		if imgPath == "" {
			fmt.Println("No image path provided")
			os.Exit(1)
		}

		termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		imgPadded, err := image.PathToPaddedCode(imgPath, imgWidth, imgHeight)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		center := lipgloss.NewStyle().Align(lipgloss.Center).Width(termWidth)

		result := lipgloss.JoinHorizontal(lipgloss.Center, imgPadded, "Hello, World!")

		result = center.Render(result)

		fmt.Println(result)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file (default is $HOME/.shwelcome.yaml)")

	rootCmd.Flags().StringP("image", "i", "", "Path to an image file")
	rootCmd.Flags().Uint("width", 0, "Width of the image")
	rootCmd.Flags().Uint("height", 0, "Height of the image")

	viper.BindPFlag("image.path", rootCmd.Flags().Lookup("image"))
	viper.BindPFlag("image.width", rootCmd.Flags().Lookup("width"))
	viper.BindPFlag("image.height", rootCmd.Flags().Lookup("height"))

	viper.SetDefault("image.width", 20)
	viper.SetDefault("image.height", 10)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgPath)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".shwelcome" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".shwelcome")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
