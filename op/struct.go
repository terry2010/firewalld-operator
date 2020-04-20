package op

var err error

type FirewallRichRule struct {
	Operate string
	Value   []FirewallRichRuleValue
}

type FirewallRichRuleValue struct {
	Property string
	Value    string
}

type FirewallRichRuleMap map[string]FirewallRichRule

//map[:{ []} accept:{accept []} port:{port [{port 80} {protocol tcp}]} rule:{rule [{family ipv4}]} source:{source [{address 192.168.1.196}]}]
////map[:{ []} accept:{accept []} port:{port [{port 80} {protocol tcp}]}
// rule:{rule [{family ipv4}]} source:{source [{address 192.168.1.196}]}]
var FireRichRuleTemplate = FirewallRichRuleMap{
	"accept": {Operate: "accept"},
	"port":   {Operate: "port", Value: []FirewallRichRuleValue{{Property: "port", Value: "80"}, {Property: "protocol", Value: "tcp"}}},
	"rule":   {Operate: "rule", Value: []FirewallRichRuleValue{{Property: "family", Value: "ipv4"}}},
	"source": {Operate: "source", Value: []FirewallRichRuleValue{{Property: "address", Value: "192.168.1.196"}}},
}
