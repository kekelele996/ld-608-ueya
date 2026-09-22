-- ground-turn 航空地勤周转保障平台 · MySQL 8 初始化脚本
-- 说明：docker-compose 中 mysql:8.0 首次启动时自动执行本文件。
-- 应用启动后 GORM AutoMigrate 会再做一次幂等对账（补缺列/索引），
-- 因此此处保留完整 DDL 便于 DBA 审阅与离线部署。

CREATE DATABASE IF NOT EXISTS app_db
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE app_db;

CREATE TABLE IF NOT EXISTS app_user (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username      VARCHAR(64)  NOT NULL,
  password      VARCHAR(128) NOT NULL COMMENT 'bcrypt 哈希',
  name          VARCHAR(64)  NOT NULL,
  role          VARCHAR(32)  NOT NULL COMMENT 'DISPATCHER/CREW/RESOURCE/SUPERVISOR',
  team_id       VARCHAR(32)  NULL,
  created_at    DATETIME(3)  NULL,
  updated_at    DATETIME(3)  NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_app_user_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS flight_turnaround (
  id                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  flight_no          VARCHAR(32)  NOT NULL,
  aircraft_reg       VARCHAR(32)  NOT NULL,
  stand_no           VARCHAR(32)  NOT NULL,
  arrival_time       DATETIME(3)  NOT NULL,
  departure_time     DATETIME(3)  NOT NULL,
  turnaround_status  VARCHAR(32)  NOT NULL DEFAULT 'ARRIVING',
  delay_reason       VARCHAR(255) NULL,
  accumulated_delay  INT          NOT NULL DEFAULT 0,
  plan_generated     TINYINT(1)   NOT NULL DEFAULT 0,
  created_at         DATETIME(3)  NULL,
  updated_at         DATETIME(3)  NULL,
  PRIMARY KEY (id),
  KEY idx_turnaround_status (turnaround_status),
  KEY idx_turnaround_arrival (arrival_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ground_task (
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  turnaround_id  BIGINT UNSIGNED NOT NULL,
  task_type      VARCHAR(32)  NOT NULL,
  team_id        VARCHAR(32)  NOT NULL,
  planned_start  DATETIME(3)  NOT NULL,
  planned_end    DATETIME(3)  NOT NULL,
  deadline       DATETIME(3)  NOT NULL,
  actual_finish  DATETIME(3)  NULL,
  status         VARCHAR(32)  NOT NULL DEFAULT 'PLANNED',
  blocker_note   VARCHAR(255) NULL,
  signed_by      VARCHAR(64)  NULL,
  created_at     DATETIME(3)  NULL,
  updated_at     DATETIME(3)  NULL,
  PRIMARY KEY (id),
  KEY idx_task_turnaround (turnaround_id),
  KEY idx_task_type (task_type),
  KEY idx_task_team (team_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ground_resource (
  id                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_code        VARCHAR(32)  NOT NULL,
  resource_type        VARCHAR(32)  NOT NULL,
  location             VARCHAR(64)  NOT NULL,
  availability_status  VARCHAR(32)  NOT NULL DEFAULT 'AVAILABLE',
  maintenance_due_at   DATETIME(3)  NULL,
  owner_team           VARCHAR(32)  NOT NULL,
  available_from       DATETIME(3)  NULL,
  available_to         DATETIME(3)  NULL,
  created_at           DATETIME(3)  NULL,
  updated_at           DATETIME(3)  NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_resource_code (resource_code),
  KEY idx_resource_type (resource_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS resource_booking (
  id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_id     BIGINT UNSIGNED NOT NULL,
  turnaround_id   BIGINT UNSIGNED NOT NULL,
  task_id         BIGINT UNSIGNED NULL,
  start_time      DATETIME(3)  NOT NULL,
  end_time        DATETIME(3)  NOT NULL,
  booking_status  VARCHAR(32)  NOT NULL DEFAULT 'CONFIRMED',
  conflict_reason VARCHAR(255) NULL,
  created_at      DATETIME(3)  NULL,
  updated_at      DATETIME(3)  NULL,
  PRIMARY KEY (id),
  KEY idx_booking_resource (resource_id),
  KEY idx_booking_turnaround (turnaround_id),
  KEY idx_booking_task (task_id),
  KEY idx_booking_start (start_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS delay_event (
  id                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  turnaround_id       BIGINT UNSIGNED NOT NULL,
  delay_type          VARCHAR(32)  NOT NULL,
  minutes             INT          NOT NULL,
  root_cause          VARCHAR(255) NULL,
  responsibility_team VARCHAR(32)  NULL,
  resolved_at         DATETIME(3)  NULL,
  created_at          DATETIME(3)  NULL,
  updated_at          DATETIME(3)  NULL,
  PRIMARY KEY (id),
  KEY idx_delay_turnaround (turnaround_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS audit_log (
  id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  actor       VARCHAR(64)  NOT NULL,
  role        VARCHAR(32)  NOT NULL,
  action      VARCHAR(64)  NOT NULL,
  target_type VARCHAR(32)  NOT NULL,
  target_id   VARCHAR(64)  NULL,
  detail      VARCHAR(512) NULL,
  created_at  DATETIME(3)  NULL,
  PRIMARY KEY (id),
  KEY idx_audit_actor (actor),
  KEY idx_audit_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
