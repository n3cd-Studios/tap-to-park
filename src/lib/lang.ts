
export const pluralize = (val: number, singular: string, plural: string = `${singular}s`) => `${val} ${val == 1 ? singular : plural}`; 

export const Formats = {
    USDollar: new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
    })
} 