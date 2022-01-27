package stream

var cipherMethods = map[string]cipherInfo{
	"aes-128-cfb":   newCipherInfo(16, 16, newAesCfbReader, newAesCfbWriter),
	"aes-192-cfb":   newCipherInfo(24, 16, newAesCfbReader, newAesCfbWriter),
	"aes-256-cfb":   newCipherInfo(32, 16, newAesCfbReader, newAesCfbWriter),
	"aes-128-ctr":   newCipherInfo(16, 16, newAesCtrReader, newAesCtrWriter),
	"aes-192-ctr":   newCipherInfo(24, 16, newAesCtrReader, newAesCtrWriter),
	"aes-256-ctr":   newCipherInfo(32, 16, newAesCtrReader, newAesCtrWriter),
	"des-cfb":       newCipherInfo(8, 8, newDESCFBReader, newDESCFBWriter),
	"rc4-md5":       newCipherInfo(16, 16, newRC4MD5Reader, newRC4MD5Writer),
	"rc4-md5-6":     newCipherInfo(16, 6, newRC4MD5Reader, newRC4MD5Writer),
	"chacha20":      newCipherInfo(32, 8, newChaCha20Reader, newChaCha20Writer),
	"chacha20-ietf": newCipherInfo(32, 12, newChaCha20Reader, newChaCha20Writer),
	"cast5-cfb":     newCipherInfo(16, 8, newCast5Reader, newCast5Writer),
	"bf-cfb":        newCipherInfo(16, 8, newBlowfishReader, newBlowfishWriter),
}
