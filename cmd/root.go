package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Lichthagel/shwelcome/anki"
	"github.com/Lichthagel/shwelcome/image"
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		ankiDbPath := viper.GetString("anki.db_path")
		ankiDeckID := viper.GetUint64("anki.deck_id")

		// if imgPath == "" {
		// 	fmt.Println("No image path provided")
		// 	os.Exit(1)
		// }

		if ankiDbPath == "" {
			fmt.Println("No Anki path provided")
			os.Exit(1)
		}

		styleSidePad := lipgloss.NewStyle().Padding(0, 1)

		currentTime := time.Now()

		timeRes := currentTime.Format(time.UnixDate)

		renderedTime := lipgloss.NewStyle().Foreground(lipgloss.Color(catppuccin.Mocha.Overlay2().Hex)).PaddingBottom(1).Render(timeRes)

		var imgBlock string

		if imgPath != "" {
			var err error
			imgBlock, err = image.PathToImgBlock(imgPath, imgWidth, imgHeight)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		db, err := sql.Open("sqlite3", ankiDbPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ankiCard, err := anki.RandomCard(db, ankiDeckID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		renderAnkiCard := anki.RenderCard(ankiCard)

		rightBlock := lipgloss.JoinVertical(0, renderedTime, renderAnkiCard)

		var result string

		if imgPath != "" {
			result = lipgloss.JoinHorizontal(lipgloss.Center, styleSidePad.Render(imgBlock), styleSidePad.Render(rightBlock))
		} else {
			result = styleSidePad.Render(rightBlock)
		}

		result = lipgloss.NewStyle().Padding(1).PaddingBottom(0).Render(result)

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

	rootCmd.Flags().StringP("anki-db", "a", "", "Path to an Anki-exported text file")
	rootCmd.Flags().Uint64P("deck-id", "d", 0, "Deck ID to use")

	viper.BindPFlag("image.path", rootCmd.Flags().Lookup("image"))
	viper.BindPFlag("image.width", rootCmd.Flags().Lookup("width"))
	viper.BindPFlag("image.height", rootCmd.Flags().Lookup("height"))
	viper.BindPFlag("anki.db_path", rootCmd.Flags().Lookup("anki-db"))
	viper.BindPFlag("anki.deck_id", rootCmd.Flags().Lookup("deck-id"))

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
		home, err := os.UserConfigDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".shwelcome" (without extension).
		viper.AddConfigPath(path.Join(home, "shwelcome"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("shwelcome.yml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
