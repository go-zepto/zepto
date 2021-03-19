import React from 'react';
import { NumberField, NumberInput } from 'react-admin';
import { ComponentGeneratorFunc, Field, Schema } from '../../../../../types';


export const NumberFieldGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => {
  return <NumberField {...props} />;
}

export const NumberInputGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => (
  <NumberInput {...props} />
);
