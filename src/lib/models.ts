
export interface Coords {
    latitude: number;
    longitude: number;
};

export interface Spot {
    guid: string;
    name: string;
    coords: Coords;
}

export interface User {
    guid: string;
    email: string;
}

export interface Organization {
    name: string;
    spots: Spot[]
}