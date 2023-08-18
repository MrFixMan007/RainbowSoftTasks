let arr = ['str', 4, false]; //let arr: (string | number | boolean)[]
arr[0] = true;
arr[1] = true;
arr[2] = 5;

//let tup: [string | number, number, (boolean | undefined | null)?]
let tup : [string | number, number, boolean?]; //let tup: [string | number, number, (boolean | undefined)?]
tup = [4, 5, false]
tup[0] = 'str'
tup[1] = 4
tup[2] = undefined

let tup1: ['hah', 4];