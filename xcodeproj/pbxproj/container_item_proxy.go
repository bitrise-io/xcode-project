package pbxproj

// PBXContainerItemProxy ...
type PBXContainerItemProxy struct {
	ISA                  string
	ID                   string
	ContainerPortal      string
	ProxyType            string
	RemoteGlobalIDString string
	RemoteInfo           string
}

var pbxContainerItemProxyByID = map[string]PBXContainerItemProxy{}

// GetPBXContainerItemProxy ...
func GetPBXContainerItemProxy(id string, raw map[string]interface{}) PBXContainerItemProxy {
	if containerItemProxy, ok := pbxContainerItemProxyByID[id]; ok {
		return containerItemProxy
	}

	rawPBXContainerItemProxy := raw[id].(map[string]interface{})

	containerPortal := rawPBXContainerItemProxy["containerPortal"].(string)
	proxyType := rawPBXContainerItemProxy["proxyType"].(string)
	remoteGlobalIDString := rawPBXContainerItemProxy["remoteGlobalIDString"].(string)
	remoteInfo := rawPBXContainerItemProxy["remoteInfo"].(string)

	containerItemProxy := PBXContainerItemProxy{
		ISA:                  "PBXContainerItemProxy",
		ID:                   id,
		ContainerPortal:      containerPortal,
		ProxyType:            proxyType,
		RemoteGlobalIDString: remoteGlobalIDString,
		RemoteInfo:           remoteInfo,
	}
	pbxContainerItemProxyByID[id] = containerItemProxy
	return containerItemProxy
}
