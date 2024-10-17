import { DayOfWeek } from "./lang";

export interface Coords {
    latitude: number;
    longitude: number;
};

export interface Spot {
    guid: string;
    name: string;
    coords: Coords;
    price?: number;
    table: Pricing;
}

export interface User {
    guid: string;
    email: string;
}

export interface Organization {
    name: string;
    spots: Spot[]
}

export interface Invite {
    code: string,
    expiration: string,
    organization: number,
    createdBy: number,
    usedBy: number
}

export type Pricing = {
    [key in DayOfWeek]: number[]
}