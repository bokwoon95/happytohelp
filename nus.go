package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

func nusRedirect(returnTo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Construct NUS OpenID url to redirect user to
		q := url.Values{}
		q.Add("openid.ns", "http://specs.openid.net/auth/2.0")
		q.Add("openid.mode", "checkid_setup")
		q.Add("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")
		q.Add("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
		q.Add("openid.return_to", returnTo)
		q.Add("openid.sreg.required", "email,nickname,fullname")
		u, err := url.Parse("https://openid.nus.edu.sg/server/")
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}
}

func nusAuthenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("openid.sreg.nickname")
		displayname := r.FormValue("openid.sreg.fullname")
		email := r.FormValue("openid.sreg.email")
		if username == "" || displayname == "" || email == "" {
			fmt.Fprintf(w, "Either nickname/fullname/email is empty nickname:%s fullname:%s email:%s", username, displayname, email)
			return
		}
		queries := r.URL.Query()
		queries["openid.mode"] = []string{"check_authentication"}
		resp, err := http.PostForm("https://openid.nus.edu.sg/server", queries)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		sbody := string(body)
		match := regexp.MustCompile(`is_valid:(\w+)`).FindStringSubmatch(sbody)
		if match == nil {
			fmt.Fprintf(w, "is_valid missing from nus openid response %s", sbody)
			return
		}
		if match[1] != "true" {
			fmt.Fprintf(w, "is_valid is not true from NUS OpenID response %s", sbody)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", username)
		ctx = context.WithValue(ctx, "displayname", displayname)
		ctx = context.WithValue(ctx, "email", email)
		next(w, r.WithContext(ctx))
	}
}
