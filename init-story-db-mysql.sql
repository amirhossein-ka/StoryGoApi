# CREATE DATABASE story CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE story_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; -- for testing only
GRANT ALL PRIVILEGES ON story_test.* TO 'story'@'%';
