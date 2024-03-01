# 登录日志表
create table login_logs
(
    id         int(10) unsigned    not null auto_increment comment 'ID',
    username   varchar(32)         not null default '' comment '用户名',
    nickname   varchar(32)         not null default '' comment '昵称',
    ip         varchar(16)         not null default '' comment 'IP地址',
    device     varchar(64)         not null default '' comment '设备型号',
    os         varchar(32)         not null default '' comment '操作系统',
    browser    varchar(32)         not null default '' comment '浏览器',
    type       tinyint(1) unsigned not null default '2' comment '日志类型：1-登录失败，2-登录成功，3-退出登录',
    remark     varchar(255)        not null default '' comment '备注',
    created_at datetime            not null default current_timestamp comment '创建时间',
    updated_at datetime            not null default current_timestamp comment '更新时间',
    primary key (id),
    key (created_at)
) engine = InnoDB
  auto_increment = 1
  default charset = utf8mb4 comment '登录日志表';
