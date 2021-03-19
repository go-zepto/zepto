import { TextInput, TextField } from 'react-admin';
import { ComponentGeneratorFunc, Field, Schema } from '../../../../../types';


export const TextFieldGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => {
  return (props: any) => {
    return <TextField {...props} />;
  };
}

export const TextInputGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => (
  <TextInput {...props} />
);
