-- # Dump of table article_category
DROP TABLE IF EXISTS `article_category`;

CREATE TABLE `article_category` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `article_id` int(11) unsigned NOT NULL,
  `category_id` int(11) unsigned NOT NULL,
  PRIMARY key (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
