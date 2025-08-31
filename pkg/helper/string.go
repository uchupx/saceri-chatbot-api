package helper

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// GenerateRandomString creates a random string of the specified length
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func StringPointer(v string) *string {
	return &v
}

func StructToString(v interface{}) (*string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return StringPointer(string(b)), nil
}

func StringToTime(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}
	}

	return t
}
