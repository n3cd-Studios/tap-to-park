
export interface GetParams {
    route: string;
    method?: "GET" | "POST";
    params?: URLSearchParams;
    headers?: GetHeaderParams;
    body?: any;
}

interface GetHeaderParams { 
    [key: string]: string
}

const getHelper = <T, R>({ route, params, body, headers, method = "GET" }: GetParams, defaultValue: R): Promise<T | R> => 
    fetch(`http://localhost:8080/api/${route}${params ? `?${params.toString()}` : ""}`, { method, headers, body: JSON.stringify(body) })
        .then(r => r.json() as T)
        .catch(_ => defaultValue);

export const getWithDefault = <T>(params: GetParams, defaultValue: T) => getHelper<T, T>(params, defaultValue);
export const get = <T>(params: GetParams) => getHelper<T, null>(params, null);
