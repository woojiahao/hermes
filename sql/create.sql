CREATE TABLE tag
(
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "content" TEXT UNIQUE NOT NULL CHECK (TRIM("content") != ''),
    hex_code  CHAR(7)     NOT NULL CHECK (hex_code SIMILAR TO '^#[A-Fa-f0-9]{6}$')
);

CREATE TYPE "role" AS ENUM ('ADMIN', 'USER');

CREATE TABLE "user"
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      TEXT UNIQUE NOT NULL CHECK (TRIM(username) SIMILAR TO '^[a-zA-Z]\w{2,}$'),
    password_hash TEXT        NOT NULL,
    "role"        "role"      NOT NULL
);

-- No updated_by as only the user who posted the thread can edit it
CREATE TABLE thread
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    is_published BOOL NOT NULL    DEFAULT TRUE,
    is_open      BOOL NOT NULL    DEFAULT TRUE,
    is_pinned    BOOL NOT NULL    DEFAULT FALSE,
    title        TEXT NOT NULL CHECK (LENGTH(TRIM(title)) >= 5),
    "content"    TEXT NOT NULL CHECK (LENGTH(TRIM("content")) >= 30),
    created_at   DATE NOT NULL    DEFAULT now(),
    created_by   UUID NOT NULL,
    updated_at   DATE NULL,
    deleted_at   DATE NULL,
    deleted_by   UUID NULL,
    FOREIGN KEY (created_by) REFERENCES "user" (id),
    FOREIGN KEY (deleted_by) REFERENCES "user" (id)
);

CREATE TABLE vote
(
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id   UUID NOT NULL,
    thread_id UUID NOT NULL,
    is_upvote BOOL NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    FOREIGN KEY (thread_id) REFERENCES "thread" (id) ON DELETE CASCADE,
    UNIQUE (user_id, thread_id)
);

CREATE TABLE thread_tag
(
    thread_id UUID NOT NULL,
    tag_id    UUID NOT NULL,
    FOREIGN KEY (thread_id) REFERENCES thread (id),
    FOREIGN KEY (tag_id) REFERENCES tag (id),
    UNIQUE (thread_id, tag_id)
);

CREATE TABLE comment
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "content"  TEXT NOT NULL CHECK(LENGTH(TRIM("content")) >= 5),
    created_at DATE NOT NULL    DEFAULT now(),
    created_by UUID NOT NULL,
    thread_id  UUID NOT NULL,
    deleted_at DATE NULL,
    deleted_by UUID NULL,
    FOREIGN KEY (thread_id) REFERENCES thread (id),
    FOREIGN KEY (created_by) REFERENCES "user" (id),
    FOREIGN KEY (deleted_by) REFERENCES "user" (id)
);
