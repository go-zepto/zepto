package linkeradmin

type Admin struct {
	Resources []*Resource `json:"resources"`
}

func NewAdmin() *Admin {
	return &Admin{
		Resources: make([]*Resource, 0),
	}
}

func (a *Admin) AddResource(res *Resource) *Admin {
	a.Resources = append(a.Resources, res)
	return a
}

func (a *Admin) findResourceIndexByName(name string) int {
	for i := 0; i < len(a.Resources); i++ {
		if a.Resources[i].Name == name {
			return i
		}
	}
	return -1
}

func (a *Admin) removeResourceAtIndex(index int) *Admin {
	a.Resources = append(a.Resources[:index], a.Resources[index+1:]...)
	return a
}

func (a *Admin) RemoveResource(name string) *Admin {
	idx := a.findResourceIndexByName(name)
	if idx >= 0 {
		a.removeResourceAtIndex(idx)
	}
	return a
}
