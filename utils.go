package main

import (
	"filippo.io/age"
	"flag"
	"github.com/mrnakumar/e2g_utils"
	"io/ioutil"
	"log"
	"strings"
)

const decryptIdentityFlag = "decrypt_identity"
const toDecodeFlag = "to_decode"

type utilFlags struct {
	generateKey bool
	toEncode    string
	toDecode    string
	decrypt     *decryptFlag
}

type decryptFlag struct {
	input    string
	output   string
	identity *age.X25519Identity
}

func main() {
	utilFlags := parseFlags()
	if utilFlags.generateKey {
		keyPair := e2g_utils.GenerateX25519Identity()
		log.Printf("Public key: '%s'\n", keyPair.Public)
		log.Printf("Private key: '%s'\n", keyPair.Private)
	}
	if len(utilFlags.toEncode) > 0 {
		log.Printf("input '%s' is encoded to '%s'", utilFlags.toEncode, e2g_utils.Base64Encode(utilFlags.toEncode))
	}
	if len(utilFlags.toDecode) > 0 {
		log.Printf("input '%s' is decoded to '%s'", utilFlags.toDecode, e2g_utils.Base64DecodeWithKill(utilFlags.toDecode, toDecodeFlag))
	}
	if utilFlags.decrypt != nil {
		decrypt(utilFlags)
	}
}

func decrypt(utilFlags utilFlags) {
	input := utilFlags.decrypt.input
	decrypted, err := e2g_utils.Decrypt(input, utilFlags.decrypt.identity)
	if err != nil {
		log.Fatalf("failed to decrypt file at path '%s'. caused by: '%v'", input, err)
	}
	output := utilFlags.decrypt.output
	err = ioutil.WriteFile(output, decrypted, 0400)
	if err != nil {
		log.Fatalf("failed to save decrypted content to path '%s'. caused by: '%v'", output, err)
	}
}

func parseFlags() utilFlags {
	generateKey := flag.String("generate_X25519_key", "", "generate key? [true | false]")
	toEncode := flag.String("to_encode", "", "input string to encode")
	toDecode := flag.String(toDecodeFlag, "", "input string to be base64 decode")
	decryptInput := flag.String("decrypt_input", "", "path of the file to decrypt")
	decryptOutput := flag.String("decrypt_output", "", "path to save decrypted file contents")
	decryptIdentity := flag.String(decryptIdentityFlag, "", "private key file path to use for decryption")

	flag.Parse()
	shouldGenerateKey := false
	if *generateKey == "true" {
		shouldGenerateKey = true
	}
	return utilFlags{
		toEncode:    *toEncode,
		toDecode:    *toDecode,
		generateKey: shouldGenerateKey,
		decrypt:     getDecryptFlag(decryptInput, decryptOutput, decryptIdentity),
	}
}

func getDecryptFlag(decryptInput *string, decryptOutPath *string, decryptIdentity *string) *decryptFlag {
	var decrypt *decryptFlag = nil
	if len(*decryptInput) > 0 {
		decrypt = &decryptFlag{}
		decrypt.input = *decryptInput
		decrypt.output = *decryptInput
		if len(*decryptOutPath) > 0 {
			decrypt.output = *decryptOutPath
		}
		if len(*decryptIdentity) == 0 {
			log.Fatalf("missing flag '%s'", decryptIdentityFlag)
		}
		identity, err := parseIdentity(*decryptIdentity)
		if err != nil {
			log.Fatalf("failed to parse flag '%s'. caused by: '%v'", decryptIdentityFlag, err)
		}
		decrypt.identity = identity
	}
	return decrypt
}

func parseIdentity(identityFilePath string) (*age.X25519Identity, error) {
	contents, err := ioutil.ReadFile(identityFilePath)
	if err != nil {
		return nil, err
	}
	contentsDecoded := e2g_utils.Base64DecodeWithKill(strings.TrimSuffix(string(contents), "\n"), decryptIdentityFlag)
	return age.ParseX25519Identity(contentsDecoded)
}
