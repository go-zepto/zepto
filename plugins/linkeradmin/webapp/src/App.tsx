/* eslint-disable react/jsx-pascal-case */
import React from 'react';
import { Admin } from 'react-admin';
import { LinkerDataProvider } from './dataProvider';


export const App = (props: any) => {
  return (
    <Admin dataProvider={LinkerDataProvider('http://localhost:8000/api')} layout={props.layoutComponent}>
      {props.resComponents}
    </Admin>
  );
};
