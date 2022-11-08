# 授权关系表
create table auth_relations
(
    typ        tinyint(1) unsigned not null default '1' comment '类型：1-账号角色，2-角色菜单',
    aid        int(10) unsigned    not null default '0' comment '账号表ID',
    rid        int(10) unsigned    not null default '0' comment '角色表ID',
    mid        int(10) unsigned    not null default '0' comment '菜单表ID',
    created_at datetime            not null default current_timestamp comment '创建时间',
    key (aid),
    key (rid),
    key (mid)
) engine = InnoDB
  auto_increment = 101
  default charset = utf8mb4 comment '授权关系表';
