package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"crypto/elliptic"

	"golang.org/x/crypto/ripemd160"
)

const (
	version            = byte(0x00)
	addressChecksumLen = 4
)

// newKeyPair creates a new cryptographic key pair
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privSt, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, nil
	}

	pubKey := pubKeyToByte(privSt.PublicKey)
	return *privSt, pubKey
}

// pubKeyToByte converts the ecdsa.PublicKey to a concatenation of its coordinates in bytes
func pubKeyToByte(pubkey ecdsa.PublicKey) []byte {
	var add []byte
	add = append(add, pubkey.X.Bytes()...)
	add = append(add, pubkey.Y.Bytes()...)
	return add
}

// GetAddress returns address
// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
func GetAddress(pubKeyBytes []byte) []byte {
	var versionedKey []byte
	hashedKey := HashPubKey(pubKeyBytes)
	versionedKey = append(versionedKey, []byte{version}...)
	versionedKey = append(versionedKey, hashedKey...)
	checksum := checksum(versionedKey)
	versionedKey = append(versionedKey, checksum...)
	return Base58Encode(versionedKey)
}

// GetStringAddress returns address as string
func GetStringAddress(pubKeyBytes []byte) string {
	return string(pubKeyBytes[:])
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	shaPubKey := sha256.Sum256(pubKey)
	rip := ripemd160.New()
	rip.Write(shaPubKey[:])
	hashBytes := rip.Sum(nil)
	return hashBytes
}

// GetPubKeyHashFromAddress returns the hash of the public key
// discarding the version and the checksum
func GetPubKeyHashFromAddress(address string) []byte {
	decodedAddr := Base58Decode([]byte(address))
	return decodedAddr[1 : len(decodedAddr)-4]
}

// ValidateAddress check if an address is valid
func ValidateAddress(address string) bool {
	// TODO(student)
	// Validate a address by decoding it, extracting the
	// checksum, re-computing it using the "checksum" function
	// and comparing both.
	decodedAddr := Base58Decode([]byte(address))
	pubKey := decodedAddr[0 : len(decodedAddr)-4]
	oldchksum := decodedAddr[len(decodedAddr)-4:]
	newChecksum := checksum(pubKey)
	return bytes.Equal(oldchksum, newChecksum)
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[0:4]
}

func encodeKeyPair(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	return encodePrivateKey(privateKey), encodePublicKey(publicKey)
}

func encodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	return string(pemEncoded)
}

func encodePublicKey(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncodedPub)
}

func decodeKeyPair(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	return decodePrivateKey(pemEncoded), decodePublicKey(pemEncodedPub)
}

func decodePrivateKey(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	privateKey, _ := x509.ParseECPrivateKey(block.Bytes)

	return privateKey
}

func decodePublicKey(pemEncodedPub string) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	genericPubKey, _ := x509.ParsePKIXPublicKey(blockPub.Bytes)
	publicKey := genericPubKey.(*ecdsa.PublicKey) // cast to ecdsa

	return publicKey
}
