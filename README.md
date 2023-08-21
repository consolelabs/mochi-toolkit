# Mochi Toolkit

## Format number

0. value >= 1000, not show decimal part. <br />Ex: 1000.382 -> 1000
1. 1 <= value < 1000, show only 2 decimals<br />
   Ex: 23.42347 -> 23.42, 100.00000153271831 -> 100.00, 100.012324523 -> 100.01, 100.19231123 -> 100.19
2. value < 1

- show to the first non-zero decimal
- number of zero decimals > 6 then show scientific number<br />
  Ex: 0.00012435 -> 0.0001, 0.0001999999 -> 0.00019, 0.1999999 -> 0.1
  Ex: 0.000000000000000001 -> 1e-18

## Format username

Profile will associate with accounts in many platform: Telegram, Discord, Web, Application,... Logic to render username is choose in current platform first, if not have then get from list fallback platform

0. Platform

- Context platform: where action happens -> ["web", "discord", "telegram"]
- Account platform: platform of all associated account -> ["app", "discord", "telegram", "mochi"]

1. Fallback platform order

- context platform = web. fallbackOrder = ["app", "mochi", "discord", "telegram"]
- context platform = discord. fallbackOrder = ["app", "discord", "telegram", "mochi"]
- context platform = telegram. fallbackOrder = ["app", "telegram", "discord", "mochi"]

2. Render username. Get username from context platform first, if not have then use fallback list<br />
   Ex: profile 1 tip profile 2 on platform = "telegram". profile 1: {telegram - username1Tel, discord - username1Disc, username1Mochi}, profile2: {discord - username2Disc}
   => show username1Tel and username2Disc<br />

3. Prefix:

- discord: "disc:"
- telegram: "tg:"
- app: "app:"
- mochi: "mochi:"

4. Render prefix. Base on platform and username platform

- if platform = username platform -> @
- if platform != username platform -> get from list prefix<br />
  Ex: from Ex 2, platform = "telegram" and profile 1 = username1Tel (on telegram) and profile 2 = username2Disc (on discord)
  => show @username1Tel and disc:username2Disc

5. In step 2 and 4. If cannot get any username from fallback list => show profile_name -> profile_id
