// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package googlecompute

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"code.google.com/p/go.crypto/ssh"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// stepCreateSSHKey represents a Packer build step that generates SSH key pairs.
type stepCreateSSHKey int

// Run executes the Packer build step that generates SSH key pairs.
func (s *stepCreateSSHKey) Run(state multistep.StateBag) multistep.StepAction {
	var (
		ui = state.Get("ui").(packer.Ui)
	)
	ui.Say("Creating temporary ssh key for instance...")
	priv, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		err := fmt.Errorf("Error creating temporary ssh key: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	priv_der := x509.MarshalPKCS1PrivateKey(priv)
	priv_blk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priv_der,
	}
	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		err := fmt.Errorf("Error creating temporary ssh key: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("ssh_private_key", string(pem.EncodeToMemory(&priv_blk)))
	state.Put("ssh_public_key", string(ssh.MarshalAuthorizedKey(pub)))
	return multistep.ActionContinue
}

// Cleanup.
// Nothing to clean up. SSH keys are associated with a single GCE instance.
func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {}
