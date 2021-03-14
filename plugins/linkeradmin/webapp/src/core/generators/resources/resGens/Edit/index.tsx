import React from "react";
import { Edit, SimpleForm } from "react-admin";
import { ResourceGenerator } from "../../../../../types/generators";
import { Field, Resource } from "../../../../../types/schema";
import { generateInputCompFromField } from "../../../fields";


const EditGenerator: ResourceGenerator =  (r: Resource) => {
  const fields = r.update_inputs.length > 0 ? r.update_inputs : r.create_inputs;
  const comps = fields.map((f: Field) => {
    const Comp = generateInputCompFromField(f);
    const optProps = (f.options && f.options.props) || {};
    return (
      <Comp key={f.name} source={f.name} {...optProps} />
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
