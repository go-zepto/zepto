export interface Field {
    name: string;
    type: string;
    options: any;
}
declare type Input = Field;
export interface Resource {
    name: string;
    endpoint: string;
    list_fields: Field[];
    create_inputs: Input[];
    update_inputs: Input[];
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
export interface Schema {
    menu: Menu;
    resources: Resource[];
}
export {};
