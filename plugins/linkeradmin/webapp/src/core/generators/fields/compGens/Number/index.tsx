import React from 'react';
import { NumberField, NumberInput } from 'react-admin';
import { ComponentGenerator } from '../../../../../types/generators';
import { Field, Schema } from '../../../../../types/schema';


export const NumberFieldGenerator: ComponentGenerator = (s: Schema, f: Field) => (props: any) => {
  return <NumberField {...props} />;
}

export const NumberInputGenerator: ComponentGenerator = (s: Schema, f: Field) => (props: any) => (
  <NumberInput {...props} />
);
