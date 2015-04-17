drop database if exists test;
create database test;
use test;

drop table if exists account;
create table `account`(
    `guid` int unsigned not null auto_increment comment '账号ID',
    `name` char(32) not null comment '账号名称',
    primary key(`guid`),
    key(`name`)
)engine=innodb default character set=utf8 collate=utf8_general_ci;

drop table if exists role;
create table `role`(
    `guid` int unsigned not null comment '角色ID',
    `name` char(32) not null default ''  comment '玩家名字',
    `occupation` tinyint unsigned not null default 0 comment '职业',
    `level` smallint unsigned not null default 0 comment '等级',
    primary key(`guid`)
)engine=innodb default character set=utf8 collate=utf8_general_ci;

drop table if exists item;
create table `item`(
    `role_id` int unsigned not null default 0 comment '玩家id',
    `guid` smallint unsigned not null default 0 comment '物品guid',
    `item_id` smallint unsigned not null default 0 comment '物品id',
    `count` smallint unsigned not null default 0 comment '数量',
    key(`role_id`)
)engine=innodb default character set=utf8 collate=utf8_general_ci;

drop table if exists task;
create table `task`(
    `role_id` int unsigned not null default 0 comment '玩家id',
    `task_id` smallint unsigned not null default 0 comment '任务id',
    `progress` smallint unsigned not null default 0 comment '进度',
    key(`role_id`)
)engine=innodb default character set=utf8 collate=utf8_general_ci;
