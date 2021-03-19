/* eslint-disable react/jsx-pascal-case */
// import * as React from 'react';
import React from 'react';
import { Admin } from 'react-admin';
import { LinkerDataProvider } from './dataProvider';


export const App = (props: any) => {
  return (
    <Admin dataProvider={LinkerDataProvider('http://localhost:8000/api')}>
      {props.resComponents}
    </Admin>
  );
};
