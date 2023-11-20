package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	"gorm.io/datatypes"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func BodyBuilder(status bool) map[string]interface{} {
	return map[string]interface{}{"success": status}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RandStringRunes(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func ChangeTimeFormat(myTime string) string {
	layout1 := "03:04 PM"
	layout2 := "15:04"
	t, err := time.Parse(layout1, myTime)
	if err != nil {
		fmt.Println(err)
		return "err"
	}
	return t.Format(layout2)
}

func JSONToSlice(jsonData *datatypes.JSON, slice *[]string) error {
	err := json.Unmarshal(*jsonData, slice)
	fmt.Println(jsonData)
	fmt.Println(slice)
	return err
}

func SliceToJson(slice *[]string) datatypes.JSON {
	inBytes, _ := json.Marshal(slice)

	inJson := postgres.Jsonb{inBytes}	
	return datatypes.JSON(inJson.RawMessage)
}
