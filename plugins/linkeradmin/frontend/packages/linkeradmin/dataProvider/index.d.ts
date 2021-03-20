import { fetchUtils, DataProvider } from 'ra-core';
export declare const LinkerDataProvider: (baseURL: string, httpClient?: (url: any, options?: fetchUtils.Options) => Promise<{
    status: number;
    headers: Headers;
    body: string;
    json: any;
}>) => DataProvider;
