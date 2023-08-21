package formatter

const (
	PlatformDiscord  = "discord"
	PlatformTelegram = "telegram"
	PlatformApp      = "app"
	PlatformMochi    = "mochi"
	PlatformWeb      = "web"
)

type MochiProfile struct {
	ProfileId string
	Discord   *MochiProfileDiscord
	Email     *MochiProfileEmail
	Telegram  *MochiProfileTelegram
}

type MochiProfileAccount struct {
	Id               string
	Platform         string
	PlatformMetadata map[string]interface{}
}

type MochiProfileDiscord struct {
	Id               string
	PlatformMetadata map[string]interface{}
}

type MochiProfileEmail struct {
	Id               string
	PlatformMetadata map[string]interface{}
}

type MochiProfileTelegram struct {
	Id               string
	PlatformMetadata map[string]interface{}
}
