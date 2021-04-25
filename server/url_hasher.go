package main

import(
    "math/rand"
    "time"
    "strconv"
    "crypto/sha1"
    "encoding/base64"
    "io"
)

func getValidTinyID(dataStore *DBStore, URI string) string {
    hasher := sha1.New()
    io.WriteString(hasher, URI)
    io.WriteString(hasher, getRandomNumber())

    // TODO: Use base62 instead of base64 
    hashValue := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

    for keyLen := 6; keyLen <= 8; keyLen++ {
        for i := 0; i < 3; i++ {
            key := hashValue[i*keyLen: (i+1)*keyLen]
            shortItem, err := dataStore.getItem(key)
            if shortItem == nil && err == nil {
                return key
            }
        }
    }
    return ""
}

func getRandomNumber() string {
    seed := rand.NewSource(time.Now().UnixNano())
    r := rand.New(seed)
    return strconv.FormatUint(r.Uint64(), 10)
}

