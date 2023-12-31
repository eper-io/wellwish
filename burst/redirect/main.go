package main

import (
	"fmt"
	"gitlab.com/eper.io/engine/metadata"
	"net/http"
)

// This document is Licensed under Creative Commons CC0.
// To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights
// to this document to the public domain worldwide.
// This document is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along with this document.
// If not, see https://creativecommons.org/publicdomain/zero/1.0/legalcode.

// This code shows how to run a web server

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://eper.io", http.StatusTemporaryRedirect)
	})

	err := http.ListenAndServe(metadata.Http11Port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
