package fipslint

/*
 eed to implement these configirations
Digest algorithms: MD2, MD4, MD5, RIPEMD160, and SHA-0
Message Authentication Code (MAC) algorithms: HMAC algorithms that involve any of the unsupported digest algorithms
Symmetric algorithms: Blowfish, CAST, RC2, RC4, and single DES (3DES is still allowed)
Asymmetric algorithms: DSA with a key length of less than 1024 bits; RSA with a key length of less than 1024 bits
*/
var (
	excludedFunctions = []function{
		{"crypto/md5", "New"},
		{"crypto/md4", "New"},
	}
)
