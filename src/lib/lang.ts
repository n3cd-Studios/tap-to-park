
export const pluralize = (val: number, singular: string, plural: string = `${singular}s`) => `${val} ${val == 1 ? singular : plural}`;
export const properNoun = (val: string) => `${val[0].toUpperCase()}${val.substring(1).toLowerCase()}`;

export const Formats = {
    USDollar: new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
    })
}

export enum DayOfWeek {
    sunday = "sunday",
    monday = "monday",
    tuesday = "tuesday",
    wednesday = "wednesday",
    thursday = "thursday",
    friday = "friday",
    saturday = "saturday"
};

export const daysOfWeek: DayOfWeek[] = [  DayOfWeek.sunday, DayOfWeek.monday, DayOfWeek.tuesday, DayOfWeek.wednesday, DayOfWeek.thursday, DayOfWeek.friday, DayOfWeek.saturday];
