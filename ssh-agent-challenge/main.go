package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"text/tabwriter"
)

// agentFromAuthSock returns an Agent that talks to the local ssh-agent on SSH_AUTH_SOCK
func agentFromAuthSock() (agent.Agent, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, errors.New("SSH_AUTH_SOCK environment variable is not set. Verify ssh-agent is running. See https://github.com/coreos/fleet/blob/master/Documentation/using-the-client.md for help.")
	}

	a, err := net.Dial("unix", sock)
	if err != nil {
		return nil, err
	}

	return agent.NewClient(a), nil
}

func listKeys(a agent.Agent) {
	keys, err := a.List()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for i, k := range keys {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i, k.Format, base64.StdEncoding.EncodeToString(k.Blob), k.Comment)
		w.Flush()
	}
}

func sign(a agent.Agent, b string) {
	keys, err := a.List()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	data, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for i, k := range keys {
		sig, err := a.Sign(k, data)
		if err != nil {
			println(err)
			os.Exit(1)
		}
		fmt.Fprintf(w, "%d\t%i\t%s\t%s\t%s\n", i, sig.Format, base64.StdEncoding.EncodeToString(sig.Blob))
		w.Flush()
	}
}

func verify(a agent.Agent, b string, sFormat string, s string) {
	keys, err := a.List()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	data, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	sigData, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	sig := ssh

	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for i, k := range keys {
		sig, err := k.Verify(data, sig)
		if err != nil {
			println(err)
			os.Exit(1)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i, sig.Format, base64.StdEncoding.EncodeToString(sig.Blob))
		w.Flush()
	}
}

func main() {
	a, err := agentFromAuthSock()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	var toolCmd = &cobra.Command{
		Use:   "ssh-agent-tool",
		Short: "ssh-agent-tool is a swiss army knife for an ssh-agent",
		Long: `ssh clients have an agent for protecting and forwarding access to challenge-response
of public keys. This tool exercises the agent and exposes these
features for use by other tools besides ssh clients.`,
	}

	var listKeysCmd = &cobra.Command{
		Use:   "list-keys",
		Short: "List all public keys added in the agent",
		Run: func(cmd *cobra.Command, args []string) {
			listKeys(a)
		},
	}

	var signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Sign a base64 encoded string",
		Run: func(cmd *cobra.Command, args []string) {
			sign(a, args[0])
		},
	}

	var verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify a base64 encoded string and signature",
		Run: func(cmd *cobra.Command, args []string) {
			verify(a, args[0], args[1], args[2])
		},
	}

	toolCmd.AddCommand(listKeysCmd)
	toolCmd.AddCommand(signCmd)
	toolCmd.AddCommand(verify)
	toolCmd.Execute()
}
