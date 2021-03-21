import { fetchUtils, DataProvider } from 'ra-core';
import { stringify } from 'query-string';

const VALID_WHERE_FILTER_REGEX = /(or_)?([_a-zA-Z0-9]+)_(eq|gt|gte|lt|lte|between|in|nin|like|nlike)/

// Linker Filter
type LinkerFilter = {
  skip?: number;
  limit?: number;
  where?: any;
  include?: {
    relation: string;
    where: any;
  }[]
};

// React-Admin filter
type RAFilter = {
  [key: string]: string;
};

const parseWhereFromFilter = (filter: RAFilter): any => {
  const where = {} as any;
  Object.keys(filter).forEach(f => {
    const match = VALID_WHERE_FILTER_REGEX.exec(f);
    if (match?.length === 4) {
      const boolOperator = match[1] === "or_" ? "or" : "and";
      const fieldName = match[2] as string;
      const filterType = match[3] as string;
      const item: any = {};
      item[fieldName] = {};
      item[fieldName][filterType] = filter[f];
      if (!where["and"]) {
        where["and"] = [{}];
      }
      if (!where["and"][0][boolOperator]) {
        where["and"][0][boolOperator] = [];
      }
      where["and"][0][boolOperator].push(item);
    }
  });
  return where;
};

export const LinkerDataProvider = (
  baseURL: string,
  httpClient = fetchUtils.fetchJson,
): DataProvider => ({
  getList: (resource, params) => {
    const where = parseWhereFromFilter(params.filter);
    const { page, perPage } = params.pagination;
    const linkerFilter: LinkerFilter = {
      skip: (page-1) * perPage,
      limit: perPage,
      where,
    };
    const query = {
      filter: JSON.stringify(linkerFilter),
    }
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
