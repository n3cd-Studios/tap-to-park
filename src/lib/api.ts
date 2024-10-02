
// The params that can be passed to a get function
export interface GetParams {
    route: string;
    method?: "GET" | "POST";
    params?: GetPairs;
    headers?: GetPairs;
    body?: any;
}

// A helper object for key-value string pairs
interface GetPairs { 
    [key: string]: string
}

// Just a helper that does all of the heavy lifting, notice how it is not exported
const getHelper = <T, R>({ route, params, body, headers, method = "GET" }: GetParams, defaultValue: R): Promise<T | R> => 
    fetch(`http://localhost:8080/api/${route}${params ? `?${new URLSearchParams(params).toString()}` : ""}`, { method, headers, body: JSON.stringify(body) })
        .then(r => r.json() as T)
        .catch(_ => defaultValue);

/**
 * A wrapper that simplifies the fetch operation and returns a default 
 *
 * @param params - The params of the function
 * @param defaultValue - The default value to return if this fails
 */
export const getWithDefault = <T>(params: GetParams, defaultValue: T) => getHelper<T, T>(params, defaultValue);

/**
 * A wrapper that simplifies the fetch operation and returns null if fails 
 *
 * @param params - The params of the function
 */
export const get = <T>(params: GetParams) => getHelper<T, null>(params, null);
