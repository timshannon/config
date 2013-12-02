package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

//Cfg is a container for reading and writing to a simple
// JSON config file, nothing fancy.  Easy to parse,
// easy to read and edit by humans
type Cfg struct {
	FileName string
	values   map[string]interface{}
}

//LoadCfg loads the config file from the passed in config path
// If the config file cannot be parsed (i.e. isn't valid json) or
// cannot be found an error will be returned
func LoadCfg(filename string) (*Cfg, error) {
	c := &Cfg{FileName: filename}
	err := c.Load()

	return c, err

}

//CreateAndLoad automatically creates the passed in config file if it
// doesn't already exist, then load it.
func LoadAndCreate(filename string) (*Cfg, error) {
	var err error
	cfg := &Cfg{FileName: filename}
	if !cfg.FileExists() {
		err = cfg.Write()
		if err != nil {
			return nil, err
		}
	}

	err = cfg.Load()

	return cfg, err
}

func (c *Cfg) FileExists() bool {
	file, err := os.Open(c.FileName)
	file.Close()

	return err == nil || !os.IsNotExist(err)
}

//Load loads a config file from the passed in location
func (c *Cfg) Load() error {
	if c.FileName == "" {
		err := errors.New("No Filename set for Cfg object")
		return err
	}

	data, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &c.values); err != nil {
		return err
	}
	return nil
}

//Value returns the raw interface value for the config entry with the given name
//  The type is left up to the consumer to determine
func (c *Cfg) Value(name string) interface{} {
	return c.values[name]
}

//Int retrieves an integer config value with the given name
// if a value with the given name is not found 0 is returned
// and an error is logged
func (c *Cfg) Int(name string, defaultValue int) int {
	value, ok := c.values[name].(float64)
	if !ok {
		return defaultValue
	}
	return int(value)
}

//String retrieves a string config value with the given name
// if a value with the given name is not found "" is returned
// and an error is logged
func (c *Cfg) String(name string, defaultValue string) string {
	value, ok := c.values[name].(string)
	if !ok {
		return defaultValue
	}
	return value
}

//Bool retrieves a bool config value with the given name
// if a value with the given name is not found false is returned
// and an error is logged
func (c *Cfg) Bool(name string, defaultValue bool) bool {
	value, ok := c.values[name].(bool)
	if !ok {
		return defaultValue
	}
	return value
}

//Float retrieves a float config value with the given name
// if a value with the given name is not found 0.0 is returned
// and an error is logged
func (c *Cfg) Float(name string, defaultValue float32) float32 {
	value, ok := c.values[name].(float64)
	if !ok {
		return defaultValue
	}
	return float32(value)
}

//SetValue sets a config value.  It is left up to the end user
// to then write out the new values with the .Write() function
func (c *Cfg) SetValue(name string, value interface{}) {
	c.values[name] = value
}

//Write writes the config values to the config's file location
func (c *Cfg) Write() error {
	if c.FileName == "" {
		return errors.New("No FileName set for this config")
	}
	data, err := json.MarshalIndent(c.values, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.FileName, data, 0644)

	return err

}