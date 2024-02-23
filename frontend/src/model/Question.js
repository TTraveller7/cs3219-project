export class Question {
    constructor(id, name, description) {
        this.id = id;
        this.name = name;
        this.description = description;
    }
}

export function createQuestion(question) {
    return new Question(question.ID, question.name, question.description)
}