# 文件表
create table files
(
    id         int(10) unsigned not null auto_increment comment 'ID',
    file_name  varchar(128)     not null default '' comment '文件名',
    file_size  int(10) unsigned not null default '0' comment '文件大小，单位：B',
    file_type  varchar(32)      not null default '' comment '文件类型',
    file_url   varchar(255)     not null default '' comment '文件地址',
    thumbnail  varchar(512)     not null default '' comment '缩略图地址',
    provider   varchar(32)      not null default '' comment '提供商：qiniu-七牛云，ali-阿里云，tencent-腾讯云',
    username   varchar(32)      not null default '' comment '用户名',
    nickname   varchar(32)      not null default '' comment '昵称',
    remark     varchar(500)     not null default '' comment '备注',
    created_at datetime         not null default current_timestamp comment '创建时间',
    updated_at datetime         not null default current_timestamp comment '更新时间',
    primary key (id),
    key (created_at)
) engine = InnoDB
  auto_increment = 1
  default charset = utf8mb4 comment '文件表';
