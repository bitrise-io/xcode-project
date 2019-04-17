package xcodeproj

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/xcode-project/serialized"
)

func filterResourcesBuildPhase(buildPhases []string, objects serialized.Object) (resourcesBuildPhase, error) {
	for _, buildPhaseUUID := range buildPhases {
		rawBuildPhase, err := objects.Object(buildPhaseUUID)
		if err != nil {
			return resourcesBuildPhase{}, err
		}
		if isResourceBuildPhase(rawBuildPhase) {
			buildPhrase, err := parseResourcesBuildPhase(buildPhaseUUID, objects)
			if err != nil {
				return resourcesBuildPhase{}, fmt.Errorf("failed to parse ResourcesBuildPhase, error: %s", err)
			}
			return buildPhrase, nil
		}
	}
	return resourcesBuildPhase{}, fmt.Errorf("not found")
}

func filterAssetCatalogs(buildPhase resourcesBuildPhase, objects serialized.Object) ([]string, error) {
	var assetCatalogs []string
	for _, fileUUID := range buildPhase.files {
		buildFile, err := parseBuildFile(fileUUID, objects)
		if err != nil {
			return nil, err
		}

		// can be PBXVariantGroup or PBXFileReference
		rawElement, err := objects.Object(buildFile.fileRef)
		if err != nil {
			return nil, err
		}
		if ok, err := isFileReference(rawElement); err != nil {
			return nil, err
		} else if !ok {
			// ignore PBXVariantGroup
			continue
		}

		fileReference, err := parseFileReference(buildFile.fileRef, objects)
		if err != nil {
			return nil, err
		}

		const xcassetsExt = ".xcassets"
		if strings.HasSuffix(fileReference.path, xcassetsExt) {
			assetCatalogs = append(assetCatalogs, fileReference.path)
		}
	}
	return assetCatalogs, nil
}
