package main

import(
    "testing"
)

func TestGetValidShortID(t *testing.T) {
    db := createFakeStorage()
    shortID := getValidShortID(db, "https://www.google.com")
    if len(shortID) != 6 {
        t.Errorf("Want len(shortID) = 6, but got len(%v) = %v", shortID, len(shortID))
    }
}
