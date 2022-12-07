CREATE TABLE tag
(
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "content" TEXT UNIQUE NOT NULL,
    hex_code  TEXT        NOT NULL
);

CREATE TABLE role
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT NOT NULL,
    permissions TEXT NOT NULL
);

CREATE TABLE "user"
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      TEXT UNIQUE NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL
);

CREATE TABLE user_role
(
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id),
    FOREIGN KEY (role_id) REFERENCES role (id)
);

CREATE TABLE thread
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    is_published BOOL NOT NULL    DEFAULT FALSE,
    is_open      BOOL NOT NULL    DEFAULT TRUE,
    "content"    TEXT NOT NULL,
    created_at   DATE NOT NULL    DEFAULT now(),
    created_by   UUID NOT NULL,
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
    FOREIGN KEY (thread_id) REFERENCES thread (id),
    FOREIGN KEY (created_by) REFERENCES "user" (id)
);

INSERT INTO role (title, permissions)
VALUES ('ADMIN', 'RWA');
INSERT INTO role (title, permissions)
VALUES ('USER', 'RW');
