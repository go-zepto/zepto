import React from "react";
import { ComponentGenerator, FieldProps } from "../../../types/generators";
import { Field } from "../../../types/schema";
import { DatetimeFieldGenerator, DatetimeInputGenerator } from "./compGens/Datetime";
import { NumberFieldGenerator, NumberInputGenerator } from "./compGens/Number";
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
  }
};

export const generateFieldCompFromField = (field: Field): React.FC<FieldProps> => DEFAULT_FIELD_GENERATORS[field.type].fieldCompGen(field);
export const generateInputCompFromField = (field: Field): React.FC<FieldProps> => DEFAULT_FIELD_GENERATORS[field.type].inputCompGen(field);
