package xcodeproj

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-tools/xcode-project/serialized"
)

func parseShowBuildSettingsOutput(out string) (serialized.Object, error) {
	settings := serialized.Object{}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Build settings for") {
			continue
		}

		split := strings.Split(line, " = ")

		if len(split) < 2 {
			return nil, fmt.Errorf("unkown build settings: %s", line)
		}

		key := strings.TrimSpace(split[0])
		value := strings.TrimSpace(strings.Join(split[1:], " = "))

		settings[key] = value
	}

	return settings, nil
}

func showBuildSettings(project, target, configuration string) (serialized.Object, error) {
	cmd := command.New("xcodebuild", "-project", project, "-target", target, "-configuration", configuration, "-showBuildSettings")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s failed: %s", cmd.PrintableCommandArgs(), err)
	}

	return parseShowBuildSettingsOutput(out)
}
