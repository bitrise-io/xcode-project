package xcodeproj

import "github.com/bitrise-tools/xcode-project/serialized"

// BuildConfiguration ..
type BuildConfiguration struct {
	ISA           string
	ID            string
	Name          string
	BuildSettings serialized.Object
}

func (p XcodeProj) buildConfiguration(id string) (BuildConfiguration, error) {
	if buildConfiguration, ok := p.buildConfigurationByID[id]; ok {
		return buildConfiguration, nil
	}

	raw, err := p.raw.Object(id)
	if err != nil {
		return BuildConfiguration{}, err
	}

	name, err := raw.String("name")
	if err != nil {
		return BuildConfiguration{}, err
	}

	buildSettings, err := raw.Object("buildSettings")
	if err != nil {
		return BuildConfiguration{}, err
	}

	buildConfiguration := BuildConfiguration{
		ISA:           "XCBuildConfiguration",
		ID:            id,
		Name:          name,
		BuildSettings: buildSettings,
	}

	p.buildConfigurationByID[id] = buildConfiguration

	return buildConfiguration, nil
}
