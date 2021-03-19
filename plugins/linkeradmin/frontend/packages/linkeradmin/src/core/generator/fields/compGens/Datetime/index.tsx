import React from 'react';
import { DateField, DateTimeInput } from 'react-admin';
import { ComponentGeneratorFunc, Field, Schema } from '../../../../../types';


export const DatetimeFieldGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => {
  return <DateField {...props} />;
}

export const DatetimeInputGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => (
  <DateTimeInput {...props} />
);
