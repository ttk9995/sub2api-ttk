package service

import (
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
)

func supportsGPT61SolLongContextPlatform(platform string) bool {
	return platform == PlatformOpenAI || platform == PlatformComposite
}

func validateGPT61SolLongContextThreshold(threshold int) error {
	if threshold < 1 || threshold > 2_147_483_647 {
		return infraerrors.BadRequest("INVALID_LONG_CONTEXT_THRESHOLD", "GPT-6.1 Sol long-context threshold must be an integer between 1 and 2147483647")
	}
	return nil
}

// applyGroupLongContextPricingPolicy changes only GPT-6.1 Sol's model ladder.
// Call after model defaults are applied; explicit channel intervals and disabled
// ladders keep their own behavior. Clone shared cards to isolate groups.
func applyGroupLongContextPricingPolicy(model string, pricing *ModelPricing, group *Group) *ModelPricing {
	if pricing == nil || group == nil || !group.GPT61SolLongContextEnabled ||
		!supportsGPT61SolLongContextPlatform(group.Platform) ||
		!openai.IsGPT61SolModelSpelling(normalizeKnownOpenAICodexModel(model)) ||
		pricing.LongContextInputThreshold <= 0 ||
		validateGPT61SolLongContextThreshold(group.GPT61SolLongContextThreshold) != nil {
		return pricing
	}
	cloned := *pricing
	cloned.LongContextInputThreshold = group.GPT61SolLongContextThreshold
	return &cloned
}
