package tools

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
)

func Encrypt(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func ConvertStruct2Map(structObject any) map[string]interface{} {
	data := make(map[string]interface{})
	result, err := json.Marshal(structObject)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(result, &data)
	return data
}
