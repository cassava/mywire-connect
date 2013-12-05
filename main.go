// Copyright (c) 2013, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// mywire-connect is a program to initiate the connection to mywire.
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

const (
	// pingURL is the URL to use for testing internet connectivity, we are
	// using Google because they are usually pretty quick in their responses.
	pingURL = "http://www.google.com"

	// wait waitTime seconds if mywire tells us to wait; make sure that this
	// is always greater than 0, otherwise we will send requests as fast as
	// we can, and that would not be very good.
	waitTime = 2

	// confUser and confPass are the environment variables that are read for
	// the mywire username and password, respectively.
	confUser = "MYWIRE_USER"
	confPass = "MYWIRE_PASS"
)

func readUserAndPass() (user, pass string) {
	user = os.Getenv(confUser)
	if user == "" {
		fmt.Printf("Fatal error: environment variable %v not set.\n", confUser)
	}
	pass = os.Getenv(confPass)
	if pass == "" {
		fmt.Printf("Fatal error: environment variable %v not set.\n", confPass)
	}

	if user == "" || pass == "" {
		os.Exit(1)
	}

	return
}

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
	fmt.Printf("Checking online connectivity... ")
	online, mywire := isOnline(pingURL)
	if online {
		fmt.Println("success.")
		return
	}
	fmt.Println("failed.")
	if !mywire {
		fmt.Println("Error: mywire service not available for login process.")
		os.Exit(1)
	}

	user, pass := readUserAndPass()
	for {
		resp, err := login(user, pass)
		switch resp {
		case mywireWait:
			// This is the only case where the loop continues
			fmt.Printf("Waiting %v seconds...\n", waitTime)
			time.Sleep(waitTime * time.Second)
		case mywireSuccess:
			if online, _ = isOnline(pingURL); online {
				fmt.Println("ONLINE")
				return
			} else {
				fmt.Println("Fatal error: mywire response indicates success, but still no connectivity.")
				os.Exit(1)
			}
		case mywireFailure:
			fmt.Printf("Fatal error: mywire responded, %v\n", err)
			os.Exit(1)
		default:
			fmt.Printf("Fatal error: %v.\n", err)
			os.Exit(1)
		}
	}
}
