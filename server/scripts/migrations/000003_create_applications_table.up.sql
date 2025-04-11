CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    config JSONB DEFAULT '{}'::jsonb,
    model_config JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_applications_project_id ON applications(project_id);
CREATE INDEX idx_applications_type ON applications(type);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_applications_deleted_at ON applications(deleted_at); 