CREATE DATABASE IF NOT EXISTS `year_end` CHARSET utf8;

USE `year_end`;

CREATE TABLE `users` (
    `id` VARCHAR(10) PRIMARY KEY,
    `password` VARCHAR(100)
);

CREATE TABLE `books` (
    `count` INT PRIMARY KEY  AUTO_INCREMENT,
    `name` VARCHAR(50),
    `id` VARCHAR(10),
    FOREIGN KEY (`id`) REFERENCES  `users`(`id`)
) AUTO_INCREMENT = 1;

CREATE TABLE `card_records` (
    `count` INT PRIMARY KEY AUTO_INCREMENT,
    `id` VARCHAR(10),
    `trans_money` FLOAT,
    `location` VARCHAR(50),
    `method` VARCHAR(20),
    `date` DATETIME,
    FOREIGN KEY(`id`) REFERENCES `users`(`id`)
) AUTO_INCREMENT = 1;