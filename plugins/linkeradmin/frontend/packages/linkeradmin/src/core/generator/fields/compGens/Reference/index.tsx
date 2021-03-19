import React from 'react';
import { TextField, ReferenceInput, SelectInput, AutocompleteInput, ReferenceField } from 'react-admin';
import { ComponentGeneratorFunc, Field, Schema } from '../../../../../types';
import { guessMainTitleField } from '../../../../utils/field';


export const ReferenceFieldGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => {
  const resource = s.resources.find(r => r.name === f.options["ref_resource"]);
  if (!resource) {
    console.error(`[ReferenceInput] Resource not found "${resource}"`);
    return null;
  }
  const optTextFieldName = f.options["text_field_name"];
  const textFieldName = optTextFieldName && optTextFieldName !== "" ? optTextFieldName : guessMainTitleField(resource);
  return (
  <ReferenceField {...props}  reference={resource?.endpoint}>
    <TextField source={textFieldName} />
  </ReferenceField>  
  );
}

export const ReferenceInputGenerator: ComponentGeneratorFunc = (s: Schema, f: Field) => (props: any) => {
  const resource = s.resources.find(r => r.name === f.options["ref_resource"]);
  if (!resource) {
    console.error(`[ReferenceInput] Resource not found "${resource}"`);
    return null;
  }
  const optionText = f.options.option_text_field;
  const { autocomplete } = f.options;
  const filterToQuery = (q: string) => {
    if (q === "") {
      return {};
    }
    const filter: any = {};
    const sf = autocomplete.searchable_fields;
    const guessedTitle = guessMainTitleField(resource);
    const searchableFields = sf.length > 0 ? sf.length : (
      guessedTitle != null ? [guessedTitle] : []
    );
    searchableFields.forEach((f: string) => {
      filter[`or_${f}_like`] = `%${q}%`;
    })
    return filter;
  };
  return (
    <ReferenceInput {...props} reference={resource?.endpoint} filterToQuery={filterToQuery}>
      {
        (autocomplete && autocomplete.enabled) ? (
          <AutocompleteInput optionText={optionText} />
          ) : (
          <SelectInput optionText={optionText} />
        )
      }
    </ReferenceInput>
  )
}
