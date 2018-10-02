package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/leophys/aurbuilder/utils"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var baseDir string
var pkgName string
var options []string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Fetch and build a package from the AUR",
	Long: `This command clones a git repo supposed to be located at

	https://aur.archlinux.org/$pkg.git

	and builds it with "makepkg -s".`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]
		fmt.Println("build called")
		log.WithFields(log.Fields{
			"subcommand": "build",
		}).Info("Build begun")
		var err error
		err = fetchSources(pkgName, baseDir)
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": "build",
				"error":      err,
				"pkg":        pkgName,
			}).Fatal("Fetch pkg failed")
		}
		err = updatePacman()
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": "build",
				"error":      err,
				"pkg":        pkgName,
			}).Error("Package database sync failed")
		}
		err = maybeEditPkgbuild(filepath.Join(baseDir, pkgName))
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": "build",
				"error":      err,
				"pkg":        pkgName,
			}).Fatal("PKGBUILD edit failed")
		}
		err = doBuild(filepath.Join(baseDir, pkgName))
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": "build",
				"error":      err,
				"pkg":        pkgName,
			}).Fatal("Build failed")
		}
	},
}

func fetchSources(pkgName string, baseDir string) error {
	var archAURrepoURL string
	var err error
	archAURrepoURL = "https://aur.archlinux.org/" + pkgName + ".git"
	log.WithFields(log.Fields{
		"subcommand": "build",
		"url":        archAURrepoURL,
		"pkg":        pkgName,
	}).Debug("Arch pkg url")
	repoDir := filepath.Join(baseDir, pkgName)
	gitCloneCmd := exec.Command("git", "clone", archAURrepoURL, repoDir)
	gitPullCmd := exec.Command("git", "pull")
	gitStashCmd := exec.Command("git", "stash")
	_, err = os.Stat(baseDir)
	if os.IsNotExist(err) {
		home, err := homedir.Dir()
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": nil,
				"err":        err,
			}).Fatal("Home dir not found")
		}
		baseDir = home
		log.WithFields(log.Fields{
			"subcommand": "build",
			"baseDir":    baseDir,
			"home":       home,
		}).Debug("baseDir is home")
	}
	log.WithFields(log.Fields{
		"subcommand": "build",
		"repoDir":    repoDir,
		"baseDir":    baseDir,
		"pkgName":    pkgName,
	}).Debug("Input data check")
	_, err = os.Stat(repoDir)
	if os.IsNotExist(err) {
		if logLevel == "debug" {
			gitCloneCmd.Stdout = os.Stdout
			gitCloneCmd.Stderr = os.Stderr
		}
		err = gitCloneCmd.Run()
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": "build",
				"error":      err,
				"repo url":   archAURrepoURL,
			}).Fatal("Could not clone repo")
		}
	} else if err != nil {
		log.WithFields(log.Fields{
			"subcommand": "build",
			"error":      err,
		}).Fatal("Could not find directory")
	} else {
		os.Chdir(repoDir)
		log.WithFields(log.Fields{
			"subcommand": "build",
			"repo":       repoDir,
			"pkg":        pkgName,
		}).Info("Stashing the local repo")
		err = gitStashCmd.Run()
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": "build",
				"pkg":        pkgName,
				"error":      err,
			}).Warn("Could not stash")
			return err
		}
		log.WithFields(log.Fields{
			"subcommand": "build",
			"repo":       repoDir,
			"pkg":        pkgName,
		}).Info("Pulling from remote to the local repo")
		err = gitPullCmd.Run()
		if err != nil {
			log.WithFields(log.Fields{
				"subcommand": "build",
				"pkg":        pkgName,
				"error":      err,
			}).Warn("Could not pull")
			return err
		}
		return nil
	}
	_, err = os.Stat(repoDir)
	return err
}

func maybeEditPkgbuild(basePath string) error {
	defaultEditor := os.Getenv("EDITOR")
	if defaultEditor == "" {
		defaultEditor = "vi"
	}
	log.WithFields(log.Fields{
		"subcommand":    "build",
		"defaultEditor": defaultEditor,
	}).Debug("Default editor")
	editCmd := exec.Command(defaultEditor, filepath.Join(basePath, "PKGBUILD"))
	editCmd.Stdin = os.Stdin
	editCmd.Stdout = os.Stdout
	editCmd.Stderr = os.Stderr
	promptQuestion := "Would you like to edit the PKGBUILD? [y/N]"
	doEdit, err := utils.AskConfirmation(promptQuestion, false)
	if err != nil {
		log.WithFields(log.Fields{
			"subcommand": "build",
			"err":        err,
		}).Error("Error while reading response")
	}
	if doEdit {
		err = editCmd.Run()
	} else {
		err = nil
	}
	return err
}

func updatePacman() error {
	updatePacmanCmd := exec.Command("/usr/bin/sh", "-c", "sudo pacman -Syu --noconfirm")
	updatePacmanCmd.Stdin = os.Stdin
	updatePacmanCmd.Stdout = os.Stdout
	updatePacmanCmd.Stderr = os.Stderr
	err := updatePacmanCmd.Run()
	return err
}

func doBuild(repoDir string) error {
	os.Chdir(repoDir)
	cwd, _ := os.Getwd()
	log.WithFields(log.Fields{
		"subcommand": "build",
		"pwd":        cwd,
	}).Info("Show build pwd")
	fmt.Println("pwd: ", cwd)
	makepkgCmd := exec.Command("makepkg", "-s")
	makepkgCmd.Stdin = os.Stdin
	makepkgCmd.Stdout = os.Stdout
	makepkgCmd.Stderr = os.Stderr
	err := makepkgCmd.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"subcommand": "build",
			"pwd":        cwd,
			"error":      err,
			"test":       "test_field",
		}).Error("Makepkg failed")
	}
	return err
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")
	pwd, _ := os.Getwd()
	log.WithFields(log.Fields{
		"subcommand": "build",
		"pwd":        pwd,
	}).Debug("The current working directory")
	defaultBasePath := filepath.Join(pwd, "store")
	buildCmd.Flags().StringVar(&baseDir, "basepath", defaultBasePath,
		`Set the path where to store pkg
		git repositories cloned from source.`)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
