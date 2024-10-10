use `airplanes`;

DROP TABLE IF EXISTS user_tokens_versions;

CREATE TABLE user_tokens_versions (
    user_id INT NOT NULL UNIQUE,
    version INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);