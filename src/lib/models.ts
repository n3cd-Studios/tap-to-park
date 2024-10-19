import { DayOfWeek } from "./lang";

export interface Coords {
    latitude: number;
    longitude: number;
};

export interface Spot {
    guid: string;
    coords: Coords;
    name: string;
    price?: number;
    reservation?: Reservation;
    table: Pricing;
}

export interface User {
    guid: string;
    email: string;
}

export interface Session {
    guid: string;
	device: string;
    ip: string;
	expires: string;
	lastUsed: string;
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

export interface Reservation {
	guid: string;
	start: string;
	end: string;
	price: number;
	email: string;
}

export type Pricing = {
    [key in DayOfWeek]: number[]
}