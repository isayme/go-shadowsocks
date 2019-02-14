package cipher

var cipherMethods = map[string]*cipherInfo{
	"aes-128-cfb": &cipherInfo{16, 16, newAESCFBEncryptStream, newAESCFBDecryptStream},
	"aes-256-cfb": &cipherInfo{32, 16, newAESCFBEncryptStream, newAESCFBDecryptStream},

	"rc4-md5": &cipherInfo{16, 16, newRC4MD5Stream, newRC4MD5Stream},
}
