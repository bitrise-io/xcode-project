package main

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xcode-project-parser/xcodeproj"
)

func main() {
	pbxproj := xcodeproj.New("pbxproj")

	// logPretty(pbxproj.Raw())
	logPretty(pbxproj.Json())
}

func logPretty(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Printf("%+v\n", string(b))
}
