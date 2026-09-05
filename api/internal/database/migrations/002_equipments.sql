CREATE TABLE equipments (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(150) NOT NULL,
  category VARCHAR(50) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'available',
  PRIMARY KEY (id),
  UNIQUE KEY uq_equipments_code (code),
  CONSTRAINT chk_equipments_status CHECK (status IN ('available', 'borrowed', 'repair'))
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci