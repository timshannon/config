//Copyright (c) 2013 Tim Shannon
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

//autoWrite will automatically write out non-existant entries
// that are defaulted
var autoWrite bool = false
var ErrNotFound = errors.New("Value not found")

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

//LoadOrCreate automatically creates the passed in config file if it
// doesn't already exist, then load it.  If file gets created, then all
// values loaded afterwards will be "defaulted" and will be written to this
// new file
func LoadOrCreate(filename string) (*Cfg, error) {
	var err error
	cfg := &Cfg{FileName: filename}
	if !cfg.FileExists() {
		autoWrite = true
		cfg.values = make(map[string]interface{})
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
func (c *Cfg) Value(name string, defaultValue interface{}) interface{} {
	value, ok := c.values[name]
	if !ok {
		if autoWrite {
			c.SetValue(name, defaultValue)
			c.Write()
		}
		return defaultValue
	}

	return value
}

//ValueToType allows you to pass in a struct as the result
// for which you want to load the config entry into
// Marshalls the JSON data directly into your passed in type
func (c *Cfg) ValueToType(name string, result interface{}) error {
	value, ok := c.values[name]
	if !ok {
		return ErrNotFound
	}
	//marshall value
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, result)
	return err
}

//Int retrieves an integer config value with the given name
// if a value with the given name is not found the default is returned
func (c *Cfg) Int(name string, defaultValue int) int {
	value, ok := c.values[name].(float64)
	if !ok {
		if autoWrite {
			c.SetValue(name, defaultValue)
			c.Write()
		}
		return defaultValue
	}
	return int(value)
}

//String retrieves a string config value with the given name
// if a value with the given name is not found the default is returned
func (c *Cfg) String(name string, defaultValue string) string {
	value, ok := c.values[name].(string)
	if !ok {
		if autoWrite {
			c.SetValue(name, defaultValue)
			c.Write()
		}
		return defaultValue
	}
	return value
}

//Bool retrieves a bool config value with the given name
// if a value with the given name is not found the default is returned
func (c *Cfg) Bool(name string, defaultValue bool) bool {
	value, ok := c.values[name].(bool)
	if !ok {
		if autoWrite {
			c.SetValue(name, defaultValue)
			c.Write()
		}
		return defaultValue
	}
	return value
}

//Float retrieves a float config value with the given name
// if a value with the given name is not found the default is returned
func (c *Cfg) Float(name string, defaultValue float32) float32 {
	value, ok := c.values[name].(float64)

	if !ok {
		if autoWrite {
			c.SetValue(name, defaultValue)
			c.Write()
		}

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
