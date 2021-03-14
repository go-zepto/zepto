import React from "react";
import { Create, SimpleForm } from "react-admin";
import { ResourceGenerator } from "../../../../../types/generators";
import { Field, Resource, Schema } from "../../../../../types/schema";
import { generateInputCompFromField } from "../../../fields";


const CreateGenerator: ResourceGenerator =  (s: Schema, r: Resource) => {
  const fields = r.create_inputs;
  const comps = fields.map((f: Field) => {
    const Comp = generateInputCompFromField(s, f);
    const optProps = (f.options && f.options.props) || {};
    return (
      <Comp key={f.name} source={f.name} {...optProps} />
    );
  })
  return (props: any) => (
    <Create {...props}>
        <SimpleForm>
            {comps.map((Comp: any) => Comp)}
        </SimpleForm>
    </Create>
  );
};

export default CreateGenerator;
