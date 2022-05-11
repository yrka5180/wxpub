-- ----------------------------
-- 初始化数据库pub_platform_mgr
-- ----------------------------

CREATE DATABASE IF NOT EXISTS pub_platform_mgr
    DEFAULT CHARSET utf8mb4
    COLLATE utf8mb4_general_ci;

USE pub_platform_mgr;
SET NAMES utf8mb4;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
                         `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
                        `open_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'open id',
                        `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
                        `delete_time` int(11) NOT NULL DEFAULT 0 COMMENT '删除时间',
                        PRIMARY KEY (`id`) USING BTREE,
                        UNIQUE INDEX `app_name`(`open_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户user表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for failure_msg_log
-- ----------------------------
DROP TABLE IF EXISTS `failure_msg_log`;
CREATE TABLE `failure_msg_log`  (
                         `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
                         `msg_id` bigint(64) NOT NULL DEFAULT 0 COMMENT '消息id',
                         `to_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接收者openid',
                         `template_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '模板id',
                         `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '模板内容',
                         `cause` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '失败原因',
                         `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '发送状态，1为正常，2为重试中，3为失败',
                         `count` tinyint(1) NOT NULL DEFAULT 0 COMMENT '发送次数',
                         `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
                         `update_time` int(11) NOT NULL DEFAULT 0 COMMENT '回调更新时间',
                         PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '消息发送失败日志表' ROW_FORMAT = DYNAMIC;