import { TextField, SelectInput } from 'react-admin';
import { ComponentGeneratorFunc, Field, Schema } from '../../../../../types';


export const SelectFieldGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => {
  return (props: any) => {
    return <TextField {...props} />;
  };
}

export const SelectInputGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => (
  <SelectInput {...props} />
);
