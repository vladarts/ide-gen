package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	appName = "ide-gen"
	//: env variables

	//: exit codes
	exitOk           = 0
	exitCommandError = 2
)

var (
	logger = logrus.New()

	appVersion = "dev"

	// Exit function
	Exit = func(code int) {
		os.Exit(code)
	}

	rootCmd = &cobra.Command{
		Use:   appName,
		Short: "IntelliJ IDEA and other familiar IDE's project manager",
	}
	versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"V"},
		Example: "version",
		Run: func(cmd *cobra.Command, args []string) {
			version := fmt.Sprintf(
				"%[1]s version: %[2]s, %[3]s/%[4]s %[5]s",
				appName, appVersion, runtime.GOOS, runtime.GOARCH, runtime.Version())
			fmt.Println(version)
		},
	}
)

func exit(err error) {
	if err == nil {
		Exit(exitOk)
		return
	}
	Exit(exitCommandError)
}

func init() {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(NewGenerateCommand().Register())
	rootCmd.AddCommand(NewJsonSchemaCommandCommand().Register())
}

func main() {
	err := rootCmd.Execute()
	exit(err)
}
