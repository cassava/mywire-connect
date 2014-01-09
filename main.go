// Copyright (c) 2013, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// mywire-connect is a program to initiate the connection to mywire.
//
// Bonus to those on Linux: mywire-connect can show tray notifications as to
// what it's status is.
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type mywireStatus int

const (
	mywireUnknown mywireStatus = iota
	mywireSuccess
	mywireFailure
	mywireWait

	// loginURL is the URL that we post our login information to
	loginURL = "https://login.my-wire.de/index.php"

	// startStr and endStr determine where the mywire response is to be found
	startStr = `<div id="content_popup">`
	endStr   = `</div>`
)

// login does the actual work of trying to logon to mywire.
func login(user, pass string) (status mywireStatus, err error) {
	resp, err := http.PostForm(loginURL, url.Values{
		"user":   {user},
		"pass":   {pass},
		"action": {"login"},
	})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return mywireUnknown, err
		}

		if strings.Contains(line, startStr) {
			if !strings.Contains(line, endStr) {
				// If this happens, you have to update this application
				return mywireUnknown, fmt.Errorf("unknown response format")
			}

			if strings.Contains(line, "WARTEN") {
				return mywireWait, nil
			} else if strings.Contains(line, "Erfolgreich") {
				return mywireSuccess, nil
			} else {
				// Find out what the mywire response actually is
				start := strings.Index(line, startStr) + len(startStr)
				end := strings.LastIndex(line, endStr)
				msg := line[start:end]
				return mywireFailure, fmt.Errorf(msg)
			}
		}
	}

	return mywireUnknown, fmt.Errorf("failure reading mywire response")
}

// isOnline checks if we are online.
func isOnline(path string) (online, mywire bool) {
	catch := fmt.Errorf("detected my-wire.de redirection")

	c := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if strings.Contains(req.URL.Host, "my-wire.de") {
			return catch
		}
		return nil
	}}

	resp, err := c.Get(path)
	if resp != nil && resp.Close {
		resp.Body.Close()
	}
	if this, ok := err.(*url.Error); ok {
		mywire = this.Err == catch
	}
	return err == nil, mywire
}

func main() {
	// Load the configuration, and exit if we cannot.
	conf := LoadConfiguration()
	say := conf.NotifyLevel.TrayNotify

	fmt.Printf("Checking online connectivity... ")
	online, mywire := isOnline(conf.PingURL)
	if online {
		fmt.Println("success.")
		return
	}
	fmt.Println("failed.")
	if !mywire {
		fmt.Println("Error: mywire service not available for login process.")
		say(NotifyCritical, "Connection NON-EXISTENT.")
		os.Exit(1)
	}

	say(NotifyAll, "Connecting to mywire...")
	waited := false
	for {
		resp, err := login(conf.User, conf.Pass)
		switch resp {
		case mywireWait:
			// This is the only case where the loop continues
			if waited {
				fmt.Print(".")
			} else {
				fmt.Printf("Trying every %v milliseconds", conf.WaitTime)
				waited = true
			}
			time.Sleep(time.Duration(conf.WaitTime) * time.Millisecond)
		case mywireSuccess:
			if waited {
				fmt.Println()
			}
			if online, _ = isOnline(conf.PingURL); online {
				fmt.Println("ONLINE")
				say(NotifyNormal, "Connection ONLINE.")
				return
			} else {
				fmt.Println("Fatal error: mywire response indicates success, but still no connectivity.")
				say(NotifyCritical, "Connection FAILED")
				os.Exit(1)
			}
		case mywireFailure:
			if waited {
				fmt.Println()
			}
			fmt.Printf("Fatal error: mywire responded, %v\n", err)
			say(NotifyCritical, "Connection FAILED")
			os.Exit(1)
		default:
			if waited {
				fmt.Println()
			}
			fmt.Printf("Fatal error: %v.\n", err)
			say(NotifyCritical, "Connection FAILED")
			os.Exit(1)
		}
	}
}
