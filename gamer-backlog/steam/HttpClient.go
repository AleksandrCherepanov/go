package steam

import (
	"strings"

	httpClient "alexanderlab.club/game_priority/http_client"
)

const STORE_API = "https://store.steampowered.com"



func GetUserGame(userId string) ([]byte, *error) {
	body := strings.NewReader("")
	url := STEAM_API + "/IPlayerService/GetOwnedGames/v0001/?key=" + API_KEY + "&steamid=" + userId + "&format=" + FORMAT
	response, err := httpClient.Request("GET", url, body)
	if err != nil {
		return make([]byte, 0), err
	}

	return response, nil
}

func GetUserWishlist(userId string) ([]byte, *error) {
	body := strings.NewReader("")
	url := STORE_API + "/wishlist/profiles/" + userId + "/wishlistdata"
	response, err := httpClient.Request("GET", url, body)
	if err != nil {
		return make([]byte, 0), err
	}

	return response, nil
}

func GetGameDetails(gameId string) ([]byte, *error) {
	body := strings.NewReader("")
	url := STORE_API + "/api/appdetails?appids=" + gameId
	response, err := httpClient.Request("GET", url, body)
	if err != nil {
		return make([]byte, 0), err
	}

	return response, nil
}
