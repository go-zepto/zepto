import { Resource } from "../../types/schema";

const COMMON_TITLE_FIELD_NAMES = [
  'title',
  'name',
  'first_name',
  'full_name',
];

export const guessMainTitleField = (res: Resource): (string|null) => {
  if (res.list_endpoint.fields.length === 0) {
    return null
  }
  const listFieldNames = res.list_endpoint.fields.map(f => f.name);
  for (const cfn of COMMON_TITLE_FIELD_NAMES) {
    if (listFieldNames.includes(cfn)) {
      return cfn;
    }
  }
  return null;
}
