import React from "react";
import { Field, Resource, Schema } from "./schema";

export interface FieldProps {
  source: string;
}

export type ComponentGeneratorFunc = (schema: Schema, field: Field) => React.FC<FieldProps>;
export type ResourceGeneratorFunc = (schema: Schema, resource: Resource) => React.FC<FieldProps>;
export type LayoutComponentGeneratorFunc = (schema: Schema) => React.FC<FieldProps>;


export interface ComponentGenerator {
  fieldCompGen: ComponentGeneratorFunc
  inputCompGen: ComponentGeneratorFunc
};

export type FieldGenerators = {
  [key: string]: ComponentGenerator;
};

export type GenerateResourceCompTypeFunc = (schema: Schema, res: Resource) => React.FC;

export interface ResourceGenerators {
  list: ResourceGeneratorFunc
  create: ResourceGeneratorFunc
  edit: ResourceGeneratorFunc
};
