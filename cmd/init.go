/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	//Functionality of command
	Run: func(cmd *cobra.Command, args []string) {

		//Fetch values of flags
		confDir, _ := cmd.Flags().GetString("db")
		privKey, _ := cmd.Flags().GetString("priv")
		pubKey, _ := cmd.Flags().GetString("pub")

		//Validation
		if _, err := os.Stat(privKey); os.IsNotExist(err) {
			log.Fatal("Private Key does not exist")
		}

		if _, err := os.Stat(pubKey); os.IsNotExist(err) {
			log.Fatal("Public Key does not exist")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Get home directory of user
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Database - $HOME/.pawman/
	rootDir := fmt.Sprintf("%s/.pawman", user.HomeDir)

	//Default SSH keys location :
	//$HOME/.ssh/id_rsa --private key
	//$HOME/.ssh/id_rsa.pub --public key

	privKey := fmt.Sprintf("%s/.ssh/id_rsa", user.HomeDir)
	pubKey := fmt.Sprintf("%s/.ssh/id_rsa.pub", user.HomeDir)

	// Define Flags for rootDir, private, public keys
	initCmd.Flags().StringP("db", "d", rootDir, "Absolute Path")
	initCmd.Flags().StringP("priv", "r", privKey, "Absolute Path")
	initCmd.Flags().StringP("pub", "u", pubKey, "Absolute Path")

}
