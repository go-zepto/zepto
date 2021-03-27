export interface Field {
  name: string;
  type: string;
  options: any; 
}

type Input = Field


export interface ResourceFieldEndpoint {
  fields: Field[];
}

export interface ResourceInputEndpoint {
  inputs: Input[];
}

export interface Resource {
  name: string;
  endpoint: string;
  list_endpoint: ResourceFieldEndpoint;
  create_endpoint: ResourceInputEndpoint;
  update_endpoint: ResourceInputEndpoint;
}

export interface MenuLink {
  icon: string;
  label: string;
  link_to_resource_name: string;
  link_to_path: string;
}

export interface Menu {
  links: MenuLink[];
}

export interface Admin {
  resources: Resource[];
}

export interface Schema {
  menu: Menu;
  admin: Admin;
}
