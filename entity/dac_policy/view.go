package dac_policy

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (p *Policy) PrintYamlHeader(silent bool) *Policy {
	if silent {
		return p
	}
	fmt.Printf("# begin of dac policy v0.9\n")
	return p
}

func (p *Policy) PrintYaml(silent bool) *Policy {
	logrus.Debug(p)
	if silent {
		return p
	}
	slice := []Policy{*p}
	_yaml, err := yaml.Marshal(slice)
	if err == nil {
		fmt.Print(string(_yaml))
	} else {
		logrus.Errorf("Failed to marshal Policy: %v", err)
	}
	return p
}

func (p *Policy) PrintYamlFooter(silent bool, count int) *Policy {
	if silent {
		return p
	}
	fmt.Printf("# end of dac policy v0.9, total elements: %d\n", count)
	logrus.Infof("TotalElements: %v", count)
	return p
}
