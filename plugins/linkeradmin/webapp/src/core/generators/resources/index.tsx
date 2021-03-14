import React from "react";
import { Resource as RAResource } from "react-admin";
import { ResourceGenerator } from "../../../types/generators";
import { Resource } from "../../../types/schema";
import CreateGenerator from "./resGens/Create";
import EditGenerator from "./resGens/Edit";
import ListGenerator from "./resGens/List";

type GenerateResourceCompType = (res: Resource) => any;

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

export const generateResourceComp: GenerateResourceCompType = (res: Resource) => (
  <RAResource
    name={res.endpoint.toLowerCase()}
    list={Generator.list(res)}
    create={Generator.create(res)}
    edit={Generator.edit(res)}
  />
);
