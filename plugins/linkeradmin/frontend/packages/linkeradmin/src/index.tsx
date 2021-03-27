import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './App';
import { fetchUtils } from 'ra-core';
import { Resource, Schema, ComponentGenerator} from './types';
import Generator from './core/generator';


interface LinkerAdminOptions {
  schemaPath?: string;
  target: string;
  defaultRowClick?: 'edit' | 'show';
}

export default class LinkerAdmin {
  schemaPath: string;
  target: string;
  defaultRowClick: 'edit' | 'show';

  constructor({
    schemaPath = '/admin/_schema',
    target = 'root',
    defaultRowClick,
  }: LinkerAdminOptions) {
    this.schemaPath = schemaPath;
    this.target = target;
    this.defaultRowClick = defaultRowClick ?? 'edit';
  }

  registerComponentGenerator(name: string, compGen: ComponentGenerator) {
    // this.fieldGenerators[name] = compGen;
  }

  async init() {
    const targetEl = document.getElementById(this.target);
    const res = await fetchUtils.fetchJson(this.schemaPath);
    try {
      const schema: Schema = res.json;
      const gen: Generator = new Generator({
        schema,
        defaultRowClick: this.defaultRowClick,
      });
      const resComps = schema.admin.resources.map((r: Resource) =>  gen.generateResourceComp(r));
      const AdminApp = () => (
        <App resComponents={resComps} />
      );
      ReactDOM.render(<AdminApp />, targetEl);
    } catch (error) {
      // TODO: Handle this error
      console.error(error);
    }
  }
}
