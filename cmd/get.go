/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/user"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get pwd key from cli
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			log.Fatal("Invalid Key")
		}

		// Check if pwd key exists iun DB and Get priv key location from db
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		// Database - $HOME/.pawman/
		dbDir := fmt.Sprintf("%s/.pawman", user.HomeDir)
		db, err := badger.Open(badger.DefaultOptions(dbDir).WithLogger(nil))
		if err != nil {
			log.Fatal("Unable to create database")
		}

		defer db.Close()

		// Check if pwd for key exists
		var encryptedPwd []byte
		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(key))
			if err != nil {
				log.Fatal("Key doesn't exist")
			}
			encryptedPwd, err = item.ValueCopy(nil)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})

		// Fetch private key from db
		var privKey []byte
		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("privKey"))
			if err != nil {
				log.Fatal(err)
			}
			// Alternatively, you could also use item.ValueCopy().
			privKey, err = item.ValueCopy(nil)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("The answer is: %s\n", pubKey)

			return nil
		})

		// Decrypt with private key
		unwrapped, err := DecryptWithPrivateKey(string(privKey), encryptedPwd)
		if err != nil {
			log.Fatal(err)
		}
		var unwrap bytes.Buffer
		unwrap.WriteString(string(unwrapped))
		dec, err := DecryptBox(unwrap)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dec))

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("key", "k", "", "Key for which password is to be retrieved")
	getCmd.MarkFlagRequired("key")
}
