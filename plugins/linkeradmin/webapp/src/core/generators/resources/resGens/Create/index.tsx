import React from "react";
import { Create, SimpleForm } from "react-admin";
import { ResourceGenerator } from "../../../../../types/generators";
import { Field, Resource } from "../../../../../types/schema";
import { generateInputCompFromField } from "../../../fields";


const CreateGenerator: ResourceGenerator =  (r: Resource) => {
  const fields = r.create_inputs;
  const comps = fields.map((f: Field) => {
    const Comp = generateInputCompFromField(f);
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
