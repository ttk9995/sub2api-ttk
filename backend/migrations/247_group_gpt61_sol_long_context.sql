ALTER TABLE groups
    ADD COLUMN IF NOT EXISTS gpt61_sol_long_context_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS gpt61_sol_long_context_threshold INTEGER NOT NULL DEFAULT 500000
        CHECK (gpt61_sol_long_context_threshold > 0);

COMMENT ON COLUMN groups.gpt61_sol_long_context_enabled IS
    'Use the group-specific GPT-6.1 Sol model long-context threshold; disabled defaults to 272000';
COMMENT ON COLUMN groups.gpt61_sol_long_context_threshold IS
    'Input plus cache tokens; the model surcharge applies strictly above this threshold when the override is enabled';
