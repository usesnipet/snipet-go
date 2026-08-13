// Command license is the offline issuer for Snipet Enterprise license keys.
//
//	license gen-keys [-dir DIR] [-force]
//	license issue -key FILE [-payload FILE]
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/usesnipet/snipet/internal/license"
)

const (
	pubFileName  = "ed25519.pub"
	privFileName = "ed25519.key"
)

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "gen-keys":
		err = genKeys(os.Args[2:])
	case "issue":
		err = issue(os.Args[2:])
	case "-h", "-help", "--help", "help":
		usage(os.Stdout)
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		usage(os.Stderr)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage(w io.Writer) {
	fmt.Fprintf(w, `Usage:
  license gen-keys [-dir DIR] [-force]
      Generate an Ed25519 key pair and write ed25519.pub / ed25519.key.

  license issue -key FILE [-payload FILE]
      Sign a JSON payload with the private key and print a LICENSE_KEY.
      Payload is a JSON object:
        {"licensee":"Acme Inc","issued_at":"2026-01-01","expires_at":"2027-01-01","max_tenants":0}
      max_tenants 0 means unlimited. -payload defaults to stdin.

The public key must be baked into internal/license (publicKeyBase64) for
issued keys to verify. The private key must never be committed.
`)
}

func genKeys(args []string) error {
	fs := flag.NewFlagSet("gen-keys", flag.ExitOnError)
	dir := fs.String("dir", ".", "directory to write ed25519.pub and ed25519.key")
	force := fs.Bool("force", false, "overwrite existing key files")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := os.MkdirAll(*dir, 0o700); err != nil {
		return err
	}
	pubPath := filepath.Join(*dir, pubFileName)
	privPath := filepath.Join(*dir, privFileName)
	if !*force {
		for _, p := range []string{pubPath, privPath} {
			if _, err := os.Stat(p); err == nil {
				return fmt.Errorf("%s already exists (use -force to overwrite)", p)
			} else if !os.IsNotExist(err) {
				return err
			}
		}
	}

	pubB64, privB64, err := license.GenerateKeyPair()
	if err != nil {
		return err
	}
	if err := os.WriteFile(pubPath, []byte(pubB64+"\n"), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(privPath, []byte(privB64+"\n"), 0o600); err != nil {
		return err
	}

	fmt.Printf("wrote %s\n", pubPath)
	fmt.Printf("wrote %s\n", privPath)
	fmt.Printf("\npublic key — paste into internal/license/license.go as publicKeyBase64:\n%s\n", pubB64)
	return nil
}

func issue(args []string) error {
	fs := flag.NewFlagSet("issue", flag.ExitOnError)
	keyPath := fs.String("key", "", "path to ed25519.key (required)")
	payloadPath := fs.String("payload", "-", "JSON payload file, or - for stdin")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *keyPath == "" {
		return fmt.Errorf("-key is required")
	}

	privB64, err := os.ReadFile(*keyPath)
	if err != nil {
		return err
	}

	var payloadJSON []byte
	if *payloadPath == "-" {
		payloadJSON, err = io.ReadAll(os.Stdin)
	} else {
		payloadJSON, err = os.ReadFile(*payloadPath)
	}
	if err != nil {
		return err
	}

	key, err := license.Issue(string(privB64), payloadJSON)
	if err != nil {
		return err
	}
	fmt.Println(key)
	return nil
}
