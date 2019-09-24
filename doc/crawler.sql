DROP TABLE IF EXISTS `deliveroo`;
CREATE TABLE `deliveroo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '商家名称',
  `score` varchar(45) NOT NULL COMMENT '评分',
  `evaluation` varchar(45) NOT NULL COMMENT '评价',
  `food_type` varchar(255) NOT NULL COMMENT '菜品类型 `,`分隔',
  `area` varchar(45) NOT NULL COMMENT '行政区',
  `place` varchar(45) NOT NULL COMMENT '地名',
  `banner` varchar(45) NOT NULL COMMENT '广告位',
  `created_at` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=13046 DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `foodpanda`;
CREATE TABLE `foodpanda` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '商家名称',
  `score` float NOT NULL COMMENT '评分',
  `evaluation` int(11) NOT NULL COMMENT '评价',
  `food_type` varchar(255) NOT NULL COMMENT '菜品类型 `,`分隔',
  `address` varchar(255) NOT NULL COMMENT '地址',
  `latitude` varchar(45) NOT NULL COMMENT '经度',
  `longitude` varchar(45) NOT NULL COMMENT '纬度',
  `banner` varchar(45) NOT NULL COMMENT '广告位',
  `created_at` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3175 DEFAULT CHARSET=utf8mb4;