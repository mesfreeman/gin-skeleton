# 公共配置表
create table common_configs
(
    id         int(10) unsigned not null auto_increment comment 'ID',
    module     varchar(64)      not null default '' comment '模块',
    keyword    varchar(64)      not null default '' comment '关键词',
    value      varchar(2000)    not null default '' comment '配置值',
    remark     varchar(255)     not null default '' comment '备注',
    created_at datetime         not null default current_timestamp comment '创建时间',
    updated_at datetime         not null default current_timestamp comment '更新时间',
    primary key (id),
    unique key (module, keyword)
) engine = InnoDB
  auto_increment = 1
  default charset = utf8mb4 comment '公共配置表';
