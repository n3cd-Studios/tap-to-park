
export type Point = [number, number];
export class Region {
    lower: Point;
    upper: Point;

    constructor() {
        this.lower = [0, 0];
        this.upper = [0, 0];
    }

    size = () => Math.sqrt(Math.pow(this.upper[0] - this.lower[0], 2) + Math.pow(this.upper[1] - this.lower[1], 2));
    in = ([x, y]: Point) => this.lower[0] <= x && x <= this.upper[0] && this.lower[1] <= y && y <= this.upper[1];

};
