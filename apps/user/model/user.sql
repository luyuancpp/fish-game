CREATE TABLE `user` (
                        `uid` BIGINT NOT NULL AUTO_INCREMENT,
                        `nickname` VARCHAR(50) NOT NULL DEFAULT '',
                        `gold` INT NOT NULL DEFAULT 0,
                        PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
