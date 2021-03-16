import { TextInput, TextField } from 'react-admin';
import { ComponentGenerator } from '../../../../../types/generators';
import { Field, Schema } from '../../../../../types/schema';


export const TextFieldGenerator: ComponentGenerator = (s: Schema, f: Field) => {
  return (props: any) => {
    return <TextField {...props} />;
  };
}

export const TextInputGenerator: ComponentGenerator = (s: Schema, f: Field) => (props: any) => (
  <TextInput {...props} />
);
