package dac_policy

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (r *SensitiveDataRule) PrintYamlHeader(silent bool) *SensitiveDataRule {
	if silent {
		return r
	}
	fmt.Printf("# begin of dac SensitiveDataRule v0.9\n")
	return r
}

func (r *SensitiveDataRule) PrintYaml(silent bool) *SensitiveDataRule {
	logrus.Debug(r)
	if silent {
		return r
	}
	slice := []SensitiveDataRule{*r}
	_yaml, err := yaml.Marshal(slice)
	if err == nil {
		fmt.Print(string(_yaml))
	} else {
		logrus.Errorf("Failed to marshal SensitiveDataRule: %v", err)
	}
	return r
}

func (r *SensitiveDataRule) PrintYamlFooter(silent bool, count int) *SensitiveDataRule {
	if silent {
		return r
	}
	fmt.Printf("# end of dac SensitiveDataRule v0.9, total elements: %d\n", count)
	logrus.Infof("TotalElements: %v", count)
	return r
}
