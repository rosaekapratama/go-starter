package yaml

import (
	"errors"
	"gopkg.in/yaml.v3"
	"time"
)

var invalidDurationErr = errors.New("invalid duration value")

type Duration struct {
	time.Duration
}

func (d Duration) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(d.String())
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	var v interface{}
	if err := yaml.Unmarshal([]byte(node.Value), &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case int:
		d.Duration = time.Duration(value)
		return nil
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return invalidDurationErr
	}
}
