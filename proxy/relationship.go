package proxy

type Relationship struct {
	User           string
	Domain         string
	OriginalDomain string
	Conn           Conn
}
