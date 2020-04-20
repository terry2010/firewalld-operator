package op

import (
	"../common"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func init() {

}

func FirewallAddRichRules(theRule FirewallRichRuleMap) error {
	_cmd := `firewall-cmd --permanent --add-rich-rule="` + FirewallCreatRichRule(theRule) + `"`
	fileName, err := writeCMD(_cmd)

	if nil == err {
		cmd := exec.Command("bash", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}
	deleteCMD(fileName)
	return err
}
func FirewallRemoveRichRules(theRule FirewallRichRuleMap) error {
	_cmd := `firewall-cmd --permanent --remove-rich-rule="` + FirewallCreatRichRule(theRule) + `"`
	fileName, err := writeCMD(_cmd)

	if nil == err {
		cmd := exec.Command("bash", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}
	deleteCMD(fileName)
	return err
}

func FirewallRemoveRichRules2(theDelRule FirewallRichRuleMap) error {
	//content := `rule family="ipv4" source address="192.168.1.196" port port="80" protocol="tcp" accept`
	//firewall-cmd --add-rich-rule='rule family="ipv4" port port="80" protocol="tcp" source address="192.168.1.196" accept '
	var rules, err = FirewallListRichRules()
	if nil == err {
		if len(rules) < 1 {
			return nil
		}
		for _, v := range rules {
			if strings.EqualFold(Common.FastJsonMarshal(v), Common.FastJsonMarshal(theDelRule)) {
				_cmd := `firewall-cmd --remove-rich-rule='` + FirewallCreatRichRule(theDelRule) + `'`
				log.Println(_cmd)
				output, err := exec.Command(_cmd).Output()
				log.Println(err, output)
				return err
			}
		}

	} else {
		return err
	}
	return nil
}

func ResetWhiteIPPort(id string) {
	port := Common.Config.GetString("firewall.port")
	_cmd := "firewall-cmd  --list-rich-rules"
	output, err := exec.Command(_cmd).Output()
	log.Println(port, output, err)
}

func FirewallListRichRules() (rules []FirewallRichRuleMap, err error) {
	_cmd := "firewall-cmd  --list-rich-rules"
	var output []byte
	output, err = exec.Command(_cmd).Output()

	if nil == err {
		tmpRules := strings.Split(string(output), "\n")
		if len(tmpRules) > 0 {
			for _, v := range tmpRules {
				if len(v) > 0 {
					rules = append(rules, FirewallParseRichRule(v))
				}
			}
		}
	}

	return
}

func FirewallParseRichRule(content string) (ruleList FirewallRichRuleMap) {

	//content := `rule family="ipv4" source address="192.168.1.196" port port="80" protocol="tcp" accept`

	//var ruleList = make(FirewallRichRuleMap)
	var newRule = true
	parts := strings.Split(content, " ")
	if len(parts) > 0 {
		var tmpOprator = ""
		var tmpRule FirewallRichRule

		for _, v := range parts {

			if false == strings.Contains(v, "=") {
				newRule = true
			}
			if true == newRule {

				ruleList[tmpOprator] = tmpRule

				tmpOprator = v
				newRule = false
				tmpRule = FirewallRichRule{Operate: v}

			} else {
				_split := strings.Split(v, "=")
				log.Println(_split)
				if len(_split[1]) > 0 && (_split[1][0] == '"' || _split[1][0] == '\'') {
					_split[1] = _split[1][1:]
				}
				if len(_split[1]) > 0 && (_split[1][len(_split[1])-1] == '"' || _split[1][len(_split[1])-1] == '\'') {
					_split[1] = _split[1][:len(_split[1])-1]
				}
				tmpRule.Value = append(tmpRule.Value, struct {
					Property string
					Value    string
				}{Property: _split[0], Value: _split[1]})
			}
		}

		if _, ok := ruleList[tmpOprator]; false == ok {
			ruleList[tmpOprator] = FirewallRichRule{Operate: tmpOprator}
		}
	}
	return
}

func FirewallCreatRichRule(ruleList FirewallRichRuleMap) (content string) {

	var richRule string
	var richRuleParts []string
	for k, v := range ruleList {
		var tmpRule string
		if len(k) > 0 {
			tmpRule = tmpRule + v.Operate + " "
			if len(v.Value) > 0 {
				for _, v1 := range v.Value {
					//tmpRule = tmpRule + v1.Property + "=\"" + v1.Value + "\" "
					tmpRule = tmpRule + v1.Property + "=" + v1.Value + " "
				}
			}
			richRuleParts = append(richRuleParts, tmpRule)
			//richRule = richRule + tmpRule
			//fmt.Println(richRuleParts)
		}
	}

	for k, v := range richRuleParts {
		if strings.HasPrefix(v, "rule") {
			richRuleParts[k] = ""
			richRule = richRule + v
		}
	}
	for k, v := range richRuleParts {
		if strings.HasPrefix(v, "accept") ||
			strings.HasPrefix(v, "reject") {

		} else {
			richRuleParts[k] = ""
			richRule = richRule + v
		}

	}

	for k, v := range richRuleParts {
		if strings.HasPrefix(v, "accept") ||
			strings.HasPrefix(v, "reject") {
			richRuleParts[k] = ""
			richRule = richRule + v
		}
	}
	return richRule
}

func writeCMD(_cmd string) (fileName string, err error) {
	fileName = "./data/cmd_" + time.Now().Format("200612_150405_") + strconv.FormatInt(rand.Int63n(10000), 10) + ".sh"
	_cmd = _cmd + "\n\nfirewall-cmd --reload"
	_cmd = _cmd + "\n\nfirewall-cmd --list-all"
	log.Println(_cmd)
	err = ioutil.WriteFile(fileName, []byte(_cmd), 0777)
	return
}
func deleteCMD(fileName string) error {
	return os.Remove(fileName)
}
