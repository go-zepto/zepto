import React from "react";
import { Datagrid, List } from "react-admin";
import { ResourceGenerator } from "../../../../../types/generators";
import { Field, Resource, Schema } from "../../../../../types/schema";
import { generateFieldCompFromField } from "../../../fields";


const ListGenerator: ResourceGenerator =  (s: Schema, r: Resource) => {
  const fields = r.list_fields;
  const comps = fields.map((f: Field) => {
    const Comp = generateFieldCompFromField(s, f);
    const optProps = (f.options && f.options.props) || {};
    return (
      <Comp key={f.name} source={f.name} {...optProps} />
    );
  })
  return (props: any) => (
    <List {...props}>
      <Datagrid rowClick="edit">
        {comps.map((Comp: any) => Comp)}
      </Datagrid>
    </List>
  );
};

export default ListGenerator;
