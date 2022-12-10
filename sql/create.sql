CREATE TABLE tag
(
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "content" TEXT UNIQUE NOT NULL,
    hex_code  TEXT        NOT NULL
);

CREATE TYPE "role" AS ENUM ('ADMIN', 'USER');

CREATE TABLE "user"
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    "role"        "role"      NOT NULL
);

-- No updated_by as only the user who posted the thread can edit it
CREATE TABLE thread
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    is_published BOOL NOT NULL    DEFAULT TRUE,
    is_open      BOOL NOT NULL    DEFAULT TRUE,
    title        TEXT NOT NULL,
    "content"    TEXT NOT NULL,
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
    FOREIGN KEY (user_id) REFERENCES "user" (id),
    FOREIGN KEY (thread_id) REFERENCES "thread" (id)
);

CREATE TABLE thread_tag
(
    thread_id UUID NOT NULL,
    tag_id    UUID NOT NULL,
    FOREIGN KEY (thread_id) REFERENCES thread (id),
    FOREIGN KEY (tag_id) REFERENCES tag (id)
);

CREATE TABLE comment
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "content"  TEXT NOT NULL,
    created_at DATE NOT NULL    DEFAULT now(),
    created_by UUID NOT NULL,
    thread_id  UUID NOT NULL,
    deleted_at DATE NULL,
    deleted_by UUID NULL,
    FOREIGN KEY (thread_id) REFERENCES thread (id),
    FOREIGN KEY (created_by) REFERENCES "user" (id),
    FOREIGN KEY (deleted_by) REFERENCES "user" (id)
);