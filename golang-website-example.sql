-- MariaDB dump 10.19  Distrib 10.5.18-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: golang-website-example
-- ------------------------------------------------------
-- Server version	10.5.18-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `locations`
--

DROP TABLE IF EXISTS `locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `locations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `location` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=95 DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `locations`
--

LOCK TABLES `locations` WRITE;
/*!40000 ALTER TABLE `locations` DISABLE KEYS */;
INSERT INTO `locations` VALUES (11,'ACEH'),(12,'SUMATERA UTARA'),(13,'SUMATERA BARAT'),(14,'RIAU'),(15,'JAMBI'),(16,'SUMATERA SELATAN'),(17,'BENGKULU'),(18,'LAMPUNG'),(19,'KEPULAUAN BANGKA BELITUNG'),(21,'KEPULAUAN RIAU'),(31,'DKI JAKARTA'),(32,'JAWA BARAT'),(33,'JAWA TENGAH'),(34,'DI YOGYAKARTA'),(35,'JAWA TIMUR'),(36,'BANTEN'),(51,'BALI'),(52,'NUSA TENGGARA BARAT'),(53,'NUSA TENGGARA TIMUR'),(61,'KALIMANTAN BARAT'),(62,'KALIMANTAN TENGAH'),(63,'KALIMANTAN SELATAN'),(64,'KALIMANTAN TIMUR'),(65,'KALIMANTAN UTARA'),(71,'SULAWESI UTARA'),(72,'SULAWESI TENGAH'),(73,'SULAWESI SELATAN'),(74,'SULAWESI TENGGARA'),(75,'GORONTALO'),(76,'SULAWESI BARAT'),(81,'MALUKU'),(82,'MALUKU UTARA'),(91,'PAPUA BARAT'),(94,'PAPUA');
/*!40000 ALTER TABLE `locations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `role` enum('admin','user') NOT NULL,
  `username` varchar(15) NOT NULL,
  `email` varchar(30) NOT NULL,
  `password` varchar(255) NOT NULL,
  `name` varchar(30) NOT NULL,
  `location` int(10) unsigned DEFAULT NULL,
  `photo` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email_UNIQUE` (`email`),
  UNIQUE KEY `username_UNIQUE` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2020-02-22 15:28:40','2022-02-05 10:14:04',NULL,'admin','admin','admin@website.com','$2a$10$XJAj65HZ2c.n1iium4qUEeGarW0PJsqVcedBh.PDGMXdjqfOdN1hW','Admin',32,''),(2,'2020-03-17 06:56:10','2021-11-25 09:30:58',NULL,'user','sugriwa','sugriwa@wanara.com','$2a$10$bVVMuFHe/iaydX9yO2AttOPT8WyhMPe9F8nDflEqEyJbGRD5.guFu','Sugriwa',0,''),(3,'2020-03-17 06:57:47','2021-11-25 09:30:58',NULL,'user','subali','subali@wanara.com','$2a$10$eO8wPLSfBU.8KLUh/T9kDeBm0vIRjiCvsmWe8ou5fZHJ3cYAUcg6y','Subali',0,'');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-02-28 14:35:09
