package shared

type ExecDto struct {
	Executor string
	Exe      string
	Args     []string
}

type GetTokenRequestDto struct {
	ClientPubPKIXBytes []byte
}

type GenTokenResponseDto struct {
	ServerPubPKIXBytes []byte
	SessionToken       string
}
