import React from "react";
import { Datagrid, List } from "react-admin";
import { ResourceGenerator } from "../../../../../types/generators";
import { Field, Resource } from "../../../../../types/schema";
import { generateFieldCompFromField } from "../../../fields";


const ListGenerator: ResourceGenerator =  (r: Resource) => {
  const fields = r.list_fields;
  const comps = fields.map((f: Field) => {
    const Comp = generateFieldCompFromField(f);
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