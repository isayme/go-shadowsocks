package aead

// https://shadowsocks.org/en/wiki/AEAD-Ciphers.html
var cipherMethods = map[string]cipherInfo{
	"aes-128-gcm":            newCipherInfo(16, newAesGcmCipher, newAesGcmCipher),
	"aes-192-gcm":            newCipherInfo(24, newAesGcmCipher, newAesGcmCipher),
	"aes-256-gcm":            newCipherInfo(32, newAesGcmCipher, newAesGcmCipher),
	"chacha20-ietf-poly1305": newCipherInfo(32, newChacha20Poly1305Cipher, newChacha20Poly1305Cipher),
}
