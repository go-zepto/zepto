import React from "react";
import { ComponentGenerator, FieldProps } from "../../../types/generators";
import { Field, Schema } from "../../../types/schema";
import { DatetimeFieldGenerator, DatetimeInputGenerator } from "./compGens/Datetime";
import { NumberFieldGenerator, NumberInputGenerator } from "./compGens/Number";
import { ReferenceFieldGenerator, ReferenceInputGenerator } from "./compGens/Reference";
import { TextFieldGenerator, TextInputGenerator } from "./compGens/Text";

interface FieldGenerator {
  fieldCompGen: ComponentGenerator
  inputCompGen: ComponentGenerator
};


type FieldGenerators = {
  [key: string]: FieldGenerator;
};

export const DEFAULT_FIELD_GENERATORS: FieldGenerators = {
  "text": {
    fieldCompGen: TextFieldGenerator,
    inputCompGen: TextInputGenerator,
  },
  "number": {
    fieldCompGen: NumberFieldGenerator,
    inputCompGen: NumberInputGenerator,
  },
  "datetime": {
    fieldCompGen: DatetimeFieldGenerator,
    inputCompGen: DatetimeInputGenerator,
  },
  "reference": {
    fieldCompGen: ReferenceFieldGenerator,
    inputCompGen: ReferenceInputGenerator,
  }
};

export const generateFieldCompFromField = (s: Schema, f: Field): React.FC<FieldProps> => DEFAULT_FIELD_GENERATORS[f.type].fieldCompGen(s, f);
export const generateInputCompFromField = (s: Schema, f: Field): React.FC<FieldProps> => DEFAULT_FIELD_GENERATORS[f.type].inputCompGen(s, f);
