package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func (c *config) FromYaml(in []byte) error {
	err := yaml.UnmarshalStrict(in, c)

	if err != nil {
		return err
	}

	return nil
}

func (c *config) FromYamlFile(path string) error {
	in, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	err = c.FromYaml(in)

	if err != nil {
		return err
	}

	return nil
}
