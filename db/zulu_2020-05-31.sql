# ************************************************************
# Sequel Pro SQL dump
# Version 5224
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.19)
# Database: zulu
# Generation Time: 2020-05-31 01:55:17 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table institutions
# ------------------------------------------------------------

DROP TABLE IF EXISTS `institutions`;

CREATE TABLE `institutions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL,
  `name` varchar(200) NOT NULL,
  `street_address` varchar(255) NOT NULL,
  `street_map_id` int DEFAULT NULL,
  `bill_address` varchar(100) NOT NULL,
  `bill_map_id` int DEFAULT NULL,
  `pic_name` varchar(255) NOT NULL,
  `pic_phone` varchar(255) NOT NULL,
  `expired_at` date DEFAULT NULL,
  `status` int DEFAULT (1),
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` varchar(100) NOT NULL DEFAULT (_utf8mb4'Admin'),
  `updated_at` timestamp NULL DEFAULT NULL,
  `updated_by` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `institutions` WRITE;
/*!40000 ALTER TABLE `institutions` DISABLE KEYS */;

INSERT INTO `institutions` (`id`, `code`, `name`, `street_address`, `street_map_id`, `bill_address`, `bill_map_id`, `pic_name`, `pic_phone`, `expired_at`, `status`, `created_at`, `created_by`, `updated_at`, `updated_by`)
VALUES
	(1,'UI','Universitas Indonesia','Kampus Baru UI Depok',1,'Kampus Baru UI Depok',1,'ASD','0812','2020-07-29',1,'2020-05-30 11:16:58','Admin',NULL,NULL),
	(2,'UIA','Universitas Islam Assyafi`iyah','Jl. Raya Jatiwaringin No.12',1,'Jl. Raya Jatiwaringin No.12',1,'ASD','0812','2020-07-29',1,'2020-05-30 11:16:58','Admin',NULL,NULL),
	(3,'GUNDAR','Universitas Gunadarma','Jl. Margonda Raya No.100',1,'Jl. Margonda Raya No.100',1,'ASD','0812','2020-07-29',1,'2020-05-30 11:16:58','Admin',NULL,NULL);

/*!40000 ALTER TABLE `institutions` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table references
# ------------------------------------------------------------

DROP TABLE IF EXISTS `references`;

CREATE TABLE `references` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_name` varchar(100) NOT NULL,
  `sub_id` int NOT NULL,
  `name` varchar(100) NOT NULL,
  `status` int DEFAULT (1),
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` varchar(100) NOT NULL DEFAULT (_utf8mb4'Admin'),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `references` WRITE;
/*!40000 ALTER TABLE `references` DISABLE KEYS */;

INSERT INTO `references` (`id`, `group_name`, `sub_id`, `name`, `status`, `created_at`, `created_by`)
VALUES
	(1,'student_type_univ',1,'Reguler',1,'2020-05-30 10:49:19','Admin'),
	(2,'student_type_univ',2,'Paruh Waktu',1,'2020-05-30 10:49:19','Admin'),
	(3,'student_type_univ',3,'Malam',1,'2020-05-30 10:49:19','Admin'),
	(4,'student_type_reguler',1,'Pagi',1,'2020-05-30 10:49:19','Admin'),
	(5,'student_type_reguler',2,'Siang',1,'2020-05-30 10:49:19','Admin'),
	(6,'student_type_reguler',3,'Sore',1,'2020-05-30 10:49:19','Admin'),
	(7,'status',0,'Inactive',1,'2020-05-30 10:49:19','Admin'),
	(8,'status',1,'Active',1,'2020-05-30 10:49:19','Admin'),
	(9,'status',2,'On Leave',1,'2020-05-30 10:49:19','Admin'),
	(10,'tutor_type',1,'Tetap',1,'2020-05-30 10:49:19','Admin'),
	(11,'tutor_type',2,'Tidak Tetap',1,'2020-05-30 10:49:19','Admin'),
	(12,'tutor_type',3,'Honorer',1,'2020-05-30 10:49:19','Admin'),
	(38,'sd_degree',1,'Kelas I',1,'2020-05-30 11:00:42','Admin'),
	(39,'sd_degree',2,'Kelas II',1,'2020-05-30 11:00:42','Admin'),
	(40,'sd_degree',3,'Kelas III',1,'2020-05-30 11:00:42','Admin'),
	(41,'sd_degree',4,'Kelas IV',1,'2020-05-30 11:00:42','Admin'),
	(42,'sd_degree',5,'Kelas V',1,'2020-05-30 11:00:42','Admin'),
	(43,'sd_degree',6,'Kelas VI',1,'2020-05-30 11:00:42','Admin'),
	(44,'smp_degree',7,'Kelas VII',1,'2020-05-30 11:00:42','Admin'),
	(45,'smp_degree',8,'Kelas VIII',1,'2020-05-30 11:00:42','Admin'),
	(46,'smp_degree',9,'Kelas IX',1,'2020-05-30 11:00:42','Admin'),
	(47,'sma_degree',10,'Kelas X',1,'2020-05-30 11:00:42','Admin'),
	(48,'sma_degree',11,'Kelas XI',1,'2020-05-30 11:00:42','Admin'),
	(49,'sma_degree',12,'Kelas XII',1,'2020-05-30 11:00:42','Admin'),
	(50,'univ_degree',13,'D1 - Diploma Satu',1,'2020-05-30 11:00:42','Admin'),
	(51,'univ_degree',14,'D2 - Diploma Dua',1,'2020-05-30 11:00:42','Admin'),
	(52,'univ_degree',15,'D3 - Diploma Tiga',1,'2020-05-30 11:00:42','Admin'),
	(53,'univ_degree',16,'D4 - Diploma Empat',1,'2020-05-30 11:00:42','Admin'),
	(54,'univ_degree',17,'S1 - Sarjana',1,'2020-05-30 11:00:42','Admin'),
	(55,'univ_degree',18,'S2 - Magister',1,'2020-05-30 11:00:42','Admin'),
	(56,'univ_degree',19,'S3 - Doktor',1,'2020-05-30 11:00:42','Admin'),
	(57,'univ_degree',20,'Professor',1,'2020-05-30 11:00:42','Admin');

/*!40000 ALTER TABLE `references` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table roles
# ------------------------------------------------------------

DROP TABLE IF EXISTS `roles`;

CREATE TABLE `roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` varchar(100) NOT NULL DEFAULT (_utf8mb4'Admin'),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;

INSERT INTO `roles` (`id`, `name`, `created_at`, `created_by`)
VALUES
	(1,'Super Admin','2020-05-30 11:20:25','Admin'),
	(2,'Admin','2020-05-30 11:20:25','Admin'),
	(3,'Tutor','2020-05-30 11:20:25','Admin'),
	(4,'Student','2020-05-30 11:20:25','Admin'),
	(5,'Parent','2020-05-30 11:20:25','Admin');

/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tutors
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tutors`;

CREATE TABLE `tutors` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `nomor_induk` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `type` varchar(100) NOT NULL,
  `status` int DEFAULT (1),
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` varchar(100) NOT NULL DEFAULT (_utf8mb3'Admin'),
  `updated_at` timestamp NULL DEFAULT NULL,
  `updated_by` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `tutors` WRITE;
/*!40000 ALTER TABLE `tutors` DISABLE KEYS */;

INSERT INTO `tutors` (`id`, `nomor_induk`, `email`, `password`, `type`, `status`, `created_at`, `created_by`, `updated_at`, `updated_by`)
VALUES
	(1,'12312312','asdi@asd.com','pass1','permanent',1,'2020-05-25 21:51:45','Admin',NULL,NULL),
	(2,'12312332','asd@asd.com','pass2','permanent',2,'2020-05-25 21:51:45','Admin',NULL,NULL),
	(3,'2312332','aasd.com','pass3','permanent',3,'2020-05-25 21:51:45','Admin',NULL,NULL);

/*!40000 ALTER TABLE `tutors` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_roles
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_roles`;

CREATE TABLE `user_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  `institution_id` bigint NOT NULL,
  `status` int DEFAULT (1),
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` varchar(100) NOT NULL DEFAULT (_utf8mb4'Admin'),
  `updated_at` timestamp NULL DEFAULT NULL,
  `updated_by` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `user_roles` WRITE;
/*!40000 ALTER TABLE `user_roles` DISABLE KEYS */;

INSERT INTO `user_roles` (`id`, `user_id`, `role_id`, `institution_id`, `status`, `created_at`, `created_by`, `updated_at`, `updated_by`)
VALUES
	(1,5,0,0,1,'2020-05-30 11:51:39','Admin',NULL,NULL),
	(2,6,0,0,1,'2020-05-30 11:51:39','Admin',NULL,NULL),
	(4,9,1,0,1,'2020-05-30 11:52:31','Admin',NULL,NULL),
	(5,10,0,0,1,'2020-05-30 11:52:31','Admin',NULL,NULL),
	(6,1,1,0,1,'2020-05-30 11:58:59','Admin',NULL,NULL),
	(7,2,1,0,1,'2020-05-30 11:58:59','Admin',NULL,NULL),
	(8,3,1,0,1,'2020-05-30 11:58:59','Admin',NULL,NULL),
	(10,13,1,1,1,'2020-05-30 12:00:01','Admin',NULL,NULL),
	(11,14,0,0,1,'2020-05-30 12:00:01','Admin',NULL,NULL),
	(12,1,2,1,1,'2020-05-30 13:05:18','Admin',NULL,NULL),
	(13,1,3,2,1,'2020-05-30 13:05:18','Admin',NULL,NULL),
	(17,21,1,1,1,'2020-05-30 13:52:08','Admin',NULL,NULL),
	(18,22,0,0,1,'2020-05-30 13:52:08','Admin',NULL,NULL),
	(19,23,1,1,1,'2020-05-30 13:53:34','Admin',NULL,NULL),
	(20,24,0,0,1,'2020-05-30 13:53:34','Admin',NULL,NULL);

/*!40000 ALTER TABLE `user_roles` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `is_first` int DEFAULT (1),
  `status` int DEFAULT (1),
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` varchar(100) NOT NULL DEFAULT (_utf8mb3'Admin'),
  `updated_at` timestamp NULL DEFAULT NULL,
  `updated_by` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;

INSERT INTO `users` (`id`, `username`, `password`, `is_first`, `status`, `created_at`, `created_by`, `updated_at`, `updated_by`)
VALUES
	(1,'superadmin@admin.com','$2a$04$V7JViZ9YlkEzgYAAz.hFdOP5hkBCesU2ijiU9Yeh9fbpdP6UsJhQu',1,1,'2020-05-30 02:57:20','Admin',NULL,NULL),
	(2,'admin@admin.com','$2a$04$O6vhq4KvHz8YoFiUrEloXe3eXoAXbF2oGW4lZfOn.LtoXubRpOWL6',1,1,'2020-05-30 02:58:23','Admin',NULL,NULL),
	(3,'admin@superadmin.com','$2a$04$ybAj1lL8rPjG1EAg1i9Be.12c3mRu2ebdIELfX.U6gLK8rA4TBNu.',1,1,'2020-05-30 03:20:31','Admin',NULL,NULL),
	(5,'admin@admin123.com','$2a$04$6y1ZKYLnSQOQ9KH7zRpF8e9Dd5uRcFLw6wRAoZzcy6Zsew8QdCV1.',1,1,'2020-05-30 11:51:39','Admin',NULL,NULL),
	(6,'admin@superadmin123.com','$2a$04$0kB1UXvucfBRyIj76vYeSeVU3yGVhlYf.hoLr6L3V49Ixv/iR5LF.',1,1,'2020-05-30 11:51:39','Admin',NULL,NULL),
	(9,'admin@admin456.com','$2a$04$MwuG1bGKjgZFRsviECp2wuqg.kxYh.C/GFywvQzeJh3MIfZquX9yu',1,1,'2020-05-30 11:52:31','Admin',NULL,NULL),
	(10,'admin@superadmin456.com','$2a$04$nMV.PWITyMk2gVuOuXZJ..L0i9Mc3JvcsRPhS0tbIvYJVqXsCX3iK',1,1,'2020-05-30 11:52:31','Admin',NULL,NULL),
	(13,'admin@admin789.com','$2a$04$TFrX6XQbyXf2Gu8KUKKq/OVoCZFu3AEnXJfNhVKTSKks5lMH0rcOK',1,1,'2020-05-30 12:00:01','Admin',NULL,NULL),
	(14,'admin@superadmin789.com','$2a$04$M3/OFmGYeE5JP0KbQOkLuuDftbxQJ0TorTJOOtmzzVwYj92TRCTsa',1,1,'2020-05-30 12:00:01','Admin',NULL,NULL),
	(21,'admin@admin11.com','$2a$04$DG6cXwzBGaNkRqZVzgbEauPS8TKr4bZthZL.fTG.drSoxZTEQD9AO',1,1,'2020-05-30 13:52:08','Admin',NULL,NULL),
	(22,'admin@superadmin11.com','$2a$04$3H1ziK56qhjtTTi/rleWyunHrVyL1clLW.TSnDn3S6f418XCCc7eW',1,1,'2020-05-30 13:52:08','Admin',NULL,NULL),
	(23,'admin@admin12.com','$2a$04$exKNmB6EsKRLdQal2mTLs.ZbDqNXM0W5IBxXOCyeyN.bq6NPqqJDi',1,1,'2020-05-30 13:53:34','Admin',NULL,NULL),
	(24,'admin@superadmin12.com','$2a$04$wn8ayCTBKPeUBo49Bzm92uWSPjZLIgtxkxa9hEQoiX9LEjmWEHkY6',1,1,'2020-05-30 13:53:34','Admin',NULL,NULL);

/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
