import { fetchUtils, DataProvider } from 'ra-core';
import { stringify } from 'query-string';

export const LinkerDataProvider = (
  baseURL: string,
  httpClient = fetchUtils.fetchJson,
): DataProvider => ({
  getList: (resource, params) => {
    const query = {};
    const url = `${baseURL}/${resource}?${stringify(query)}`;
    return httpClient(url).then(({ headers, json }) => ({
      data: json.data,
      total: json.count,
    }));
  },
  getOne: (resource, params) => {
    const url = `${baseURL}/${resource}/${params.id}`;
    return httpClient(url).then(({ headers, json }) => ({
      data: json,
    }));
  },
  getMany: (resource, params) => {
    const query = {
      filter: JSON.stringify({ id: params.ids }),
  };
  const url = `${baseURL}/${resource}?${stringify(query)}`;
  return httpClient(url).then(({ json }) => ({ data: json.data }));
  },
  getManyReference: (resource, params) => {
    return Promise.resolve({} as any);
  },
  update: (resource, params) => {
    return httpClient(`${baseURL}/${resource}/${params.id}`, {
      method: 'PUT',
      body: JSON.stringify(params.data),
    }).then(({ json }) => ({ data: json }))
  },
  updateMany: (resource, params) => {
    return Promise.resolve({} as any);
  },
  create: (resource, params) => {
    return httpClient(`${baseURL}/${resource}`, {
      method: 'POST',
      body: JSON.stringify(params.data),
    }).then(({ json }) => ({ data: json }))
  },
  delete: (resource, params) => {
    return Promise.resolve({} as any);
  },
  deleteMany: (resource, params) => {
    return Promise.resolve({} as any);
  },
});
