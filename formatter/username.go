package formatter

import "encoding/json"

// Note: this is resued logic from mochi-formatter service
// If you need any more platform or logic, please refer to it and add here
var (
	usernamePrefix = map[string]string{
		PlatformDiscord:  "dsc:",
		PlatformTelegram: "tg:",
		PlatformApp:      "app:",
		PlatformMochi:    "mochi:",
		PlatformFacebook: "fb:",
	}

	usernameEmoji = map[string]string{
		PlatformDiscord:  "👾",
		PlatformTelegram: "🔹",
		PlatformApp:      "🔌",
		PlatformMochi:    "🍡",
	}
)

func Account(platform string, profileA, profileB MochiProfile) (*MochiProfileAccount, *MochiProfileAccount) {
	accountA := map[string]interface{}{
		"telegram": profileA.Telegram,
		"discord":  profileA.Discord,
		"facebook": profileA.Facebook,
	}

	accountB := map[string]interface{}{
		"telegram": profileB.Telegram,
		"discord":  profileB.Discord,
		"facebook": profileB.Facebook,
	}

	fallBackOrder := make([]string, 0)
	switch platform {
	case PlatformWeb:
		fallBackOrder = []string{PlatformApp, PlatformDiscord, PlatformTelegram, PlatformMochi, PlatformFacebook}
	case PlatformDiscord:
		fallBackOrder = []string{PlatformApp, PlatformDiscord, PlatformTelegram, PlatformMochi, PlatformFacebook}
	case PlatformTelegram:
		fallBackOrder = []string{PlatformApp, PlatformTelegram, PlatformDiscord, PlatformMochi, PlatformFacebook}
	default:
		fallBackOrder = []string{}
	}

	var userA, userB *MochiProfileAccount
	for _, p := range fallBackOrder {
		tmpUserA := accountA[p]
		tmpUserB := accountB[p]

		if userA == nil && tmpUserA != nil {
			byteTempUserA, err := json.Marshal(tmpUserA)
			if err != nil {
				return nil, nil
			}

			var parseTempUserA *MochiProfileAccount
			err = json.Unmarshal(byteTempUserA, &parseTempUserA)
			if err != nil {
				return nil, nil
			}

			if parseTempUserA != nil {
				parseTempUserA.Platform = p
				userA = parseTempUserA
			}
		}

		if userB == nil && tmpUserB != nil {
			byteTempUserB, err := json.Marshal(tmpUserB)
			if err != nil {
				return nil, nil
			}

			var parseTempUserB *MochiProfileAccount
			err = json.Unmarshal(byteTempUserB, &parseTempUserB)
			if err != nil {
				return nil, nil
			}

			if parseTempUserB != nil {
				parseTempUserB.Platform = p
				userB = parseTempUserB
			}
		}

		if userA != nil && userB != nil {
			return userA, userB
		}
	}
	return nil, nil
}

func prefix(platform, accountPlatform string) string {
	if platform == accountPlatform {
		return "@"
	}

	return usernameEmoji[accountPlatform] + usernamePrefix[accountPlatform]
}

func FallbackUsernamePrefix(aProfileA *MochiProfileAccount, profileA *MochiProfile, platform string) (string, string) {
	if aProfileA == nil || aProfileA.PlatformMetadata["username"] == nil {
		return profileA.ProfileId, "@"
	}

	return aProfileA.PlatformMetadata["username"].(string), prefix(platform, aProfileA.Platform)
}
