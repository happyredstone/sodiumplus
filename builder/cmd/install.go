package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/helpers"
	"github.com/google/go-github/v56/github"
	"github.com/spf13/cobra"
)

var (
	JavaPath string
	PackPath string
)

var InstallCommand = CreateInstallCommand()

func CreateInstallCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:          "install",
		Aliases:      []string{"i"},
		Short:        "Install the modpack (development version).",
		Long:         `Install the modpack (development version).`,
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			installerPath := path.Join(os.TempDir(), "packwiz-installer-bootstrap.jar")
			client := github.NewClient(nil)

			err := helpers.Refresh()

			if err != nil {
				return err
			}

			fmt.Println("Downloading bootstrapper...")

			rel, _, err := client.Repositories.GetLatestRelease(context.Background(), "packwiz", "packwiz-installer-bootstrap")

			if err != nil {
				return err
			}

			for _, asset := range rel.Assets {
				name := *asset.Name

				if strings.HasSuffix(name, ".jar") {
					url := *asset.BrowserDownloadURL

					err = helpers.DownloadFile(installerPath, url)

					if err != nil {
						return err
					}

					break
				}
			}

			fmt.Println("Creating install directory...")

			os.MkdirAll(strings.Join(args, " "), os.ModePerm)

			defer os.Remove(installerPath)

			fmt.Println("Running installer...")

			cexec := exec.Command(JavaPath, "-jar", installerPath, "-g", PackPath)

			cexec.Dir = strings.Join(args, " ")
			cexec.Stdout = os.Stdout
			cexec.Stderr = os.Stderr

			err = cexec.Start()

			if err != nil {
				return err
			}

			return cexec.Wait()
		},
	}

	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVarP(&JavaPath, "java", "j", "java", "The path to the Java executable")
	cmd.Flags().StringVarP(&PackPath, "file", "f", cwd+"/pack.toml", "The path to the pack.toml file")

	return &cmd
}
