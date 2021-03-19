import React from "react";
import { ComponentGeneratorFunc } from "../../../types/generators";
import { DatetimeFieldGenerator, DatetimeInputGenerator } from "./compGens/Datetime";
import { NumberFieldGenerator, NumberInputGenerator } from "./compGens/Number";
import { ReferenceFieldGenerator, ReferenceInputGenerator } from "./compGens/Reference";
import { ReferenceListFieldGenerator, ReferenceListInputGenerator } from "./compGens/ReferenceList";
import { TextFieldGenerator, TextInputGenerator } from "./compGens/Text";

interface FieldGenerator {
  fieldCompGen: ComponentGeneratorFunc
  inputCompGen: ComponentGeneratorFunc
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
  },
  "reference_list": {
    fieldCompGen: ReferenceListFieldGenerator,
    inputCompGen: ReferenceListInputGenerator,
  }
};
