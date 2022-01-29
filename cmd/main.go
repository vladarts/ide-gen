package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	appName = "idea-pm"
	//: env variables

	//: exit codes
	exitOk           = 0
	exitCommandError = 2
)

var (
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
		// this return makes sense only for testing, due to
		// there's no real system exit from this function, thus far
		// running in tests it will continue to follow the code sequence.
		return
	}
	Exit(exitCommandError)
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(NewGenerateCommand().Register())
}

func main() {
	err := rootCmd.Execute()
	exit(err)
}
