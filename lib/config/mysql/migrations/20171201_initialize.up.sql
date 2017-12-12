-- Initialize the database

BEGIN;

-- -----------------------------------------------------
-- Table `Instances`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Instances` (
	  `Guid` VARCHAR(36) NOT NULL,
	  `Name` VARCHAR(45) NULL,
	  `TargetURL` TEXT NULL,
	  `AuthKey` VARCHAR(45) NULL,
	  `CaCert` BLOB NULL,
	  `SkipSSL` TINYINT(1) NULL,
	  PRIMARY KEY (`Guid`))
	ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `Plans`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Plans` (
	  `Guid` VARCHAR(36) NOT NULL,
	  `Name` TEXT(255) NULL,
	  `Description` TEXT(255) NULL,
	  `Free` TINYINT(1) NULL COMMENT '    ',
	  `Metadata` BLOB NULL,
	  PRIMARY KEY (`Guid`))
	ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `Dials`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Dials` (
	  `Guid` VARCHAR(36) NOT NULL,
	  `Configuration` BLOB NULL,
	  `Plans_Guid` VARCHAR(36) NOT NULL,
	  `Instances_Guid` VARCHAR(36) NOT NULL,
	  PRIMARY KEY (`Guid`),
	  INDEX `fk_Dials_Plans1_idx` (`Plans_Guid` ASC),
	  INDEX `fk_Dials_Instances1_idx` (`Instances_Guid` ASC),
	  CONSTRAINT `fk_Dials_Plans1`
	    FOREIGN KEY (`Plans_Guid`)
	    REFERENCES `Plans` (`Guid`)
	    ON DELETE NO ACTION
	    ON UPDATE NO ACTION,
	  CONSTRAINT `fk_Dials_Instances1`
	    FOREIGN KEY (`Instances_Guid`)
	    REFERENCES `Instances` (`Guid`)
	    ON DELETE NO ACTION
	    ON UPDATE NO ACTION)
	ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `Services`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Services` (
	  `Guid` VARCHAR(36) NOT NULL,
	  `Bindable` TINYINT(1) NULL,
	  `DashboardClient` BLOB NULL,
	  `Description` TEXT(255) NULL,
	  `Metadata` BLOB NULL,
	  `Name` TEXT(255) NULL,
	  `PlanUpdateable` TINYINT(1) NULL,
	  `Tags` BLOB NULL,
	  `Instances_Guid` VARCHAR(36) NOT NULL,
	  `Requires` BLOB NULL,
	  PRIMARY KEY (`Guid`),
	  INDEX `fk_Services_Instances1_idx` (`Instances_Guid` ASC),
	  CONSTRAINT `fk_Services_Instances1`
	    FOREIGN KEY (`Instances_Guid`)
	    REFERENCES `Instances` (`Guid`)
	    ON DELETE NO ACTION
	    ON UPDATE NO ACTION)
	ENGINE = InnoDB;

COMMIT;
