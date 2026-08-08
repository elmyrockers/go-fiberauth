CREATE TABLE IF NOT EXISTS users (
    id                          CHAR(36)     NOT NULL PRIMARY KEY,
    name                        VARCHAR(255) NULL DEFAULT NULL,
    email                       VARCHAR(255) NOT NULL,
    email_verified_at           DATETIME(6)  NULL,
    password                    VARCHAR(255) NOT NULL,
    remember_token              VARCHAR(100) NULL,
    two_factor_secret           TEXT         NULL,
    two_factor_recovery_codes   TEXT         NULL,
    two_factor_confirmed_at     DATETIME(6)  NULL,
    created_at                  DATETIME(6)  NOT NULL,
    updated_at                  DATETIME(6)  NOT NULL,
    UNIQUE KEY uniq_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    user_id     CHAR(36)     NOT NULL,
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  DATETIME(6)  NOT NULL,
    used_at     DATETIME(6)  NULL,
    created_at  DATETIME(6)  NOT NULL,
    KEY idx_password_reset_tokens_user_id (user_id),
    KEY idx_password_reset_tokens_hash (token_hash),
    CONSTRAINT fk_password_reset_tokens_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;