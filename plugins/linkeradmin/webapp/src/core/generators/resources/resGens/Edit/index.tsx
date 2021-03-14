import React from "react";
import { Edit, SimpleForm } from "react-admin";
import { ResourceGenerator } from "../../../../../types/generators";
import { Field, Resource } from "../../../../../types/schema";
import { generateInputCompFromField } from "../../../fields";


const EditGenerator: ResourceGenerator =  (r: Resource) => {
  const comps = r.fields.map((f: Field) => {
    const Comp = generateInputCompFromField(f);
    return (
      <Comp source={f.name} />
    );
  })
  return (props: any) => (
    <Edit {...props}>
      <SimpleForm>
          {comps.map((Comp: any) => Comp)}
      </SimpleForm>
    </Edit>
  );
};

export default EditGenerator;
