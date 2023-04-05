/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Command to initialize database and private,public keys",
	Long:  ``,

	//Functionality of command
	Run: func(cmd *cobra.Command, args []string) {

		//Fetch values of flags
		confDir, _ := cmd.Flags().GetString("db")
		privKeyLoc, _ := cmd.Flags().GetString("priv")
		pubKeyLoc, _ := cmd.Flags().GetString("pub")

		//Validation - check if files exists
		if _, err := os.Stat(privKeyLoc); os.IsNotExist(err) {
			log.Fatal("Private Key does not exist")
		}

		if _, err := os.Stat(pubKeyLoc); os.IsNotExist(err) {
			log.Fatal("Public Key does not exist")
		}

		// If conf dir doesn't exist, create the dir and the database
		if _, err := os.Stat(confDir); os.IsNotExist(err) {
			err := os.Mkdir(confDir, os.ModePerm)

			if err != nil {
				log.Fatal("Unable to create directory")
			}
			fmt.Println("Created directory")

			dbLoc := fmt.Sprintf("%s/pawman.db", confDir)

			db, err := badger.Open(badger.DefaultOptions(dbLoc))
			if err != nil {
				log.Fatal("Unable to create database")
			}

			defer db.Close()

			//Store private and public key locations
			err = db.Update(func(txn *badger.Txn) error {
				err := txn.Set([]byte("privKey"), []byte(privKeyLoc))
				return err
			})

			if err != nil {
				log.Fatalf(err.Error())
			}

			err = db.Update(func(txn *badger.Txn) error {
				err := txn.Set([]byte("pubKey"), []byte(pubKeyLoc))
				return err
			})

			if err != nil {
				log.Fatalf(err.Error())
			}

			fmt.Println("Pawman Initialized")
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

	// Define Flags for rootDir, private, public keys with default values
	initCmd.Flags().StringP("db", "d", rootDir, "Absolute Path")
	initCmd.Flags().StringP("priv", "r", privKey, "Absolute Path")
	initCmd.Flags().StringP("pub", "u", pubKey, "Absolute Path")

}
