package xcodeproj

import (
	"io/ioutil"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"howett.net/plist"
)

// Pbxproj ...
type Pbxproj struct {
	Raw  string
	JSON map[string]interface{}

	ArchiveVersion string  `json:"archiveVersion"`
	Classes        Classes `json:"classes"`
	ObjectVersion  string  `json:"objectVersion"`

	Objects map[string]interface{}
}

// Classes ...
type Classes struct {
	// Unknown TODO later
}

// New ...
func New(pth string) *Pbxproj {
	raw := raw(pth)
	p := Pbxproj{
		Raw:  raw,
		JSON: unmarshal(raw),
	}

	if _, err := plist.Unmarshal([]byte(raw), &p); err != nil {
		failf("Failed to generate json from Pbxproj - error: %s", err)
	}

	p.Objects = p.JSON["objects"].(map[string]interface{})
	return &p
}

//
// Private methods

func unmarshal(raw string) map[string]interface{} {
	var i interface{}
	if _, err := plist.Unmarshal([]byte(raw), &i); err != nil {
		failf("Failed to generate json from Pbxproj - error: %s", err)
	}

	return i.(map[string]interface{})
}

func raw(pth string) string {
	b, err := ioutil.ReadFile(pth)
	if err != nil {
		failf("Failed to read file - error: %s", err)
	}

	return string(b)
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
