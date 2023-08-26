package config

type AppConfig struct {
	AppSecret     []byte
	AccessTkn     []byte
	HubToken      []byte
	WTPhoneID     []byte
	InDevelopment bool
}
