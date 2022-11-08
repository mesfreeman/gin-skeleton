# 角色表
create table roles
(
    id         int(10) unsigned    not null auto_increment comment 'ID',
    name       varchar(64)         not null default '' comment '名称，中文名',
    status     tinyint(1) unsigned not null default '2' comment '状态：1-禁用，2-启用',
    weight     tinyint(1) unsigned not null default '0' comment '权重，值越大越靠前',
    remark     varchar(255)        not null default '' comment '备注',
    created_at datetime            not null default current_timestamp comment '创建时间',
    updated_at datetime            not null default current_timestamp comment '更新时间',
    primary key (id),
    unique key (name),
    key (updated_at)
) engine = InnoDB
  auto_increment = 101
  default charset = utf8mb4 comment '角色表';
