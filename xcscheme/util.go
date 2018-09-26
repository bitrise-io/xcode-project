package xcscheme

import (
	"path/filepath"
)

// FindSchemesIn ...
func FindSchemesIn(root string) ([]Scheme, error) {
	//
	// Add the shared schemes to the list
	pths, err := pathsByPattern(root, "xcshareddata", "xcschemes", "*.xcscheme")
	if err != nil {
		return nil, err
	}

	var schemes []Scheme
	for _, pth := range pths {
		scheme, err := Open(pth)
		if err != nil {
			return nil, err
		}
		schemes = append(schemes, scheme)
	}

	//
	// Add the non-shared user schemes to the list
	pths, err = pathsByPattern(root, "xcuserdata", "*.xcuserdatad", "xcschemes", "*.xcscheme")
	if err != nil {
		return nil, err
	}

	for _, pth := range pths {
		scheme, err := Open(pth)
		if err != nil {
			return nil, err
		}
		schemes = append(schemes, scheme)
	}
	return schemes, nil
}

func pathsByPattern(paths ...string) ([]string, error) {
	pattern := filepath.Join(paths...)
	return filepath.Glob(pattern)
}
