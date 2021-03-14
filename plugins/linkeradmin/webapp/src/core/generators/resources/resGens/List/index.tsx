import React from "react";
import { Datagrid, List } from "react-admin";
import { ResourceGenerator } from "../../../../../types/generators";
import { Field, Resource } from "../../../../../types/schema";
import { generateFieldCompFromField } from "../../../fields";


const ListGenerator: ResourceGenerator =  (r: Resource) => {
  const comps = r.fields.map((f: Field) => {
    const Comp = generateFieldCompFromField(f);
    return (
      <Comp source={f.name} />
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
