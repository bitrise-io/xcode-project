package xcodeproj

import (
	"io/ioutil"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"howett.net/plist"
)

// Pbxproj ...
type Pbxproj struct {
	raw  string
	json map[string]interface{}
	// TODO
}

// New ...
func New(pth string) *Pbxproj {
	raw := raw(pth)
	return &Pbxproj{
		raw:  raw,
		json: json(raw),
	}
}

// Raw ...
func (p Pbxproj) Raw() string {
	return p.raw
}

// Json ...
func (p Pbxproj) Json() map[string]interface{} {
	return p.json
}

//
// Private methods

func json(raw string) map[string]interface{} {
	var test interface{}
	if _, err := plist.Unmarshal([]byte(raw), &test); err != nil {
		failf("Failed to generate json from Pbxproj - error: %s", err)
	}

	return test.(map[string]interface{})
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
	log.Warnf("For more details you can enable the debug logs by turning on the verbose step input.")
	os.Exit(1)
}
