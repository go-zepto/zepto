export interface Field {
  name: string;
  type: string;
  options: Object; 
}


export interface Resource {
  name: string;
  endpoint: string;
  fields: Field[];
}

export default interface Schema {
  resources: Resource[];
}
