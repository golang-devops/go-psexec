package dtos

type EncryptedJsonContainer struct {
	EncryptedJson []byte
}

type GetTokenRequestDto struct {
	ClientPubPKIXBytes []byte
}

type GenTokenResponseDto struct {
	EncryptedSessionToken []byte //Encrypted with client public-key
	EncryptedMessage      []byte //Encrypted with the session-token, the unencrypted object is of type `GenTokenResponseMessage`
}

type GenTokenResponseMessage struct {
	SessionId                int
	TokenEncryptionSignature []byte
	ServerPubKeyBytes        []byte
}
