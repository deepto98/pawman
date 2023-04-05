/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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

		//Encrypting pwd with a symmetric key,nonce
		ciphertext, encKey, nonce, err := Encrypt([]byte(pwd))
		//Encrypt symmetric key with SSH public key
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("key", "k", "", "Key for password")
	setCmd.Flags().StringP("pwd", "p", "", "Password string")
	setCmd.MarkFlagRequired("key")
	setCmd.MarkFlagRequired("pwd")
}
