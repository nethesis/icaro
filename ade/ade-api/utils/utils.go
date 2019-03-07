package utils

import (
	"github.com/nethesis/icaro/ade/ade-api/models"
	"github.com/nethesis/icaro/sun/sun-api/database"
	wax_utils "github.com/nethesis/icaro/wax/utils"
)

func GetHotspotPrefs(hotspotId int) map[string]string {
	hotspotPerfs := []string{"captive_2_title",
		"captive_3_logo",
		"captive_7_background",
		"marketing_9_threshold",
		"marketing_10_first_url",
		"marketing_11_second_url",
		"marketing_12_third_url",
	}

	prefsMap := make(map[string]string)

	prefs := wax_utils.GetHotspotPreferencesByKeys(hotspotId, hotspotPerfs)

	for i := 0; i < len(prefs); i++ {
		prefsMap[prefs[i].Key] = prefs[i].Value
	}

	return prefsMap
}

func GetAdeTokenFromToken(token string) models.AdeToken {
	var adeToken models.AdeToken

	db := database.Instance()
	db.Where("token = ?", token).First(&adeToken)

	return adeToken
}
