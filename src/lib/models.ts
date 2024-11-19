import { DayOfWeek } from "./lang";

export interface Coords {
    latitude: number;
    longitude: number;
};

export type PartialSpot = Pick<Spot, "guid" | "coords" | "timeLeft">;

export interface Spot {
    guid: string;
    coords: Coords;
    name: string;
    maxHours: number;
    timeLeft?: number;
    handicap: boolean;
    price?: number;
    reservation?: Reservation;
    table: Pricing;
}

export enum UserRole {
    USER = 0,
    ADMIN
}

export interface User {
    guid: string;
    email: string;
    role: UserRole;
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
