import React from "react";
import { Field, Resource, Schema } from "./schema";

export interface FieldProps {
  source: string;
}

export type ComponentGenerator = (schema: Schema, field: Field) => React.FC<FieldProps>;
export type ResourceGenerator = (schema: Schema, resource: Resource) => React.FC;
export type LayoutComponentGenerator = (schema: Schema) => React.FC;
