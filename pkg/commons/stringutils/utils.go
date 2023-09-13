package stringutils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateUsername(email string) string {
	prefix := strings.Split(email, "@")[0][:8]
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d", prefix, timestamp)
}

func GetRandomImageURL(seed int64, imageList []string) string {
	if len(imageList) == 0 {
		return ""
	}
	r := rand.New(rand.NewSource(seed))
	return imageList[r.Intn(len(imageList))]
}
