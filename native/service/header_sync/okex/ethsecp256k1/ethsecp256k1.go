/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */
package ethsecp256k1

import (
	"bytes"
	"crypto/ecdsa"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	tmcrypto "github.com/tendermint/tendermint/crypto"
)

const (
	// PrivKeySize defines the size of the PrivKey bytes
	PrivKeySize = 32
	// KeyType is the string constant for the EthSecp256k1 algorithm
	KeyType = "eth_secp256k1"
)

// Amino encoding names
const (
	// PrivKeyName defines the amino encoding name for the EthSecp256k1 private key
	PrivKeyName = "ethermint/PrivKeyEthSecp256k1"
	// PubKeyName defines the amino encoding name for the EthSecp256k1 public key
	PubKeyName = "ethermint/PubKeyEthSecp256k1"
)

// ----------------------------------------------------------------------------
// secp256k1 Private Key

var _ tmcrypto.PrivKey = PrivKey{}

// PrivKey defines a type alias for an ecdsa.PrivateKey that implements
// Tendermint's PrivateKey interface.
type PrivKey []byte

// GenerateKey generates a new random private key. It returns an error upon
// failure.
func GenerateKey() (PrivKey, error) {
	priv, err := ethcrypto.GenerateKey()
	if err != nil {
		return PrivKey{}, err
	}

	return PrivKey(ethcrypto.FromECDSA(priv)), nil
}

// PubKey returns the ECDSA private key's public key.
func (privkey PrivKey) PubKey() tmcrypto.PubKey {
	ecdsaPKey := privkey.ToECDSA()
	return PubKey(ethcrypto.CompressPubkey(&ecdsaPKey.PublicKey))
}

// Bytes returns the raw ECDSA private key bytes.
func (privkey PrivKey) Bytes() []byte {
	return CryptoCodec.Amino.MustMarshalBinaryBare(privkey)
}

// Sign creates a recoverable ECDSA signature on the secp256k1 curve over the
// Keccak256 hash of the provided message. The produced signature is 65 bytes
// where the last byte contains the recovery ID.
func (privkey PrivKey) Sign(msg []byte) ([]byte, error) {
	return ethcrypto.Sign(ethcrypto.Keccak256Hash(msg).Bytes(), privkey.ToECDSA())
}

// Equals returns true if two ECDSA private keys are equal and false otherwise.
func (privkey PrivKey) Equals(other tmcrypto.PrivKey) bool {
	if other, ok := other.(PrivKey); ok {
		return bytes.Equal(privkey.Bytes(), other.Bytes())
	}

	return false
}

// ToECDSA returns the ECDSA private key as a reference to ecdsa.PrivateKey type.
// The function will panic if the private key is invalid.
func (privkey PrivKey) ToECDSA() *ecdsa.PrivateKey {
	key, err := ethcrypto.ToECDSA(privkey)
	if err != nil {
		panic(err)
	}
	return key
}

func (privkey PrivKey) Type() string {
	return KeyType
}

// ----------------------------------------------------------------------------
// secp256k1 Public Key

var _ tmcrypto.PubKey = (*PubKey)(nil)

// PubKey defines a type alias for an ecdsa.PublicKey that implements Tendermint's PubKey
// interface. It represents the 33-byte compressed public key format.
type PubKey []byte

// Address returns the address of the ECDSA public key.
// The function will panic if the public key is invalid.
func (key PubKey) Address() tmcrypto.Address {
	pubk, err := ethcrypto.DecompressPubkey(key)
	if err != nil {
		panic(err)
	}

	return tmcrypto.Address(ethcrypto.PubkeyToAddress(*pubk).Bytes())
}

// Bytes returns the raw bytes of the ECDSA public key.
// The function panics if the key cannot be marshaled to bytes.
func (key PubKey) Bytes() []byte {
	bz, err := CryptoCodec.Amino.MarshalBinaryBare(key)
	if err != nil {
		panic(err)
	}
	return bz
}

// VerifyBytes verifies that the ECDSA public key created a given signature over
// the provided message. It will calculate the Keccak256 hash of the message
// prior to verification.
func (key PubKey) VerifySignature(msg []byte, sig []byte) bool {
	if len(sig) == 65 {
		// remove recovery ID if contained in the signature
		sig = sig[:len(sig)-1]
	}

	// the signature needs to be in [R || S] format when provided to VerifySignature
	return secp256k1.VerifySignature(key, ethcrypto.Keccak256Hash(msg).Bytes(), sig)
}

// Equals returns true if two ECDSA public keys are equal and false otherwise.
func (key PubKey) Equals(other tmcrypto.PubKey) bool {
	if other, ok := other.(PubKey); ok {
		return bytes.Equal(key.Bytes(), other.Bytes())
	}

	return false
}

func (key PubKey) Type() string {
	return KeyType
}
