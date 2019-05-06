package xcodeproj

import (
	"fmt"
	"path"

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

type sourceTree int

const (
	unsupportedParent sourceTree = iota
	groupParent
	absoluteParentPath
	undefinedParent
)

// PBXFileReference
// 47C11A4921FF63970084FD7F /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
type fileReference struct {
	id         string
	path       string
	sourceTree sourceTree
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

	sourceTreeRaw, err := rawFileReference.String("sourceTree")
	if err != nil {
		return fileReference{}, err
	}

	var sourceTree sourceTree
	switch sourceTreeRaw {
	case "<group>":
		sourceTree = groupParent
	case "<absolute>":
		sourceTree = absoluteParentPath
	case "":
		sourceTree = undefinedParent
	default:
		sourceTree = unsupportedParent
	}

	return fileReference{
		id:         id,
		path:       path,
		sourceTree: sourceTree,
	}, nil
}

// PBXGroup example:
// 01801EC11A3360B1002B4718 /* Resources */ = {
// 	isa = PBXGroup;
// 	children = (
// 		A045E5E11EDC5C1700BC8A92 /* Localizable.strings */,
// 		01801EA51A32CA2A002B4718 /* Images.xcassets */,
// 	);
// 	name = Resources;
// 	sourceTree = "<group>";
// };

func resolveFileReferenceAbsolutePath(fileReference fileReference, projectID string, objects serialized.Object) (string, error) {
	project, err := objects.Object(projectID)
	if err != nil {
		return "", err
	}

	projectDirPath, err := project.String("projectDirPath")
	if err != nil {
		return "", fmt.Errorf("key projectDirPath not found, project: %s, error: %s", project, err)
	}
	projectRoot, err := project.String("projectRoot")
	if err != nil {
		return "", fmt.Errorf("key projectRoot not found, project: %s, error: %s", project, err)
	}
	mainGroup, err := project.String("mainGroup")
	if err != nil {
		return "", fmt.Errorf("key mainGroup not found, project: %s, error: %s", project, err)
	}

	fmt.Printf("dirPath: %s, root: %s, mainGroup: %s", projectDirPath, projectRoot, mainGroup)

	resolvedPath, err := resolvePath(mainGroup, fileReference.id, path.Join(projectDirPath, projectRoot), objects, &[]string{projectID})
	if err != nil {
		return "", fmt.Errorf("failed to resolve path, error: %s", err)
	}
	return resolvedPath, nil
}

func resolvePath(id string, target string, parentPath string, objects serialized.Object, visited *[]string) (string, error) {
	*visited = append(*visited, id)
	entry, err := objects.Object(id)
	if err != nil {
		return "", fmt.Errorf("object not found, id: %s, error: %s", id, err)
	}

	entryPath, err := entry.String("path")
	if err != nil {
	}
	sourceTreeRaw, err := entry.String("sourceTree")
	if err != nil {
		return "", err
	}
	var sourceTree sourceTree
	switch sourceTreeRaw {
	case "<group>":
		sourceTree = groupParent
	case "<absolute>":
		sourceTree = absoluteParentPath
	case "":
		sourceTree = undefinedParent
	default:
		sourceTree = unsupportedParent
	}

	var currentProjRelPath string
	switch sourceTree {
	case groupParent:
		currentProjRelPath = path.Join(parentPath, entryPath)
	case absoluteParentPath:
		currentProjRelPath = entryPath
	case undefinedParent:
		currentProjRelPath = parentPath
	case unsupportedParent:
		return "", fmt.Errorf("failed to resolve path, unsupported relation: %s", sourceTreeRaw)
	}

	if id == target {
		return currentProjRelPath, nil
	}

	children, err := entry.StringSlice("children")
	if err != nil {
		// return "", fmt.Errorf("key children not found, entry: %s, error: %s", entry, err)
		return "", nil
	}
	for _, child := range children {
		resolvedPath, err := resolvePath(child, target, currentProjRelPath, objects, visited)
		if err != nil {
			return "", err
		} else if resolvedPath != "" {
			return resolvedPath, nil
		}
	}
	return "", nil
}
