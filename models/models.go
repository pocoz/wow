package models

type ConnectCfg struct {
	NetworkType string
	Address     string
	Port        int
}

type Message struct {
	Body string `json:"body"`
}
