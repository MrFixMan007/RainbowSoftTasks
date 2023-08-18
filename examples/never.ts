type Bird = {kind: 'bird', legs: number; wings: 2;};
type Dog = {kind: 'dog', legs: number; };
type Fish = {kind: 'fish', finds: number; };
type Cat = {kind: 'cat', legs: number; };

type Monkey = {kind: 'monkey', legs: number};

type Animals = Bird | Dog | Fish | Cat;

function animalAppendages(animal: Animals): number {
    switch(animal.kind){
        case "bird":
            return animal.wings + animal.legs;  
        case "dog":
            return animal.legs;
        case "fish":
            return animal.finds;
        case "cat":
            return animal.legs;
        default:
            let neverHappens : never = animal; //ругается если не прописали действия 
            return neverHappens;
        //но также ругается если вообще убрать default и не исп-ть это)
    }
}

function animalAppendages1(animal: Animals) {
    if(animal.kind === 'bird') return 1;
    else if (animal.kind === 'dog') return 0;
    else if (animal.kind === 'fish') return 2;
    const neverHappens : never = animal;
}
    function animalAppendages2(animal: Animals): number {
        switch(animal.kind){
            case "bird":
        return animal.wings + animal.legs;
            case "dog":
        return animal.legs;
            case "fish":
        return animal.finds;
        }
}