import type { Coords } from "./models";

export enum IconType {
    SUCCESS = "success",
    WARNING = "warning",
    ERROR = "error",
    NONE = "none",
}

export enum ButtonType {
    POSITIVE = "positive",
    NEGATIVE = "negative",
    CAUTION = "caution",
    DEFAULT = "default",
}

export const promisifyGeolocation = (): Promise<Coords> =>
    new Promise((res, rej) => navigator.geolocation ? navigator.geolocation.getCurrentPosition(({ coords: { latitude, longitude } }) =>
      res({ latitude, longitude })) : rej(null));
