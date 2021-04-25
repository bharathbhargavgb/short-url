package main

import(
    "testing"
)

func TestGetValidTinyID(t *testing.T) {
    db := createFakeStorage()
    tinyID := getValidTinyID(db, "https://www.google.com")
    if len(tinyID) != 6 {
        t.Errorf("Want len(tinyID) = 6, but got len(%v) = %v", tinyID, len(tinyID))
    }
}
