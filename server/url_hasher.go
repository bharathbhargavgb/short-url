package main

import(
    "math/rand"
    "time"
    "strconv"
    "crypto/sha1"
    "io"

    "github.com/eknkc/basex"
)


func getValidShortID(dataStore *DBStore, URI string) string {
    hashValue := generateShortID(URI)

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

func generateShortID(URI string) string {
    hasher := sha1.New()
    io.WriteString(hasher, URI)
    io.WriteString(hasher, getRandomNumber())
    return base58Encode(hasher.Sum(nil))
}

func base58Encode(input []byte) string {
    // base58 excludes visually ambiguous characters - O 0 I l
    base58_charSet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
    encoder, _ := basex.NewEncoding(base58_charSet)
    return encoder.Encode(input)
}

func getRandomNumber() string {
    seed := rand.NewSource(time.Now().UnixNano())
    r := rand.New(seed)
    return strconv.FormatUint(r.Uint64(), 10)
}

