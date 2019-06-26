package xcodeproj

import "github.com/bitrise-io/xcode-project/serialized"

// ProjectAtributes ...
type ProjectAtributes struct {
	TargetAttributes map[string]TargetAttributes
}

// TargetAttributes ...
type TargetAttributes struct {
	CreatedOnToolsVersion string
	LastSwiftMigration    string
	DevelopmentTeam       string
	ProvisioningStyle     string
}

func parseProjectAttributes(rawPBXProj serialized.Object) (ProjectAtributes, error) {
	var attributes ProjectAtributes
	attributesObject, err := rawPBXProj.Object("attributes")
	if err != nil {
		return ProjectAtributes{}, err
	}

	attributes.TargetAttributes, err = parseTargetAttributesMap(attributesObject)
	if err != nil {
		return ProjectAtributes{}, err
	}

	return attributes, nil
}

func parseTargetAttributesMap(attributesObject serialized.Object) (map[string]TargetAttributes, error) {
	targetAttributesObject, err := attributesObject.Object("TargetAttributes")
	if err != nil {
		return nil, err
	}

	targetAttributesMap := make(map[string]TargetAttributes)
	for _, key := range targetAttributesObject.Keys() {
		obj, err := targetAttributesObject.Object(key)
		if err != nil {
			return nil, err
		}

		var t TargetAttributes
		t.CreatedOnToolsVersion, err = obj.String("CreatedOnToolsVersion")
		t.LastSwiftMigration, err = obj.String("LastSwiftMigration")
		t.DevelopmentTeam, err = obj.String("DevelopmentTeam")
		t.ProvisioningStyle, err = obj.String("ProvisioningStyle")

		targetAttributesMap[key] = t
	}
	return targetAttributesMap, nil
}
