package dtos

type PingDto struct {
	Ping string
}

func (p *PingDto) IsPong() bool {
	return p.Ping == "pong"
}
