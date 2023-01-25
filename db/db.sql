DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS user_forum CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    nickname citext COLLATE "ucs_basic" PRIMARY KEY NOT NULL,
    fullname VARCHAR                                NOT NULL,
    about    VARCHAR                                NOT NULL,
    email    citext UNIQUE                          NOT NULL
);

CREATE TABLE IF NOT EXISTS forums
(
    title           VARCHAR                                                NOT NULL,
    author_nickname citext REFERENCES users (nickname) ON DELETE NO ACTION NOT NULL,
    slug            citext PRIMARY KEY                                     NOT NULL,
    posts           INTEGER DEFAULT 0                                      NOT NULL,
    threads         INTEGER DEFAULT 0                                      NOT NULL
);

CREATE TABLE IF NOT EXISTS user_forum
(
    nickname   citext REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    forum_slug citext REFERENCES forums (slug) ON DELETE CASCADE    NOT NULL,
    PRIMARY KEY (nickname, forum_slug)
);

CREATE TABLE IF NOT EXISTS threads
(
    id              BIGSERIAL PRIMARY KEY                                  NOT NULL,
    title           VARCHAR                                                NOT NULL,
    author_nickname citext REFERENCES users (nickname) ON DELETE NO ACTION NOT NULL,
    forum           citext REFERENCES forums (slug) ON DELETE CASCADE      NOT NULL,
    slug            citext,
    votes           INTEGER                  DEFAULT 0,
    message         VARCHAR                                                NOT NULL,
    created         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP     NOT NULL
);

CREATE TABLE IF NOT EXISTS votes
(
    nickname  citext REFERENCES users (nickname) ON DELETE NO ACTION NOT NULL,
    thread_id BIGSERIAL REFERENCES threads (id) ON DELETE CASCADE    NOT NULL,
    voice     SMALLINT CHECK ( voice BETWEEN -1 AND 1 )              NOT NULL,
    PRIMARY KEY (nickname, thread_id)
);

CREATE TABLE IF NOT EXISTS posts
(
    id              BIGSERIAL PRIMARY KEY                                  NOT NULL,
    author_nickname citext REFERENCES users (nickname) ON DELETE NO ACTION NOT NULL,
    forum_slug      citext REFERENCES forums (slug) ON DELETE CASCADE      NOT NULL,
    message         text                                                   NOT NULL,
    thread_id       BIGSERIAL REFERENCES threads (id) ON DELETE CASCADE    NOT NULL,
    parent_id       BIGINT                   DEFAULT 0,
    is_edited       BOOLEAN                  DEFAULT FALSE                 NOT NULL,
    created         TIMESTAMP WITH TIME ZONE DEFAULT now()                 NOT NULL,
    path            BIGINT[]                 DEFAULT ARRAY []::BIGINT[]
);

CREATE OR REPLACE FUNCTION update_path_trigger() RETURNS TRIGGER AS
$$
BEGIN
    new.path = (SELECT path FROM posts WHERE id = new.parent_id) || new.id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_path_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path_trigger();

CREATE OR REPLACE FUNCTION insert_trigger_forum_posts() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums SET posts = posts + 1 WHERE slug = new.forum_slug;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_trigger_forum_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE insert_trigger_forum_posts();

CREATE OR REPLACE FUNCTION insert_trigger_forum_threads() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums SET threads = threads + 1 WHERE slug = new.forum;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_trigger_forum_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE insert_trigger_forum_threads();

CREATE OR REPLACE FUNCTION insert_trigger_thread_votes() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE threads SET votes = votes + new.voice WHERE id = new.thread_id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_trigger_thread_votes() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE threads SET votes = votes + new.voice - old.voice WHERE id = new.thread_id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_trigger_thread_votes
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_trigger_thread_votes();

CREATE TRIGGER update_trigger_thread_votes
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_trigger_thread_votes();

CREATE OR REPLACE FUNCTION update_trigger_user_forum_thread() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO user_forum (nickname, forum_slug)
    VALUES (new.author_nickname, new.forum)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_trigger_user_forum_posts() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO user_forum (nickname, forum_slug)
    VALUES (new.author_nickname, new.forum_slug)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_trigger_user_forum_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_trigger_user_forum_thread();

CREATE TRIGGER update_trigger_user_forum_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_trigger_user_forum_posts();

-- INDEXES

-- Forum
-- indexed based on selectivity
CREATE INDEX IF NOT EXISTS index_forum_nickname ON forums (author_nickname);

-- Users
-- indexed based on selectivity
CREATE INDEX IF NOT EXISTS index_users_email ON users (email);

-- Threads
-- index based on forum/{slug}/threads
CREATE INDEX IF NOT EXISTS index_thread_forum_nickname ON threads (forum, author_nickname);
-- indexed based on selectivity
CREATE INDEX IF NOT EXISTS index_threads_slug on threads (slug);
-- index based on forum/{slug}/threads with since option
CREATE INDEX IF NOT EXISTS index_threads_forum_created ON threads (forum, created);

-- Posts
-- index based on thread/{slug_or_id}/posts
CREATE INDEX IF NOT EXISTS index_posts_thread_id on posts (thread_id);
-- index based on flat sorting with since option
CREATE INDEX IF NOT EXISTS index_posts_flat on posts (thread_id, created, id);
-- index based on tree sorting
CREATE INDEX IF NOT EXISTS index_posts_tree on posts (thread_id, path);
-- index based on parent_tree sorting
CREATE INDEX IF NOT EXISTS index_posts_parent_tree on posts ((path[1]), path);

VACUUM ANALYZE;
