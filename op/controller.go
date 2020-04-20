package op

import (
	"../common"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func CrontrollerFirewallRichRuleUpdate(ctx *gin.Context) {
	var ip = ctx.ClientIP()

	//var port = Common.Config.GetString("http.port")
	var hole = Common.Config.GetString("firewall.hole")
	var id = ctx.Param("id")
	var dataFileName = "./data/id_" + id + ".txt"

	var tpl = FirewallRichRuleMap{
		"accept": {Operate: "accept"},
		"port":   {Operate: "port", Value: []FirewallRichRuleValue{{Property: "port", Value: hole}, {Property: "protocol", Value: "tcp"}}},
		"rule":   {Operate: "rule", Value: []FirewallRichRuleValue{{Property: "family", Value: "ipv4"}}},
		"source": {Operate: "source", Value: []FirewallRichRuleValue{{Property: "address", Value: ip}}},
	}

	oldData, err := ioutil.ReadFile(dataFileName)
	if nil == err {
		var oldTpl FirewallRichRuleMap
		Common.Json.Unmarshal(oldData, &oldTpl)
		_, ok := oldTpl["source"]
		if ok {
			for _, v := range oldTpl["source"].Value {
				if "address" == v.Property {
					if v.Value == ip {
						ctx.JSON(http.StatusOK, Common.ResultData{
							Code: http.StatusAlreadyReported,
							Err:  "",
							Msg:  "ip not changed",
							Data: "",
						})
						return
					}
				}
			}
		}

		FirewallRemoveRichRules(oldTpl)
	}

	err = FirewallAddRichRules(tpl)

	if nil == err {
		ioutil.WriteFile(dataFileName, []byte(Common.FastJsonMarshal(tpl)), 0777)
		ctx.JSON(http.StatusOK, Common.ResultData{
			Code: http.StatusOK,
			Err:  "",
			Msg:  "success",
			Data: "",
		})
	} else {
		ctx.JSON(http.StatusOK, Common.ResultData{
			Code: http.StatusInternalServerError,
			Err:  err.Error(),
			Msg:  "op failed",
			Data: "",
		})
	}
	return

}
