import React from 'react';
import { NumberField, NumberInput } from 'react-admin';
import { ComponentGenerator } from '../../../../../types/generators';
import { Field } from '../../../../../types/schema';


export const NumberFieldGenerator: ComponentGenerator = (f: Field) => (props: any) => {
  return <NumberField {...props} />;
}

export const NumberInputGenerator: ComponentGenerator = (f: Field) => (props: any) => (
  <NumberInput {...props} />
);
