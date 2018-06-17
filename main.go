package main

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xcode-project/xcodeproj"
)

func main() {
	pth := "pbxproj"
	fmt.Printf("opening project: %s\n", pth)
	proj, err := xcodeproj.Open(pth)
	if err != nil {
		panic(err)
	}

	// fmt.Println("\n\nXcodeProj:")
	// logPretty(proj)

	fmt.Printf("defaultConfigurationName: %s\n", proj.Proj.BuildConfigurationList.DefaultConfigurationName)
	for i, buildConfiguration := range proj.Proj.BuildConfigurationList.BuildConfigurations {
		fmt.Printf("%d. buildConfiguration: %s\n", i+1, buildConfiguration.Name)
	}

	fmt.Println("targets:")
	for i, target := range proj.Proj.Targets {
		fmt.Printf("%d. target: %s\n", i+1, target.Name)

		for _, buildConfiguration := range target.BuildConfigurationList.BuildConfigurations {
			if buildConfiguration.Name == target.BuildConfigurationList.DefaultConfigurationName {
				bundleID := buildConfiguration.BuildSettings["PRODUCT_BUNDLE_IDENTIFIER"].(string)
				fmt.Printf("bundleID: %s\n", bundleID)
				break
			}
		}
	}
}

func logPretty(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Printf("%+v\n", string(b))
}
