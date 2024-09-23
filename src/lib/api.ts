export interface GetParams {
    route: string;
    method?: "GET" | "POST";
    params?: URLSearchParams;
    body?: any;
}

export const get = <T>({ route, params, body, method = "GET" }: GetParams): Promise<T | null> => 
    fetch(`http://localhost:8080/api/${route}${params ? `?${params.toString()}` : ""}`, { method, body: JSON.stringify(body) })
        .then(r => r.json() as T)
        .catch(_ => null);