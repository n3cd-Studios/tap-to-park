
export interface Coords {
    latitude: number;
    longitude: number;
};

export interface Spot {
    name: string;
    coords: Coords;
}

export interface User {
    uuid: string;
    email: string;
}

export interface Organization {
    name: string;
}