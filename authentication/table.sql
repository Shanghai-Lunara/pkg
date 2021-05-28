CREATE DATABASE rbac
    CHARACTER SET utf8
    COLLATE utf8_general_ci;

use rbac;

create table accounts (
     id int(3) NOT NULL AUTO_INCREMENT,
     PRIMARY KEY(id),
     account varchar(128) default '' comment 'account',
     unique (account),
     password varchar(128) default '' comment 'password',
     createTime int(11) default 0 comment 'createTime',
     status tinyint(2) default 0 comment 'status'
);

alter table accounts add index idx1(`account`);

