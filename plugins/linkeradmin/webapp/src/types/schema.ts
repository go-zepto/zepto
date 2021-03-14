export interface Field {
  name: string;
  type: string;
  options: any; 
}

type Input = Field

export interface Resource {
  name: string;
  endpoint: string;
  list_fields: Field[];
  create_inputs: Input[];
  update_inputs: Input[];
}

export interface Schema {
  resources: Resource[];
}
