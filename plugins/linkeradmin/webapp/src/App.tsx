/* eslint-disable react/jsx-pascal-case */
import React, { useEffect, useState } from 'react';
import { Admin } from 'react-admin';
import { LinkerDataProvider } from './dataProvider';
import { fetchUtils } from 'ra-core';
import { Schema, Resource } from './types/schema';
import { generateResourceComp } from './core/generators/resources';



export const App = () => {
  const [resources, setResources] = useState<Resource[]>([]);
  useEffect(() => {
    fetchUtils.fetchJson('/admin/_schema').then(res => {
      const schema: Schema = res.json;
      const resComponents = schema.resources.map((r: Resource) =>  generateResourceComp(schema, r));
      if (resources.length === 0) {
        setResources(resComponents);
      }
    })
  }, [resources.length, setResources]);
  if (resources.length === 0) {
    return null;
  }
  return (
    <Admin dataProvider={LinkerDataProvider('http://localhost:8000/api')}>
      {resources}
    </Admin>
  );
};
