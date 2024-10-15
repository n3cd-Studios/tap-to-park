

// The params that can be passed to a get function
export interface GetParams {
    route: string;
    method?: "GET" | "POST" | "PUT" | "DELETE";
    params?: GetPairs;
    headers?: GetPairs;
    body?: any;
}

// A helper object for key-value string pairs
interface GetPairs { 
    [key: string]: string
}

// A helper type to add paginated items
interface AsPaginated<T> { 
    items: T[], 
    page: number, 
    pages: number
};

// Paginate a table of a certain datatype
export class Paginator<T> {

    private page: number = 0;
    private pages: number = 0;
    private update?: (items: typeof this.items) => any;
    items: T[] = [];

    constructor(private params: GetParams, private pageSize: number = 10) {}

    hasNext = () => this.page < this.pages;
    next = () => {
        this.page = Math.min(++this.page, this.pages)
        return this.load();
    };

    hasLast = () => this.page > 0;
    last = () => {
        this.page = Math.max(--this.page, 0);
        return this.load();
    };

    set = (page: number) => {
        this.page = Math.min(Math.max(page, 0), this.pages);
        return this.load();
    };

    subscribe = (cb: typeof this.update) => this.update = cb;

    load = async () => {
        const payload = await get<AsPaginated<T>>({ ...this.params, params: { ...this.params.params, page: `${this.page}`, size: `${this.pageSize}` }});
        if (payload) {
            this.items = payload.items;
            this.pages = payload.pages;
        }
        if (this.update) this.update(this.items);
        return this.items;
    };

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
