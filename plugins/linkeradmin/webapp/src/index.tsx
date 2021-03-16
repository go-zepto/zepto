import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './App';
import reportWebVitals from './reportWebVitals';
import { fetchUtils } from 'ra-core';
import { Resource, Schema } from './types/schema';
import { generateResourceComp } from './core/generators/resources';




const rootEl = document.getElementById('root');
const schemaPath = rootEl?.getAttribute("schema");
fetchUtils.fetchJson(schemaPath).then(res => {
  const schema: Schema = res.json;
  const resComponents = schema.resources.map((r: Resource) =>  generateResourceComp(schema, r));
  ReactDOM.render(
    <React.StrictMode>
      <App resComponents={resComponents} />
    </React.StrictMode>,
    rootEl,
  );
})


// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
