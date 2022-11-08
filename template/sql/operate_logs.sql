# 操作日志表
create table operate_logs
(
    id         int(10) unsigned not null auto_increment comment 'ID',
    username   varchar(32)      not null default '' comment '用户名',
    nickname   varchar(32)      not null default '' comment '昵称',
    ip         varchar(16)      not null default '' comment 'IP地址',
    function   varchar(32)      not null default '' comment '操作功能',
    uri        varchar(255)     not null default '' comment '请求地址',
    method     varchar(32)      not null default '' comment '请求方式',
    params     varchar(1000)    not null default '' comment '请求参数',
    status     int(10) unsigned not null default '0' comment '状态码',
    code       varchar(10)      not null default '' comment '业务码',
    spend_time int(10) unsigned not null default '0' comment '耗时，单位：ms',
    result     varchar(5000)    not null default '' comment '响应结果',
    user_agent varchar(500)     not null default '' comment '浏览器信息',
    remark     varchar(500)     not null default '' comment '备注',
    created_at datetime         not null default current_timestamp comment '创建时间',
    updated_at datetime         not null default current_timestamp comment '更新时间',
    primary key (id),
    key (spend_time),
    key (created_at)
) engine = InnoDB
  auto_increment = 1
  default charset = utf8mb4 comment '操作日志表';
