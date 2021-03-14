import { DateField, DateTimeInput } from 'react-admin';
import { ComponentGenerator } from '../../../../../types/generators';
import { Field } from '../../../../../types/schema';


export const DatetimeFieldGenerator: ComponentGenerator = (f: Field) => (props: any) => {
  return <DateField {...props} />;
}

export const DatetimeInputGenerator: ComponentGenerator = (f: Field) => (props: any) => (
  <DateTimeInput {...props} />
);
