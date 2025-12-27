CREATE TABLE users(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

-- private - 1 to 1
-- group - N to N(n writes for n)
-- channel - N to M(N writes for M(subs))
CREATE TYPE CHAT_TYPES AS ENUM ('private', 'group', 'channel');

CREATE TABLE chats(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    type CHAT_TYPES NOT NULL,
    address VARCHAR(255) UNIQUE NULLS NOT DISTINCT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE chat_members(
    role VARCHAR(32) NOT NULL DEFAULT 'user',

    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    joined_at TIMESTAMP DEFAULT current_timestamp,
    is_banned BOOLEAN DEFAULT FALSE,

    last_read_msg_id BIGINT DEFAULT NULL,

    PRIMARY KEY (chat_id, user_id)
);
CREATE INDEX idx_chat_members_user_id ON chat_members(user_id);

CREATE TABLE msgs(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    author_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
); 
CREATE INDEX idx_msgs_chat_id_id ON msgs(chat_id, id);