package v1alpha1

func UpdateStatusSecondaryResources(secondaryResources map[string][]string, kind string, resourceName string) map[string][]string {
	// If the map is nil, instansiate it
	if secondaryResources == nil {
		secondaryResources = make(map[string][]string)
	}

	// return if the resource name already exists in the slice
	for _, ele := range secondaryResources[kind] {
		if ele == resourceName {
			return secondaryResources
		}
	}
	// add the resource name to the list of secondary resources in the status
	secondaryResources[kind] = append(secondaryResources[kind], resourceName)

	// return new map
	return secondaryResources
}
