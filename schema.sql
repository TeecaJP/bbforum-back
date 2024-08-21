CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    image VARCHAR(255),
    favorite_team VARCHAR(255),
    location VARCHAR(255),
    type VARCHAR(50),
    birthday DATE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_delete BOOLEAN DEFAULT FALSE,
    is_ban BOOLEAN DEFAULT FALSE,
    is_official BOOLEAN DEFAULT FALSE
);

CREATE TABLE boards (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tags (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE posts (
    id VARCHAR(36) PRIMARY KEY,
    board_id VARCHAR(36),
    user_id VARCHAR(36),
    content TEXT NOT NULL,
    reply_to VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (board_id) REFERENCES boards(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (reply_to) REFERENCES posts(id)
);

CREATE TABLE board_tags (
    board_id VARCHAR(36),
    tag_id VARCHAR(36),
    PRIMARY KEY (board_id, tag_id),
    FOREIGN KEY (board_id) REFERENCES boards(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

CREATE TABLE post_tags (
    post_id VARCHAR(36),
    tag_id VARCHAR(36),
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
