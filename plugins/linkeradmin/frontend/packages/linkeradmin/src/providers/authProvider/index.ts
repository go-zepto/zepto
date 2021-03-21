import { fetchUtils, AuthProvider } from 'ra-core';

export const ZeptoAuthTokenProvider = (
  baseURL: string,
  httpClient = fetchUtils.fetchJson,
): AuthProvider => ({
    login: params  => {
      return httpClient(`${baseURL}`, {
        method: 'POST',
        body: JSON.stringify(params),
      }).then(({ json }) => {
        console.log(json);
        localStorage.setItem("auth_token", json.token.value);
      });
    },
    checkError: error => {
      const status = error.status;
      if (status === 401 || status === 403) {
          localStorage.removeItem('auth');
          return Promise.reject();
      }
      return Promise.resolve();
    },
    checkAuth: params => {
      return localStorage.getItem("auth_token") !== null ? Promise.resolve() : Promise.reject();
    },
    logout: () => {
      return httpClient(`${baseURL}/logout`, {
        method: 'POST',
      }).then(() => {
        localStorage.removeItem('auth_token');
      });
    },
    // @ts-ignore
    getIdentity: () => Promise.resolve({}),
    getPermissions: params => Promise.resolve(),
});
