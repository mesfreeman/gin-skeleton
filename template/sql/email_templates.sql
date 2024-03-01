# 邮件模板表
create table email_templates
(
    id         int unsigned not null auto_increment comment 'ID',
    subject    varchar(64)  not null default '' comment '邮件主题',
    content    text         not null comment '邮件内容',
    slug       varchar(32)  not null default '' comment '标识',
    remark     varchar(255) not null default '' comment '备注',
    created_at datetime     not null default current_timestamp comment '创建时间',
    updated_at datetime     not null default current_timestamp comment '更新时间',
    primary key (id),
    unique key (slug),
    key (updated_at)
) engine = InnoDB
  auto_increment = 1
  default charset = utf8mb4 comment '邮件模板表';

