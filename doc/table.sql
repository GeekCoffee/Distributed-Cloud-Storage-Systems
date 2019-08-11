-- 创建一个database
create database fileserver default character set utf8mb4;

-- tbl_file文件元信息表
create table `tbl_filemeta` (
`id` int(11) not null auto_increment comment 'tbl_file表的id',
`file_sha1` char(40) not null default '' comment '文件内容通过SHA1算法，得到的hash值',
`file_name` varchar(256) not null default '' comment '文件名',
`file_path` varchar(1024) not null default '' comment '文件路径+文件名',
`create_at` datetime default NOW() comment '文件被创建的时间',
`update_at` datetime default NOW() on update current_timestamp() comment '文件每次被修改的时间',
`status` int(11) not null default 0 comment '文件的状态(可用-0, 禁用-1, 被删除-2)',
`file_size` bigint(20) default 0 comment '文件大小',
`ext1` int(11) default 0 comment '备用字段1',
`ext2` text comment '备用字段2',
primary key(`id`),
unique key `idx_file_hash` (`file_sha1`),
key `idx_status` (`status`)
) ENGINE=InnoDB default charset=utf8;


-- tbl_user用户信息表，使用default的值而不是NULL
create table `tbl_user`(
`id` int(11) not null auto_increment,
`user_name` varchar(64) not null default '' comment '用户名',
`user_pwd` varchar(256) not null default '' comment '用户密码',
`email` varchar(64) default '' comment '邮箱，用户名也可以是邮箱',
`phone` varchar(128) default '' comment '手机号',
`email_validated` tinyint(1) default 0 comment '邮箱是否被验证过',
`phone_validated` tinyint(1) default 0 comment '手机号是否被验证过',
`signup_at` datetime default current_timestamp comment '注册日期, signup-注册、签约; login-登录',
`last_active` datetime default current_timestamp on update current_timestamp comment '最后活跃时间',
`profile` text comment '用户属性',
`status` int(11) not null default 0 comment '用户状态(可用|禁用|锁定|删除标记等)',
primary key(`id`),
unique key `idx_phone` (`phone`),
unique key `idx_email` (`email`),
key `idx_status` (`status`)
)engine=InnoDB auto_increment=5 default charset=utf8mb4;


-- tbl_user_token用户访问令牌表
create table `tbl_user_token` (
`id` int(11) not null auto_increment,
`user_name` varchar(64) not null default '' comment '用户名，唯一',
`user_token` varchar(64) not null default '' comment '用户登录token凭证',
primary key (`id`),
unique key `idx_username` (`user_name`)
)engine = InnoDB default charset=utf8mb4;

-- tbl_user_file 用户下的文件表
create table `tbl_user_file` (
`id` int(11) not null auto_increment comment 'id是主键',
`user_name` varchar(64) not null,
`file_sha1` varchar(64) not null default '' comment '文件的hash值',
`file_size` bigint(20) default 0 comment '文件大小',
`file_name` varchar(256) not null default '' comment '文件名称',
`upload_at` datetime default current_timestamp comment '上传文件的时间',
`last_update` datetime default current_timestamp on update current_timestamp comment '最后文件修订时间',
`status` int(11) not null default 0 comment '文件状态(可用0，禁用1，被删除-1)',
primary key (`id`),
unique key `idx_user_file` (`user_name`,`file_sha1`),
key `idx_status` (`status`),
key `idx_user_id` (`user_name`)
)engine=InnoDB default charset=utf8mb4;
