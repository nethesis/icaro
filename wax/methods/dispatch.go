package methods

import (
	"net/http"
        "github.com/gin-gonic/gin"
         _ "github.com/jinzhu/gorm/dialects/mysql"
        "sun-api/database"
        "sun-api/models"
)

func Reply(c *gin.Context, httpCode int, message string) {
	c.String(httpCode, "Reply-Message: %s", message)
}

func isHotspotUnit(hotspotUnitMac string) (bool){
        var unit models.Unit

        db := database.Database()
        db.Where("mac_address = ?", hotspotUnitMac).First(&unit)
        db.Close()

        return (unit.Id != 0)
}


func isHotspot(hotspotName string) (bool){
	var hotspot models.Hotspot

        db := database.Database()
        db.Where("name = ?", hotspotName).First(&hotspot)
        db.Close()

        return (hotspot.Id != 0)
}

func Dispatch(c *gin.Context) {
        stage := c.Query("stage")
        hotspotName := c.Query("nasid")
        hotspotUnitMac := c.Query("ap")

	if (stage == "") {
		c.String(http.StatusBadRequest, "No stage provided")
		return
	}

	if (!isHotspot(hotspotName)) {
		Reply(c, http.StatusNotFound, "Hotspot not found")
		return
	}

	if (!isHotspotUnit(hotspotUnitMac)) {
		Reply(c, http.StatusNotFound, "Hotspot unit not found")
		return
	}


	//TODO: should we verify uamsecret?

	switch stage {
	case "login":
		hotspotId := c.Query("nasid")
		user := c.Query("user")
		password :=  c.Query("chap_pass")
		challenge := c.Query("chap_chal")
		Login(c, hotspotId, user, password, challenge)


	case "counters":
		parameters := c.Request.URL.Query()
		delete(parameters, "stage")

		// Pass all parameters except 'stage' key
		Counters(c, parameters)

	case "register":
		c.String(http.StatusNotImplemented, "Not implemented: %s", stage)
	default:
		c.String(http.StatusNotFound, "Invalid stage: '%s'", stage)
	}

}
