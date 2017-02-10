var baseURL = '';

if (module.hot) {
    baseURL = 'http://localhost:3000';
}

export function getAPI(api) {
    return baseURL + api;
}
