-- ground-turn schema (mirrors GORM AutoMigrate in backend/src/config/database.go)
CREATE TABLE IF NOT EXISTS app_user (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL,
  display_name VARCHAR(64) NOT NULL,
  role VARCHAR(32) NOT NULL,
  team_id VARCHAR(32) NULL,
  created_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_app_user_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS flight_turnaround (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  flight_no VARCHAR(32) NOT NULL,
  aircraft_reg VARCHAR(32) NOT NULL,
  stand_no VARCHAR(16) NOT NULL,
  arrival_time DATETIME(3) NOT NULL,
  departure_time DATETIME(3) NOT NULL,
  turnaround_status VARCHAR(16) NOT NULL,
  delay_reason VARCHAR(255) NULL,
  total_delay_minutes INT NOT NULL DEFAULT 0,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_turn_status (turnaround_status),
  KEY idx_turn_arrival (arrival_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ground_task (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  turnaround_id BIGINT UNSIGNED NOT NULL,
  task_type VARCHAR(32) NOT NULL,
  team_id VARCHAR(32) NOT NULL,
  planned_start DATETIME(3) NOT NULL,
  deadline DATETIME(3) NOT NULL,
  actual_finish DATETIME(3) NULL,
  status VARCHAR(16) NOT NULL,
  blocker_note VARCHAR(255) NULL,
  signed_at DATETIME(3) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_task_turn (turnaround_id),
  KEY idx_task_type (task_type),
  KEY idx_task_status (status),
  KEY idx_task_start (planned_start)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ground_resource (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_code VARCHAR(32) NOT NULL,
  resource_type VARCHAR(32) NOT NULL,
  location VARCHAR(64) NULL,
  availability_status VARCHAR(16) NOT NULL,
  maintenance_due_at DATETIME(3) NULL,
  maintenance_end DATETIME(3) NULL,
  owner_team VARCHAR(32) NOT NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_resource_code (resource_code),
  KEY idx_resource_type (resource_type),
  KEY idx_resource_status (availability_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS resource_booking (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_id BIGINT UNSIGNED NOT NULL,
  turnaround_id BIGINT UNSIGNED NOT NULL,
  task_id BIGINT UNSIGNED NOT NULL,
  start_time DATETIME(3) NOT NULL,
  end_time DATETIME(3) NOT NULL,
  booking_status VARCHAR(16) NOT NULL,
  conflict_reason VARCHAR(64) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_booking_resource (resource_id),
  KEY idx_booking_turn (turnaround_id),
  KEY idx_booking_task (task_id),
  KEY idx_booking_status (booking_status),
  KEY idx_booking_start (start_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS delay_event (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  turnaround_id BIGINT UNSIGNED NOT NULL,
  delay_type VARCHAR(32) NOT NULL,
  minutes INT NOT NULL,
  root_cause VARCHAR(255) NULL,
  responsibility_team VARCHAR(32) NULL,
  resolved_at DATETIME(3) NULL,
  created_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_delay_turn (turnaround_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS audit_log (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  actor VARCHAR(64) NOT NULL,
  role VARCHAR(32) NOT NULL,
  action VARCHAR(255) NOT NULL,
  target_type VARCHAR(32) NOT NULL,
  target_id VARCHAR(64) NULL,
  created_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_audit_actor (actor),
  KEY idx_audit_target (target_type, target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
