// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package googlecompute

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

// processPrivateKeyFile.
func processPrivateKeyFile(privateKeyFile, passphrase string) ([]byte, error) {
	rawPrivateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Failed loading private key file: %s", err)
	}
	PEMBlock, _ := pem.Decode(rawPrivateKeyBytes)
	if PEMBlock == nil {
		return nil, fmt.Errorf("%s does not contain a vaild private key", privateKeyFile)
	}
	if x509.IsEncryptedPEMBlock(PEMBlock) {
		if passphrase == "" {
			return nil, errors.New("a passphrase must be specified when using an encrypted private key")
		}
		decryptedPrivateKeyBytes, err := x509.DecryptPEMBlock(PEMBlock, []byte(passphrase))
		if err != nil {
			return nil, fmt.Errorf("Failed decrypting private key: %s", err)
		}
		b := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: decryptedPrivateKeyBytes,
		}
		return pem.EncodeToMemory(b), nil
	}
	return rawPrivateKeyBytes, nil
}
