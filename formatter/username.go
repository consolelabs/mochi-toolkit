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
	}
)

func Account(platform string, profileA, profileB MochiProfile) (*MochiProfileAccount, *MochiProfileAccount) {
	accountA := map[string]interface{}{
		"telegram": profileA.Telegram,
		"discord":  profileA.Discord,
	}

	accountB := map[string]interface{}{
		"telegram": profileB.Telegram,
		"discord":  profileB.Discord,
	}

	fallBackOrder := make([]string, 0)
	switch platform {
	case PlatformWeb:
		fallBackOrder = []string{PlatformApp, PlatformDiscord, PlatformTelegram, PlatformMochi}
	case PlatformDiscord:
		fallBackOrder = []string{PlatformApp, PlatformDiscord, PlatformTelegram, PlatformMochi}
	case PlatformTelegram:
		fallBackOrder = []string{PlatformApp, PlatformTelegram, PlatformDiscord, PlatformMochi}
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

func Prefix(platform, accountPlatform string) string {
	if platform == accountPlatform {
		return "@"
	}

	return usernamePrefix[accountPlatform]
}

func FallbackUsernamePrefix(aProfileA *MochiProfileAccount, profileA *MochiProfile, platform string) (string, string) {
	if aProfileA == nil || aProfileA.PlatformMetadata["username"] == nil {
		return profileA.ProfileId, "@"
	}

	return aProfileA.PlatformMetadata["username"].(string), Prefix(platform, aProfileA.Platform)
}
