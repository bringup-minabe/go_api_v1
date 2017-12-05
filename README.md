# Go Api Ver1

## Database

Mysql

### users

    CREATE TABLE `users` (
      `id` bigint(20) NOT NULL AUTO_INCREMENT,
      `uuid` varchar(45) DEFAULT NULL,
      `username` varchar(255) DEFAULT NULL,
      `password` text,
      `last_name` varchar(255) DEFAULT NULL,
      `first_name` varchar(255) DEFAULT NULL,
      `created_at` datetime DEFAULT NULL,
      `updated_at` datetime DEFAULT NULL,
      `deleted_at` datetime DEFAULT NULL,
      PRIMARY KEY (`id`),
      KEY `idx_user_uuid` (`uuid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

## Setting

### .env

バイナリファイルと同階層に配置

    APP_DB_USER = hoge
    APP_DB_PASS = hoge
    APP_DB_LOCATION = hoge
    APP_DB_DATABASE = hoge