package xcworkspace

import "github.com/bitrise-tools/xcode-project/xcodeproj"

func Example() {
	workspace, err := Open("workspace.xcworkspace")
	if err != nil {
		panic(err)
	}

	var fileRefs []FileRef
	for _, fileRef := range workspace.FileRefs {
		fileRefs = append(fileRefs, fileRef)
	}
	for _, group := range workspace.Groups {
		for _, fileRef := range group.FileRefs {
			fileRefs = append(fileRefs, fileRef)
		}
	}

	var projects []xcodeproj.XcodeProj
	for _, fileRef := range fileRefs {
		pth, err := fileRef.AbsPath("workspace_dir")
		if err != nil {
			panic(err)
		}

		if !xcodeproj.IsXcodeProj(pth) {
			continue
		}

		project, err := xcodeproj.Open(pth)
		if err != nil {
			panic(err)
		}
		projects = append(projects, project)
	}
}
