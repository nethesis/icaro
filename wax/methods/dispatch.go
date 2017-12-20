package methods

import (
	"net/http"
        "github.com/gin-gonic/gin"
)

func Dispatch(c *gin.Context) {
        stage := c.Query("stage")

	switch stage {
	case "login":
		hotspotId := c.Query("nasid")
		user := c.Query("user")
		password :=  c.Query("chap_pass")
		challenge := c.Query("chap_chal")
		Login(hotspotId, user, password, challenge)
		c.String(http.StatusOK, stage)

	case "counters":
		/*
		/hotspot/hs-script.php?stage=counters&status=stop&user=aaa&ap=00-0D-B9-41-7C-F8&mac=84-3A-4B-11-44-D4&ip=10.1.0.3&sessionid=151324185800000001&nasid=hs-test&duration=702&bytes_down=47487&pkts_down=135&bytes_up=5292&pkts_up=65&md=F2D7D9B3184E2890140C9B7FE28CC0FB
/hotspot/hs-script.php?stage=counters&status=up&ap=00-0D-B9-41-7C-F8&mac=00-00-00-00-00-00&nasid=hs-test&md=F998859655817FE494BB9CEC3133BD2A
/hotspot/hs-script.php?stage=counters&status=update&user=aaa&ap=00-0D-B9-41-7C-F8&mac=84-3A-4B-11-44-D4&ip=10.1.0.3&sessionid=151318184500000001&nasid=hs-test&duration=10205&bytes_down=656025&pkts_down=1655&bytes_up=164192&pkts_up=1794&md=F614C20150DCD5A8D867A977DEFF154C
*/
		// status := c.Query("status")

		c.String(http.StatusOK, stage)

	case "register":
		c.String(http.StatusNotImplemented, "Not implemented: %s", stage)

	case "":
		c.String(http.StatusBadRequest, "No stage provided")

	default:
		c.String(http.StatusNotFound, "Invalid stage: '%s'", stage)
	}

}
