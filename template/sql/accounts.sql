# 账号表
create table accounts
(
    id         int(10) unsigned    not null auto_increment comment 'ID',
    username   varchar(32)         not null default '' comment '昵称',
    nickname   varchar(32)         not null default '' comment '昵称',
    password   varchar(32)         not null default '' comment '密码',
    avatar     varchar(255)        not null default '' comment '头像',
    email      varchar(64)         not null default '' comment '邮箱',
    phone      varchar(16)         not null default '' comment '手机号',
    status     tinyint(1) unsigned not null default '2' comment '状态：1-禁用，2-启用',
    remark     varchar(255)        not null default '' comment '备注',
    login_at   datetime                     default null comment '最后登录时间',
    created_at datetime            not null default current_timestamp comment '创建时间',
    updated_at datetime            not null default current_timestamp comment '更新时间',
    primary key (id),
    unique key (username),
    key (email),
    key (created_at),
    key (updated_at)
) engine = InnoDB
  auto_increment = 10001
  default charset = utf8mb4 comment '账号表';
