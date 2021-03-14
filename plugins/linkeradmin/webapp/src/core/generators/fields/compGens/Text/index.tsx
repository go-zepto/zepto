import { TextInput, TextField } from 'react-admin';
import { ComponentGenerator } from '../../../../../types/generators';
import { Field } from '../../../../../types/schema';


export const TextFieldGenerator: ComponentGenerator = (f: Field) => (props: any) => {
  console.log("EI. PARE DE GERAR!!");
  return <TextField {...props} />;
}

export const TextInputGenerator: ComponentGenerator = (f: Field) => (props: any) => (
  <TextInput {...props} />
);
