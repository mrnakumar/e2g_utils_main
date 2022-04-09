package main

import (
	"flag"
	"github.com/mrnakumar/e2g_utils"
)

type utilFlags struct {
	generateKey bool
	toEncode    string
}

func main() {
	utilFlags := parseFlags()
	if utilFlags.generateKey {
		e2g_utils.GenerateX25519Identity()
	}
	if len(utilFlags.toEncode) > 0 {
		e2g_utils.Base64Encode(utilFlags.toEncode)
	}
}

func parseFlags() utilFlags {
	generateKey := flag.String("generate_X25519_key", "", "generate key? [true | false]")
	toEncode := flag.String("to_encode", "", "input string to decode.")
	flag.Parse()
	shouldGenerateKey := false
	if *generateKey == "true" {
		shouldGenerateKey = true
	}
	return utilFlags{
		toEncode:    *toEncode,
		generateKey: shouldGenerateKey,
	}
}
