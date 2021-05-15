package main

import (
    "fmt"
    "testing"
)

func TestGetItem(t *testing.T) {
    db := createFakeStorage()
    want := "https://www.google.com"
    shortLink, err := db.getItem("goog")
    if err != nil {
        t.Errorf("getItem returned error: %q", err)
    }
    if want != shortLink.URI {
        t.Errorf("Got %q, Want %q", shortLink.URI, want)
    }
}

func TestGetItemMissingKey(t *testing.T) {
    db := createFakeStorage()
    shortLink, err := db.getItem("randomString")
    if err != nil {
        t.Errorf("Want nil, Got %+v", err)
    }
    if shortLink != nil {
        t.Error(fmt.Sprintf("Want nil, Got %+v", *shortLink))
    }
}

func TestPutItem(t *testing.T) {
    db := createFakeStorage()
    shortItem := &shortURI {
        ShortID: "goog",
        URI: "https://www.google.com",
    }
    err := db.putItem(shortItem)
    if err != nil {
        t.Error("Error inserting to DB")
    }
}

