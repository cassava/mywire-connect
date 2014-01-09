// Copyright (c) 2013, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	// confUser and confPass are the environment variables that are read for
	// the mywire username and password, respectively.
	confUser = "MYWIRE_USER"
	confPass = "MYWIRE_PASS"
)

type Conf struct {
	// User is the username for logging into mywire.
	User string
	// Pass is the password for logging into mywire.
	Pass string

	// PingURL is the URL to use for testing internet connectivity, we are
	// using Google because they are usually pretty quick in their responses.
	PingURL string

	// Wait WaitTime milliseconds if mywire tells us to wait; make sure that
	// this is always greater than 0, otherwise we will send requests as fast
	// as we can, and that would not be very good.
	WaitTime uint

	// NotifyLevel specifies if and what notifications should be used.
	// The possible values are: 0 = none, 1 = critical, 2 = normal, or 3 = all
	// Anything else is considered to be all.
	NotifyLevel
}

func NewDefaultConf() Conf {
	return Conf{
		PingURL:     "http://www.google.com",
		WaitTime:    2000,
		NotifyLevel: NotifyAll,
	}
}

func Version() {
	fmt.Println("mywire-connect version", version)
}

func Usage() {
	conf := NewDefaultConf()
	fmt.Printf(`Usage: mywire-connect [-version|-config CONFIG]

mywire-connect tries to login to mywire using the credentials you supply,
either from the environment or from the file given by CONFIG. If both
are supplied, then the configuration file has precedence.

The environmental variables you can set are:

	%s
	%s

If using the configuration file, then you may specify a few more options.
The configuration file is in the JSON format, with the following defaults:

	{
		"User": "",
		"Pass": "",
		"PingURL": %q,
		"WaitTime": %d,
		"NotifyLevel": %d
	}

You should be set. Do note that notifications only work on Linux or BSD,
wherever D-Bus and freedesktop.org notifications are present.
`, confUser, confPass, conf.PingURL, conf.WaitTime, conf.NotifyLevel)
}

// LoadConfiguration loads the configuration from the environment or from
// a specified configuration file.
//
// LoadConfiguration has the authority to exit at will.
func LoadConfiguration() Conf {
	path := flag.String("config", "", "use the given configuration file")
	ver := flag.Bool("version", false, "print the current version")
	flag.Usage = Usage
	flag.Parse()

	if *ver {
		Version()
		os.Exit(0)
	}

	conf := NewDefaultConf()
	readAuthFromEnvironment(&conf)
	if *path != "" {
		err := readConfFromFile(*path, &conf)
		if err != nil {
			fmt.Println("Fatal error: ", err)
			Usage()
			os.Exit(1)
		}
	}

	if conf.User == "" || conf.Pass == "" {
		Usage()
		os.Exit(1)
	}

	return conf
}

func readAuthFromEnvironment(conf *Conf) {
	conf.User = os.Getenv(confUser)
	conf.Pass = os.Getenv(confPass)
}

func readConfFromFile(path string, conf *Conf) error {
	blob, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(blob, conf)
}
