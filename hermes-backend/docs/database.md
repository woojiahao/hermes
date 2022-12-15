# ERD

```mermaid
erDiagram
  user {
    uuid id PK
    text username "unique"
    text email "unique"
    text password_hash
  }
  user_role {
    uuid user_id FK "references user.id"
    uuid role_id FK "references role.id"
  }
  role {
    uuid id PK
    text permissions
    text title
  }
  thread {
    uuid id PK
    bool is_published
    bool is_open
    text content
    date created_at
    uuid created_by FK "references user.id"
    date deleted_at
    uuid deleted_by FK "references user.id"
  }
  thread_tag {
    uuid thread_id FK "references thread.id"
    uuid tag_id FK "references tag.id"
  }
  tag {
    uuid id PK
    text content "unique"
    text hex_code
  }
  vote {
    uuid id PK
    uuid user_id FK "references user.id"
    uuid thread_id FK "references thread.id"
    bool is_upvote
  }
  comment {
    uuid id PK
    text content
    uuid created_by FK "references user.id"
    date created_at
    uuid thread_id FK "references thread.id"
  }

  thread ||--o{ thread_tag : has
  tag ||--o{ thread_tag : has
  user ||--o{ user_role : is
  role ||--o{ user_role : is
  user ||--o{ vote : cast
  thread ||--o{ vote : has
  user ||--o{ thread : create
  thread ||--o{ comment : has
  user ||--o{ comment : leave
```
