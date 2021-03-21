/* eslint-disable react/jsx-pascal-case */
// import * as React from 'react';
import React from 'react';
import { Admin, fetchUtils } from 'react-admin';
import { LinkerDataProvider } from './providers/dataProvider';
import { ZeptoAuthTokenProvider } from './providers/authProvider';

const httpClient = (url: string, options: any = {}) => {
  if (!options.headers) {
      options.headers = new Headers({ Accept: 'application/json' });
  }
  const token = localStorage.getItem('auth_token');
  if (token !== null && token !== "") {
    options.headers.set('Authorization', `Bearer ${token}`);
  }
  return fetchUtils.fetchJson(url, options);
};

export const App = (props: any) => {
  return (
    <Admin
      dataProvider={LinkerDataProvider('/api', httpClient)}
      authProvider={ZeptoAuthTokenProvider('/auth', httpClient)}
    >
      {props.resComponents}
    </Admin>
  );
};
