import React from "react";
import { Field, ComponentGenerator, FieldGenerators, FieldProps, Schema, Resource, ComponentGeneratorFunc } from "../../types";
interface GeneratorsOption {
    schema: Schema;
    defaultRowClick: 'edit' | 'show';
}
export default class Generators {
    fieldGenerators: FieldGenerators;
    schema: Schema;
    defaultRowClick: 'edit' | 'show';
    constructor(opts: GeneratorsOption);
    registerComponentGenerator(name: string, compGen: ComponentGenerator): void;
    generateFieldComp(field: Field): React.FC<FieldProps>;
    generateInputComp(field: Field): React.FC<FieldProps>;
    _generateResourceComps(fields: Field[], genFunc: ComponentGeneratorFunc, options?: {
        props: any;
    }): JSX.Element[];
    _generateShowResource(resource: Resource): React.FC;
    _generateListResource(resource: Resource): React.FC;
    _generateCreateResource(resource: Resource): (props: any) => JSX.Element;
    _generateEditResource(resource: Resource): (props: any) => JSX.Element;
    generateResourceComp(resource: Resource): JSX.Element;
}
export {};
