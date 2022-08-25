https://www.youtube.com/watch?v=RroLKn54FzE&t=706s

<!-- Sample queries -->

```gql
query getUsers {
  users {
    id
  }
}

mutation createUser {
  createUser(input: { name: "user1" }) {
    id
  }
}

query getTodos {
  todos {
    text
    user {
      name
    }
  }
}

mutation createTodo {
  createTodo(input: { text: "todo 1", userId: "user1" }) {
    id
    text
    user {
      id
    }
  }
}
```
