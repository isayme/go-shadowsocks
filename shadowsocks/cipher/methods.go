package cipher

var cipherMethods = map[string]*cipherInfo{
	"aes-128-cfb": &cipherInfo{16, 16, newAESCFBEncryptStream, newAESCFBDecryptStream},
	"aes-256-cfb": &cipherInfo{32, 16, newAESCFBEncryptStream, newAESCFBDecryptStream},
}
