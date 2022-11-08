# 菜单表
create table menus
(
    id         int(10) unsigned    not null auto_increment comment 'ID',
    pid        int(10) unsigned    not null default '0' comment '父ID',
    name       varchar(32)         not null default '' comment '名称',
    icon       varchar(128)        not null default '' comment '图标',
    path       varchar(255)        not null default '' comment '地址',
    component  varchar(255)        not null default '' comment '组件',
    type       tinyint(1) unsigned not null default '1' comment '类型：1-目录，2-菜单，3-按钮',
    mode       tinyint(1) unsigned not null default '1' comment '模式：1-组件，2-内链，3-外链',
    weight     int(10) unsigned    not null default '0' comment '权重，值越大越靠前',
    level      tinyint(1) unsigned not null default '1' comment '等级，表示几级菜单',
    is_show    tinyint(1) unsigned not null default '2' comment '是否显示：1-否，2-是',
    status     tinyint(1) unsigned not null default '2' comment '状态：1-禁用，2-启用',
    created_at datetime            not null default current_timestamp comment '创建时间',
    updated_at datetime            not null default current_timestamp comment '更新时间',
    primary key (id),
    key (name),
    key (updated_at)
) engine = InnoDB
  auto_increment = 101
  default charset = utf8mb4 comment '菜单表';
