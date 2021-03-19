import React from "react";
import {
  Field,
  ComponentGenerator,
  FieldGenerators,
  FieldProps,
  Schema,
  Resource,
  ComponentGeneratorFunc,
} from "../../types";
import { DEFAULT_FIELD_GENERATORS } from "./fields";
import { Create, Datagrid, Edit, List, Resource as ReactAdminResource, Show, SimpleForm, SimpleShowLayout } from "react-admin";


interface GeneratorsOption {
  schema: Schema;
  defaultRowClick: 'edit' | 'show';
}

export default class Generators {
  fieldGenerators: FieldGenerators;
  schema: Schema;
  defaultRowClick: 'edit' | 'show';

  constructor(opts: GeneratorsOption) {
    this.fieldGenerators = DEFAULT_FIELD_GENERATORS;
    this.schema = opts.schema;
    this.defaultRowClick = opts.defaultRowClick ?? 'edit';
  }

  registerComponentGenerator(name: string, compGen: ComponentGenerator) {
    this.fieldGenerators[name] = compGen;
  }

  generateFieldComp(field: Field): React.FC<FieldProps> {
    return this.fieldGenerators[field.type].fieldCompGen(this.schema, field);
  }

  generateInputComp(field: Field): React.FC<FieldProps> {
    return this.fieldGenerators[field.type].inputCompGen(this.schema, field);
  }

  _generateResourceComps(fields: Field[], genFunc: ComponentGeneratorFunc, options?: { props: any }) {
    const comps = fields.map((f: Field) => {
      const Comp = genFunc(this.schema, f);
      const optProps = (f.options && f.options.props) || {};
      return (
        <Comp key={f.name} source={f.name} {...optProps} {...options?.props} />
      );
    });
    return comps;
  }

  _generateShowResource(resource: Resource): React.FC {
    const fields = resource.list_fields;
    const fieldGen = (s: Schema, f: Field) => this.generateFieldComp(f);
    const comps = this._generateResourceComps(fields, fieldGen, { props: { addLabel: true }});
    return (props: any) => (
      <Show {...props}>
        <SimpleShowLayout>
          {comps}
        </SimpleShowLayout>
      </Show>
    );
  }

  _generateListResource(resource: Resource): React.FC {
    const fields = resource.list_fields;
    const fieldGen = (s: Schema, f: Field) => this.generateFieldComp(f);
    const comps = this._generateResourceComps(fields, fieldGen);
    return (props: any) => (
      <List {...props}>
        <Datagrid rowClick={this.defaultRowClick}>
          {comps}
        </Datagrid>
      </List>
    );
  }

  _generateCreateResource(resource: Resource) {
    const fields = resource.create_inputs;
    const fieldGen = (s: Schema, f: Field) => this.generateInputComp(f);
    const comps = this._generateResourceComps(fields, fieldGen);
    return (props: any) => (
      <Create {...props}>
          <SimpleForm>
            {comps.map((Comp: any) => Comp)}
          </SimpleForm>
      </Create>
    );
  }

  _generateEditResource(resource: Resource) {
    const fields = resource.update_inputs.length > 0 ? resource.update_inputs : resource.create_inputs;
    const fieldGen = (s: Schema, f: Field) => this.generateInputComp(f);
    const comps = this._generateResourceComps(fields, fieldGen);
    return (props: any) => (
      <Edit {...props}>
        <SimpleForm>
          {comps.map((Comp: any) => Comp)}
        </SimpleForm>
      </Edit>
    );
  }

  generateResourceComp(resource: Resource) {
    const list = this._generateListResource(resource);
    const show = this._generateShowResource(resource);
    const create = this._generateCreateResource(resource);
    const edit = this._generateEditResource(resource);
    return (
      <ReactAdminResource
        key={resource.name}
        name={resource.endpoint.toLowerCase()}
        list={list}
        show={this.defaultRowClick == 'show' && show || undefined}
        create={create}
        edit={edit}
      />
    );
  }
}
