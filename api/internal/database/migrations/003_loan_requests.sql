CREATE TABLE loan_requests (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  equipment_id BIGINT UNSIGNED NOT NULL,
  purpose VARCHAR(500) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  requested_at DATETIME NOT NULL,
  approved_by BIGINT UNSIGNED NULL,
  approved_at DATETIME NULL,
  returned_at DATETIME NULL,
  PRIMARY KEY (id),
  KEY idx_loans_user_requested (user_id, requested_at DESC),
  KEY idx_loans_status_requested (status, requested_at DESC),
  CONSTRAINT fk_loans_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
  CONSTRAINT fk_loans_equipment FOREIGN KEY (equipment_id) REFERENCES equipments (id) ON DELETE RESTRICT,
  CONSTRAINT fk_loans_approver FOREIGN KEY (approved_by) REFERENCES users (id) ON DELETE RESTRICT,
  CONSTRAINT chk_loans_status CHECK (
    status IN ('pending', 'approved', 'rejected', 'returned')
  )
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci