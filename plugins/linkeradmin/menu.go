package linkeradmin

type MenuLink struct {
	Icon               string `json:"icon"`
	Label              string `json:"label"`
	LinkToResourceName string `json:"link_to_resource_name"`
	LinkToPath         string `json:"link_to_path"`
}

type Menu struct {
	Links []MenuLink `json:"links"`
}
