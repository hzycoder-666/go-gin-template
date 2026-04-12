CREATE TABLE IF NOT EXISTS ai_generation_tasks(
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  task_id VARCHAR(64) NOT NULL UNIQUE,
  user_id BIGINT NOT NULL,
  type VARCHAR(20) DEFAULT 'imagine',
  platform VARCHAR(20) DEFAULT 'midjourney',
  prompt TEXT,
  status VARCHAR(20) DEFAULT 'pending',
  raw_response JSON,
  image_url TEXT,
  error_message TEXT,
  custom_id VARCHAR(200),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  finished_at DATETIME NULL,
  INDEX idx_user_id (user_id),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at)
) engine=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI生成式任务信息表';
