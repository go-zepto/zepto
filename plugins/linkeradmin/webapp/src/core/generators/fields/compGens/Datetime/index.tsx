import { DateField, DateTimeInput } from 'react-admin';
import { ComponentGenerator } from '../../../../../types/generators';
import { Field, Schema } from '../../../../../types/schema';


export const DatetimeFieldGenerator: ComponentGenerator = (s: Schema, f: Field) => (props: any) => {
  return <DateField {...props} />;
}

export const DatetimeInputGenerator: ComponentGenerator = (s: Schema, f: Field) => (props: any) => (
  <DateTimeInput {...props} />
);
