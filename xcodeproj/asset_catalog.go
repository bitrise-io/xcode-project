package xcodeproj

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/xcode-project/serialized"
)

// TargetsToAssetCatalogs maps target names to an array of asset catalog namess
type TargetsToAssetCatalogs map[string][]string

// AssetCatalogs parses a xcode project and returns targets mnapped to asset catalogs
func AssetCatalogs(projectPth string) (TargetsToAssetCatalogs, error) {
	absPth, err := pathutil.AbsPath(projectPth)
	if err != nil {
		return TargetsToAssetCatalogs{}, err
	}

	objects, projectID, err := open(absPth)

	p, err := parseProj(projectID, objects)
	if err != nil {
		return TargetsToAssetCatalogs{}, err
	}

	return assetCatalogs(p.Targets, objects)
}

func assetCatalogs(targets []Target, objects serialized.Object) (TargetsToAssetCatalogs, error) {
	targetToAssetCatalogs := map[string][]string{}
	for _, target := range targets {
		resourcesBuildPhase, err := filterResourcesBuildPhase(target.buildPhaseIDs, objects)
		if err != nil {
			return TargetsToAssetCatalogs{}, fmt.Errorf("getting resource build phases failed, error: %s", err)
		}
		assetCatalogs, err := filterAssetCatalogs(resourcesBuildPhase, objects)
		if err != nil {
			return TargetsToAssetCatalogs{}, err
		}
		targetToAssetCatalogs[target.ID] = assetCatalogs
	}
	return targetToAssetCatalogs, nil
}

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
	return resourcesBuildPhase{}, fmt.Errorf("resource build phase not found")
}

func filterAssetCatalogs(buildPhase resourcesBuildPhase, objects serialized.Object) ([]string, error) {
	assetCatalogs := []string{}
	for _, fileUUID := range buildPhase.files {
		buildFile, err := parseBuildFile(fileUUID, objects)
		if err != nil {
			// ignore:
			// D0177B971F26869C0044446D /* (null) in Resources */ = {isa = PBXBuildFile; };
			continue
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
