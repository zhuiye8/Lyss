-- 删除索引
DROP INDEX IF EXISTS idx_models_provider;
DROP INDEX IF EXISTS idx_models_type;
DROP INDEX IF EXISTS idx_models_status;
DROP INDEX IF EXISTS idx_models_org_id;

DROP INDEX IF EXISTS idx_model_configs_model_id;
DROP INDEX IF EXISTS idx_model_configs_org_id;
DROP INDEX IF EXISTS idx_model_configs_created_by;
DROP INDEX IF EXISTS idx_model_configs_shared;

-- 删除表
DROP TABLE IF EXISTS model_configs;
DROP TABLE IF EXISTS models; 