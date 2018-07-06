package xcodeproj

import (
	"fmt"
)

func Example() {
	project, err := Open("project.xcodeproj")
	if err != nil {
		panic(err)
	}

	for _, target := range project.Proj.Targets {
		fmt.Printf("%s target default configuration: %s\n", target.Name, target.BuildConfigurationList.DefaultConfigurationName)

		buildConfiguration := target.BuildConfigurationList.BuildConfigurations[0]
		bundleID := buildConfiguration.BuildSettings["PRODUCT_BUNDLE_IDENTIFIER"]
		fmt.Printf("%s target bundle id: %s\n", target.Name, bundleID)
	}

	schemes, err := project.SharedSchemes()
	if err != nil {
		panic(err)
	}

	for _, scheme := range schemes {
		entry := scheme.BuildAction.BuildActionEntries[0]
		targetID := entry.BuildableReference.BlueprintIdentifier

		target, _ := project.Proj.Target(targetID)
		fmt.Printf("%s scheme's main target: %s\n", scheme.Name, target.Name)
	}
}
