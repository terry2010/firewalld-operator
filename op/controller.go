package op

import (
	"github.com/gin-gonic/gin"
	"log"
)

func CrontrollerFirewallRichRuleAdd(ctx *gin.Context) {
	var ip = ctx.ClientIP()
	var port = "3389"
	var tpl = FirewallRichRuleMap{
		"accept": {Operate: "accept"},
		"port":   {Operate: "port", Value: []FirewallRichRuleValue{{Property: "port", Value: port}, {Property: "protocol", Value: "tcp"}}},
		"rule":   {Operate: "rule", Value: []FirewallRichRuleValue{{Property: "family", Value: "ipv4"}}},
		"source": {Operate: "source", Value: []FirewallRichRuleValue{{Property: "address", Value: ip}}},
	}

	err := FirewallAddRichRules(tpl)

	log.Println(err)

}
