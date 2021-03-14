import React from "react";
import { Field, Resource } from "./schema";

export interface FieldProps {
  source: string;
}

export type ComponentGenerator = (field: Field) => React.FC<FieldProps>;
export type ResourceGenerator = (resource: Resource) => React.FC;
