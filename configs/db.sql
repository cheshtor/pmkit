/* 用户 */
DROP TABLE IF EXISTS `pk_user`;
CREATE TABLE `pk_user`
(
    `id`          BIGINT(20) COMMENT '用户 ID',
    `phone`       VARCHAR(12) UNIQUE COMMENT '登录手机',
    `password`    VARCHAR(128) COMMENT '登录密码',
    `username`    VARCHAR(10) UNIQUE COMMENT '用户名',
    `create_time` BIGINT(20) DEFAULT 0 COMMENT '创建时间',
    `active`      TINYINT(1) DEFAULT 1 COMMENT '账号是否可用。1 - 可用; 0 - 禁用',
    PRIMARY KEY `pk_id` (`id`)
);

/* 角色 */
DROP TABLE IF EXISTS `pk_role`;
CREATE TABLE `pk_role`
(
    `id`     BIGINT(20) COMMENT '角色 ID',
    `name`   VARCHAR(32) COMMENT '角色名称',
    `remark` VARCHAR(128) COMMENT '角色备注',
    PRIMARY KEY `pk_id` (`id`)
);

/* 用户角色关联 */
DROP TABLE IF EXISTS `pk_user_role`;
CREATE TABLE `pk_user_role`
(
    `id`          BIGINT(20) COMMENT '关联 ID',
    `user_id`     BIGINT(20) COMMENT '用户 ID',
    `role_id`     BIGINT(20) COMMENT '角色 ID',
    `create_by`   BIGINT(20) DEFAULT 0 COMMENT '创建人 ID',
    `create_time` BIGINT(20) DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY `pk_id` (`id`)
);

/* 项目 */
DROP TABLE IF EXISTS `pk_project`;
CREATE TABLE `pk_project`
(
    `id`            BIGINT(20) COMMENT '项目 ID',
    `name`          VARCHAR(64) COMMENT '项目名称',
    `employer`      VARCHAR(64) COMMENT '建设单位',
    `contractor`    VARCHAR(64) COMMENT '承建单位',
    `supervisor`    VARCHAR(64) COMMENT '监理单位',
    `executor`      VARCHAR(64) COMMENT '实施单位',
    `description`   VARCHAR(512) COMMENT '项目简介',
    `start_date`    BIGINT(20) DEFAULT 0 COMMENT '项目启动日期',
    `end_date`      BIGINT(20) DEFAULT 0 COMMENT '项目完成日期',
    `status`        TINYINT(2) COMMENT '项目状态。1 - 待建; 2 - 在建; 3 - 暂停; 4 - 终止; 5 - 结项',
    `delete`        TINYINT(1) DEFAULT 0 COMMENT '删除。0 - 可用; 1 - 删除',
    `create_by`     BIGINT(20) DEFAULT 0 COMMENT '创建人 ID',
    `create_time`   BIGINT(20) DEFAULT 0 COMMENT '创建时间',
    `modified_by`   BIGINT(20) DEFAULT 0 COMMENT '修改人 ID',
    `modified_time` BIGINT(20) DEFAULT 0 COMMENT '修改时间',
    PRIMARY KEY `pk_id` (`id`)
);

/* 迭代 */
DROP TABLE IF EXISTS `pk_iteration`;
CREATE TABLE `pk_iteration`
(
    `id`            BIGINT(20) COMMENT '迭代 ID',
    `project_id`    BIGINT(20) COMMENT '项目 ID',
    `name`          VARCHAR(64) COMMENT '迭代名称',
    `remark`        VARCHAR(256) COMMENT '迭代备注',
    `start_date`    BIGINT(20) COMMENT '迭代开始日期',
    `end_date`      BIGINT(20) COMMENT '迭代完成日期',
    `create_by`     BIGINT(20) DEFAULT 0 COMMENT '创建人 ID',
    `create_time`   BIGINT(20) DEFAULT 0 COMMENT '创建时间',
    `modified_by`   BIGINT(20) DEFAULT 0 COMMENT '修改人 ID',
    `modified_time` BIGINT(20) DEFAULT 0 COMMENT '修改时间',
    `delete`        TINYINT(1) DEFAULT 0 COMMENT '删除。0 - 可用; 1 - 删除',
    PRIMARY KEY `pk_id` (`id`)
);

/* 需求 */
DROP TABLE IF EXISTS `pk_requirement`;
CREATE TABLE `pk_requirement`
(
    `id`            BIGINT(20) COMMENT '需求 ID',
    `project_id`    BIGINT(20) COMMENT '所属项目 ID',
    `iteration_id`  BIGINT(20) COMMENT '所属迭代 ID',
    `code`          VARCHAR(64) COMMENT '需求编号',
    `name`          VARCHAR(128) COMMENT '需求名称',
    `type`          TINYINT(2) COMMENT '需求类型。1 - 业务需求; 2 - BUG',
    `demander`      VARCHAR(64) COMMENT '需求方',
    `priority`      TINYINT(2) COMMENT '优先级',
    `influence`     TINYINT(2) COMMENT '影响程度',
    `status`        TINYINT(2) COMMENT '需求状态。1 - 待审核; 2 - 审核中; 3 - 待实施; 4 - 实施中; 5 - 待交付; 6 - 已交付',
    `delete`        TINYINT(1) DEFAULT 0 COMMENT '删除。0 - 可用; 1 - 删除',
    `create_by`     BIGINT(20) DEFAULT 0 COMMENT '创建人 ID',
    `create_time`   BIGINT(20) DEFAULT 0 COMMENT '创建时间',
    `modified_by`   BIGINT(20) DEFAULT 0 COMMENT '修改人 ID',
    `modified_time` BIGINT(20) DEFAULT 0 COMMENT '修改时间',
    PRIMARY KEY `pk_id` (`id`)
);

/* 需求内容 */
DROP TABLE IF EXISTS `pk_requirement_content`;
CREATE TABLE `pk_requirement_content`
(
    `requirement_id` BIGINT(20) COMMENT '需求 ID',
    `content`        LONGTEXT COMMENT '需求内容',
    PRIMARY KEY `pk_requirement_id` (`requirement_id`)
);

/* 需求状态追踪 */
DROP TABLE IF EXISTS `pk_requirement_track`;
CREATE TABLE `pk_requirement_track`
(
    `id`             BIGINT(20) COMMENT '主键 ID',
    `requirement_id` BIGINT(20) COMMENT '需求 ID',
    `status`         TINYINT(2) COMMENT '需求状态',
    `parent_id`      BIGINT(20) COMMENT '父追踪记录 ID',
    `create_by`      BIGINT(20) DEFAULT 0 COMMENT '创建人 ID',
    `create_time`    BIGINT(20) DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY `pk_id` (`id`)
);

/* 需求评论 */
DROP TABLE IF EXISTS `pk_requirement_comment`;
CREATE TABLE `pk_requirement_comment`
(
    `id`             BIGINT(20) COMMENT '需求评论 ID',
    `requirement_id` BIGINT(20) COMMENT '需求 ID',
    `comment`        VARCHAR(512) COMMENT '评论内容',
    `parent_id`      BIGINT(20) COMMENT '父评论 ID',
    `delete`         TINYINT(1) DEFAULT 0 COMMENT '删除。0 - 可用; 1 - 删除',
    `create_by`      BIGINT(20) DEFAULT 0 COMMENT '创建人 ID',
    `create_time`    BIGINT(20) DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY `pk_id` (`id`)
);