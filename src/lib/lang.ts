
export const pluralize = (val: number, singular: string, plural: string = `${singular}s`) => `${val} ${val == 1 ? singular : plural}`; 