import React from "react";
import { Resource as RAResource } from "react-admin";
import { ResourceGenerator } from "../../../types/generators";
import { Resource, Schema } from "../../../types/schema";
import CreateGenerator from "./resGens/Create";
import EditGenerator from "./resGens/Edit";
import ListGenerator from "./resGens/List";

type GenerateResourceCompType = (schema: Schema, res: Resource) => any;

interface ResourceGenerators {
  list: ResourceGenerator
  create: ResourceGenerator
  edit: ResourceGenerator
};

const Generator: ResourceGenerators = {
  list: ListGenerator,
  create: CreateGenerator,
  edit: EditGenerator
};

export const generateResourceComp: GenerateResourceCompType = (schema: Schema, res: Resource) => {
  const list = Generator.list(schema, res);
  const create = Generator.create(schema, res);
  const edit = Generator.edit(schema, res);
  return (
    <RAResource
      key={res.name}
      name={res.endpoint.toLowerCase()}
      list={list}
      create={create}
      edit={edit}
    />
  );
}
