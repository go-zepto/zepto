import React from "react";
import { Field, Resource, Schema } from "./schema";
export interface FieldProps {
    source: string;
}
export declare type ComponentGeneratorFunc = (schema: Schema, field: Field) => React.FC<FieldProps>;
export declare type ResourceGeneratorFunc = (schema: Schema, resource: Resource) => React.FC<FieldProps>;
export declare type LayoutComponentGeneratorFunc = (schema: Schema) => React.FC<FieldProps>;
export interface ComponentGenerator {
    fieldCompGen: ComponentGeneratorFunc;
    inputCompGen: ComponentGeneratorFunc;
}
export declare type FieldGenerators = {
    [key: string]: ComponentGenerator;
};
export declare type GenerateResourceCompTypeFunc = (schema: Schema, res: Resource) => React.FC;
export interface ResourceGenerators {
    list: ResourceGeneratorFunc;
    create: ResourceGeneratorFunc;
    edit: ResourceGeneratorFunc;
}
