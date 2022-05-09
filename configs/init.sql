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
                        PRIMARY KEY (`id`) USING BTREE,
                        UNIQUE INDEX `app_name`(`open_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户user表' ROW_FORMAT = DYNAMIC;