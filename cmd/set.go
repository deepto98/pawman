/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"os/user"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		pwd, _ := cmd.Flags().GetString("pwd")

		//1. Encrypting pwd with a symmetric key,nonce
		ciphertext, encKey, nonce, err := Encrypt([]byte(pwd))
		if err != nil {
			log.Fatal(err)
		}
		var encryptedPayload bytes.Buffer
		payload := gob.NewEncoder(&encryptedPayload)
		err = payload.Encode(EncryptPayload{encKey, nonce, ciphertext})
		if err != nil {
			log.Fatal(err)
		}

		//2. Encrypt symmetric key with SSH public key

		//a. Fetch public key
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

		// Fetch value
		var pubKey []byte
		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("pubKey"))
			if err != nil {
				log.Fatal(err)
			}
			// Alternatively, you could also use item.ValueCopy().
			pubKey, err = item.ValueCopy(nil)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("The answer is: %s\n", pubKey)

			return nil
		})
		// fmt.Println("Public KEY locn....................")
		// fmt.Println(encryptedPayload.String())

		enc, err := EncryptWithPublicKey(string(pubKey), encryptedPayload)
		if err != nil {
			log.Fatal(err)
		}

		// Write encrypted pwd to db
		err = db.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte(key), []byte(base64.StdEncoding.EncodeToString(enc)))
			return err
		})
		if err != nil {
			log.Fatal("Unable to write to DB")
		}

		fmt.Printf("Successfully stored the encrypted value for key: '%s'\n", key)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("key", "k", "", "Key for password")
	setCmd.Flags().StringP("pwd", "p", "", "Password string")
	setCmd.MarkFlagRequired("key")
	setCmd.MarkFlagRequired("pwd")
}
