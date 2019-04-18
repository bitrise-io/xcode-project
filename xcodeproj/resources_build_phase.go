package xcodeproj

import (
	"fmt"

	"github.com/bitrise-io/xcode-project/serialized"
)

// resourcesBuildPhase represents a PBXResourcesBuildPhase element
type resourcesBuildPhase struct {
	ID    string
	files []string
}

func isResourceBuildPhase(raw serialized.Object) bool {
	if isa, err := raw.String("isa"); err != nil {
		return false
	} else if isa != "PBXResourcesBuildPhase" {
		return false
	}
	return true
}

func parseResourcesBuildPhase(id string, objects serialized.Object) (resourcesBuildPhase, error) {
	rawResourcesBuildPhase, err := objects.Object(id)
	if err != nil {
		return resourcesBuildPhase{}, err
	}

	if !isResourceBuildPhase(rawResourcesBuildPhase) {
		return resourcesBuildPhase{}, fmt.Errorf("not a PBXResourcesBuildPhase element")
	}

	files, err := rawResourcesBuildPhase.StringSlice("files")
	if err != nil {
		return resourcesBuildPhase{}, err
	}

	return resourcesBuildPhase{
		ID:    id,
		files: files,
	}, nil
}

// buildFile represents a PBXBuildFile element
// 47C11A4A21FF63970084FD7F /* Assets.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 47C11A4921FF63970084FD7F /* Assets.xcassets */; };
type buildFile struct {
	fileRef string
}

func parseBuildFile(id string, objects serialized.Object) (buildFile, error) {
	rawBuildFile, err := objects.Object(id)
	if err != nil {
		return buildFile{}, err
	}
	if isa, err := rawBuildFile.String("isa"); err != nil {
		return buildFile{}, err
	} else if isa != "PBXBuildFile" {
		return buildFile{}, fmt.Errorf("not a PBXBuildFile element")
	}

	fileRef, err := rawBuildFile.String("fileRef")
	if err != nil {
		return buildFile{}, err
	}

	return buildFile{
		fileRef: fileRef,
	}, nil
}

// PBXFileReference
// 47C11A4921FF63970084FD7F /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
type fileReference struct {
	path string
}

const fileReferenceElementType = "PBXFileReference"

func isFileReference(raw serialized.Object) (bool, error) {
	if isa, err := raw.String("isa"); err != nil {
		return false, err
	} else if isa == fileReferenceElementType {
		return true, nil
	}
	return false, nil
}

func parseFileReference(id string, objects serialized.Object) (fileReference, error) {
	rawFileReference, err := objects.Object(id)
	if err != nil {
		return fileReference{}, err
	}

	if ok, err := isFileReference(rawFileReference); err != nil {
		return fileReference{}, err
	} else if !ok {
		return fileReference{}, fmt.Errorf("not a %s element", fileReferenceElementType)
	}

	path, err := rawFileReference.String("path")
	if err != nil {
		return fileReference{}, err
	}

	return fileReference{
		path: path,
	}, nil
}
