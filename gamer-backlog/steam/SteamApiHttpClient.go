package steam

import (
	"strings"

	httpClient "alexanderlab.club/game_priority/http_client"
)

const url = "https://api.steampowered.com"
const format = "json"
const apiKey = "*"

type SteamApi struct {
	key    string
	url    string
	format string
	http_client
}

func NewSteamApi() *SteamApi {
	steamApi := &SteamApi{}
	steamApi.format = format
	steamApi.key = apiKey
	steamApi.url = url

	return steamApi
}

func (steamApi *SteamApi) GetUserGames(userId string) ([]byte, *error) {
	body := strings.NewReader("")
	steamApi.url = steamApi.url + "/IPlayerService/GetOwnedGames/v0001/?key=" + steamApi.key +  

	url := STEAM_API + ?key=" + API_KEY + "&steamid=" + userId + "&format=" + FORMAT
	response, err := httpClient.Request("GET", url, body)
	if err != nil {
		return make([]byte, 0), err
	}

	return response, nil
}
