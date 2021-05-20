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

func TestBase58Encode(t *testing.T) {
    input := "HelloWorld!"
    expected_output := "JxF12TsQucjW69N"
    actual_output := base58Encode([]byte(input))
    if actual_output != expected_output {
        t.Errorf("Want [%v], got [%v]", expected_output, actual_output)
    }
}
