//Copyright (c) 2014 Tim Shannon
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
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

//Tries to adhere to the xdg base directory specification http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html

func userLocations() []string {
	location := os.Getenv("XDG_CONFIG_HOME")
	if location != "" {
		return []string{location}
	}
	usr, err := user.Current()
	if err != nil {
		return []string{}
	}

	return []string{filepath.Join(usr.HomeDir, ".config")}
}

func systemLocations() []string {
	defaults := []string{"/usr/local/etc/xdg", "/usr/local/etc"}
	envLocations := os.Getenv("XDG_CONFIG_DIRS")
	if envLocations == "" {
		return defaults
	}
	locations := strings.Split(envLocations, ":")

	for d := range defaults {
		found := false
		for l := range locations {
			locations[l] = filepath.Clean(locations[l])
			if locations[l] == defaults[d] {
				found = true
				break
			}
		}
		if !found {
			locations = append(locations, defaults[d])
		}
	}

	return locations
}
