use `airplanes`;

DROP TABLE IF EXISTS user_sessions;

CREATE TABLE `user_sessions` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_id` int(11) NOT NULL,
    `login_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `logout_at` timestamp,
    `invalid_logout_reason` varchar(255) DEFAULT NULL,
    `crash_reason_type` int(11) DEFAULT 0,
    
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`),
    CONSTRAINT `user_sessions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
)