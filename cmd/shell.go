package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "This is to spawn a shell to play interactively",
	Long: `If you need to interact with your environment, just spawn a shell
	with this subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		shellPath := searchShell("bash")
		shell := exec.Command(shellPath, "-i")
		shell.Stdout = os.Stdout
		shell.Stdin = os.Stdin
		shell.Stderr = os.Stderr
		fmt.Println("Bash shell invoked")
		log.WithFields(log.Fields{
			"subcommand": "shell",
		}).Info("Shell invoked")
		shell.Run()
	},
}

func searchShell(executable string) string {
	path, isPathSet := os.LookupEnv("PATH")
	absPath := ""
	if !isPathSet {
		panic("PATH is not set or is empty!")
	}
	paths := filepath.SplitList(path)
	for _, v := range paths {
		_pathExists, _ := pathExists(v)
		_isInPath, _ := isInPath(v, executable)
		if _pathExists && _isInPath {
			return filepath.Join(v, executable)
		}
	}
	if absPath == "" {
		panic("Executable not found")
	} else {
		return absPath
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func isInPath(inDir string, executable string) (bool, error) {
	result := false
	err := filepath.Walk(inDir, func(dir string, fi os.FileInfo, err error) error {
		if err != nil {
			// prevent panic by handling failure accessing a path
			return err
		}
		log.WithFields(log.Fields{
			"subcommand": "shell",
			"dir":        dir,
			"file":       fi.Name(),
		}).Debug("Currently walking")
		if !fi.IsDir() && fi.Name() == executable {
			result = true
			return filepath.SkipDir
		}
		return nil
	})
	return result, err
}

func init() {
	rootCmd.AddCommand(shellCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shellCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shellCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
