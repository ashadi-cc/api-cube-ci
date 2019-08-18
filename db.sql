-- phpMyAdmin SQL Dump
-- version 4.9.0.1
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Aug 18, 2019 at 05:40 AM
-- Server version: 10.4.6-MariaDB-1:10.4.6+maria~bionic
-- PHP Version: 7.2.19

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `cubes`
--

-- --------------------------------------------------------

--
-- Table structure for table `cubes`
--

CREATE TABLE `cubes` (
  `id` int(11) NOT NULL,
  `rate_date` date NOT NULL,
  `currency` varchar(16) NOT NULL,
  `rate` float NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------
--
-- Indexes for dumped tables
--

--
-- Indexes for table `cubes`
--
ALTER TABLE `cubes`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `cubes`
  ADD UNIQUE rate_currency_id(
    `rate_date`, `currency`
  );

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `cubes`
--
ALTER TABLE `cubes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;