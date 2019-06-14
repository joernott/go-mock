package rules

import (
	"encoding/json"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

type Ruleset struct {
	Methods map[string]MethodRules `json:"methods"`
}

type MethodRules []MethodRule
type MethodRule struct {
	Path            string     `json:"path"`
	Query           string     `json:"query"`
	ResponseCode    int        `json:"response_code"`
	ResponseHeaders HeaderList `json:"response_headers"`
	ResponseBody    string     `json:"response_body"`
}

type HeaderList map[string]string

func LoadRules(RuleFile string) (*Ruleset, error) {
	file, err := ioutil.ReadFile(RuleFile)
	if err != nil {
		log.WithFields(log.Fields{
			"File": RuleFile,
			"func": LoadRules,
		}).Error(err)
		return nil, err
	}
	data := new(Ruleset)

	err = json.Unmarshal([]byte(file), data)
	if err != nil {
		return nil, err
	}
	spew.Dump(data)
	return data, err
}
