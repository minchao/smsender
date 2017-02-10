var baseURL = '';

if (module.hot) {
    baseURL = 'http://localhost:8080';
}

export function getAPI(api) {
    return baseURL + api;
}
