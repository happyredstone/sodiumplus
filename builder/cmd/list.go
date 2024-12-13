package cmd

import (
	"fmt"
	"os"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/list"
	"github.com/spf13/cobra"
)

var (
	Json     bool
	Html     bool
	Markdown bool
	Console  bool

	OutputFile string
)

var ListCommand = CreateListCommand()

func FixArgs() {
	if Json {
		Html = false
		Markdown = false
		Console = false
	}

	if Html {
		Json = false
		Markdown = false
		Console = false
	}

	if Markdown {
		Json = false
		Html = false
		Console = false
	}

	if Console {
		Json = false
		Html = false
		Markdown = false
	}
}

func CreateListCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:          "list",
		Aliases:      []string{"l", "ls"},
		Short:        "List all mods in the modpack.",
		Long:         `List all mods in the modpack.`,
		SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			FixArgs()

			res := ""
			var err error = nil

			if Console {
				res, err = list.Console()

				if err != nil {
					return err
				}
			}

			if Markdown {
				res, err = list.Markdown()

				if err != nil {
					return err
				}
			}

			if Json {
				res, err = list.Json()

				if err != nil {
					return err
				}
			}

			if Html {
				res, err = list.Html()

				if err != nil {
					return err
				}
			}

			if OutputFile != "" {
				out, err := os.Create(OutputFile)

				if err != nil {
					return err
				}

				_, err = out.WriteString(res)

				if err != nil {
					return err
				}

				err = out.Close()

				if err != nil {
					return err
				}
			} else {
				fmt.Println(res)
			}

			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&Json, "json", "j", false, "Output as JSON.")
	cmd.Flags().BoolVarP(&Html, "html", "H", false, "Output as HTML.")
	cmd.Flags().BoolVarP(&Markdown, "markdown", "m", false, "Output as Markdown.")
	cmd.Flags().BoolVarP(&Console, "console", "c", true, "Output to console format.")
	cmd.Flags().StringVarP(&OutputFile, "output", "o", "", "Output file path.")

	cmd.MarkFlagsMutuallyExclusive("json", "html", "markdown", "console")

	return &cmd
}
