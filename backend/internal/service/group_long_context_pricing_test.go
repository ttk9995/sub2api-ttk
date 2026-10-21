//go:build unit

package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGPT61SolGroupThresholdBillingAndSchedule(t *testing.T) {
	bs := NewBillingService(nil, nil)
	resolver := NewModelPricingResolver(nil, bs)
	for _, threshold := range []int{200000, 400000, 500000, 700000} {
		for _, model := range []string{"gpt-6.1-sol", "openai/gpt-6.1-sol-max", "gpt-6.1-sol-openai-compact"} {
			for _, enabled := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%d/enabled=%t", model, threshold, enabled), func(t *testing.T) {
					group := &Group{ID: 12, Platform: PlatformOpenAI, LongContextPricingEnabled: true,
						GPT61SolLongContextEnabled: enabled, GPT61SolLongContextThreshold: threshold}
					wantThreshold := 272000
					if enabled {
						wantThreshold = threshold
					}
					for _, n := range []int{wantThreshold - 1, wantThreshold, wantThreshold + 1} {
						for _, tier := range []string{"", "priority", "flex"} {
							for _, withResolver := range []bool{false, true} {
								input := CostInput{Ctx: context.Background(), Model: model, Group: group,
									Tokens:         UsageTokens{InputTokens: n - 3000, CacheReadTokens: 2000, CacheCreationTokens: 1000, OutputTokens: 1000},
									RateMultiplier: 1, ServiceTier: tier}
								if withResolver {
									input.Resolver = resolver
								}
								cost, err := bs.CalculateCostUnified(input)
								require.NoError(t, err)
								require.Equal(t, n > wantThreshold, cost.LongContextBillingApplied)
								inputFactor, outputFactor, tierFactor := 1.0, 1.0, 1.0
								if n > wantThreshold {
									inputFactor, outputFactor = 2, 1.5
								}
								if tier == "priority" {
									tierFactor = 2
								} else if tier == "flex" {
									tierFactor = 0.5
								}
								require.InDelta(t, float64(n-3000)*2e-6*inputFactor*tierFactor, cost.InputCost, 1e-12)
								require.InDelta(t, 2000*0.1e-6*inputFactor*tierFactor, cost.CacheReadCost, 1e-12)
								require.InDelta(t, 1000*2.5e-6*inputFactor*tierFactor, cost.CacheCreationCost, 1e-12)
								require.InDelta(t, 0.01*outputFactor*tierFactor, cost.OutputCost, 1e-12)
							}
						}
					}
					schedule, err := bs.ResolveContextPricingSchedule(context.Background(), resolver, ContextPricingScheduleInput{Model: model, Group: group})
					require.NoError(t, err)
					require.Len(t, schedule.Tiers, 2)
					require.Equal(t, wantThreshold, *schedule.Tiers[0].MaxTokens)
					require.Equal(t, wantThreshold, schedule.Tiers[1].MinTokens)
					original, err := bs.GetModelPricing(model)
					require.NoError(t, err)
					require.Equal(t, 272000, original.LongContextInputThreshold, "分组配置不能污染共享价卡")
				})
			}
		}
	}
}

func TestGPT61SolGroupThresholdIsolation(t *testing.T) {
	bs := NewBillingService(nil, nil)
	resolver := NewModelPricingResolver(nil, bs)
	for _, platform := range []string{PlatformOpenAI, PlatformComposite} {
		group := &Group{Platform: platform, LongContextPricingEnabled: true,
			GPT61SolLongContextEnabled: true, GPT61SolLongContextThreshold: 600000}
		for _, source := range []string{PricingSourceLiteLLM, PricingSourceGroup, PricingSourceChannel} {
			pricing := &ModelPricing{InputPricePerToken: 3e-6, OutputPricePerToken: 12e-6,
				LongContextInputThreshold: 272000, LongContextInputMultiplier: 2, LongContextOutputMultiplier: 1.5}
			input := CostInput{Model: "gpt-6.1-sol", Group: group, Tokens: UsageTokens{InputTokens: 400000},
				RateMultiplier: 1, Resolver: resolver, Resolved: &ResolvedPricing{BasePricing: pricing, Source: source, longContextPricingEnabled: true}}
			cost, err := bs.CalculateCostUnified(input)
			require.NoError(t, err)
			require.InDelta(t, 1.2, cost.TotalCost, 1e-12)
			require.False(t, cost.LongContextBillingApplied)
			input.Resolved.Intervals = []PricingInterval{{MinTokens: 272000, InputMultiplier: pricingMultiplier(3)}}
			cost, err = bs.CalculateCostUnified(input)
			require.NoError(t, err)
			require.InDelta(t, 3.6, cost.TotalCost, 1e-12, "显式渠道区间仍优先")
		}
	}
	group := &Group{Platform: PlatformOpenAI, LongContextPricingEnabled: true, GPT61SolLongContextEnabled: true, GPT61SolLongContextThreshold: 600000}
	for _, model := range []string{"gpt-6-sol", "gpt-6-luna", "gpt-6-astra", "gpt-6.1-sol-preview"} {
		pricing := &ModelPricing{LongContextInputThreshold: 272000}
		require.Same(t, pricing, applyGroupLongContextPricingPolicy(model, pricing, group))
	}
	group.LongContextPricingEnabled = false
	cost, err := bs.CalculateTokenCostForRequest(TokenCostRequest{Model: "gpt-6.1-sol", Group: group,
		Resolver: resolver, Tokens: UsageTokens{InputTokens: 900000}, RateMultiplier: 1})
	require.NoError(t, err)
	require.False(t, cost.LongContextBillingApplied, "阈值配置不隐式开启长上下文计费")
}

func TestOpenAIRecordUsageGPT61SolCustomGroupThreshold(t *testing.T) {
	for _, withResolver := range []bool{false, true} {
		for _, threshold := range []int{400000, 500000, 650000} {
			for _, n := range []int{272001, threshold, threshold + 1} {
				for _, accountType := range []string{AccountTypeAPIKey, AccountTypeOAuth, AccountTypeSetupToken} {
					usageRepo := &openAIRecordUsageLogRepoStub{inserted: true}
					svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{}, &openAIRecordUsageSubRepoStub{}, nil)
					apiKey := &APIKey{ID: 1, Group: &Group{ID: 10, Platform: PlatformOpenAI, LongContextPricingEnabled: true,
						GPT61SolLongContextEnabled: true, GPT61SolLongContextThreshold: threshold}}
					if withResolver {
						svc.resolver = NewModelPricingResolver(nil, svc.billingService)
					}
					err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
						Result: &OpenAIForwardResult{RequestID: "custom-boundary", Model: "gpt-6.1-sol", Usage: OpenAIUsage{InputTokens: n, OutputTokens: 1000}},
						APIKey: apiKey, User: &User{ID: 2}, Account: &Account{ID: 3, Platform: PlatformOpenAI, Type: accountType,
							Extra: map[string]any{openAILongContextBillingEnabledKey: true}}})
					require.NoError(t, err)
					wantLong := accountType == AccountTypeAPIKey && n > threshold
					require.Equal(t, wantLong, usageRepo.lastLog.LongContextBillingApplied)
					inputFactor := 1.0
					if wantLong {
						inputFactor = 2
					}
					require.InDelta(t, float64(n)*2e-6*inputFactor, usageRepo.lastLog.InputCost, 1e-12)
				}
			}
		}
	}
}

func TestAdminGroupGPT61SolContextSettings(t *testing.T) {
	ctx := context.Background()
	threshold := 640000
	repo := &groupRepoStubForAdmin{createID: 1}
	svc := &adminServiceImpl{groupRepo: repo}
	group, err := svc.CreateGroup(ctx, &CreateGroupInput{Name: "custom", Platform: PlatformOpenAI,
		RateMultiplier: 1, GPT61SolLongContextEnabled: true, GPT61SolLongContextThreshold: &threshold})
	require.NoError(t, err)
	require.True(t, group.GPT61SolLongContextEnabled)
	require.Equal(t, threshold, group.GPT61SolLongContextThreshold)
	repo.getByID = group
	group, err = svc.UpdateGroup(ctx, 1, &UpdateGroupInput{Description: ptrString("unrelated")})
	require.NoError(t, err)
	require.True(t, group.GPT61SolLongContextEnabled)
	require.Equal(t, threshold, group.GPT61SolLongContextThreshold)
	group, err = svc.UpdateGroup(ctx, 1, &UpdateGroupInput{GPT61SolLongContextEnabled: boolPtr(false)})
	require.NoError(t, err)
	require.False(t, group.GPT61SolLongContextEnabled)
	require.Equal(t, threshold, group.GPT61SolLongContextThreshold)
	for _, invalid := range []int{0, -1, 2147483648} {
		_, err = svc.CreateGroup(ctx, &CreateGroupInput{GPT61SolLongContextThreshold: &invalid})
		require.Error(t, err)
		_, err = svc.UpdateGroup(ctx, 1, &UpdateGroupInput{GPT61SolLongContextThreshold: &invalid})
		require.Error(t, err)
	}
	group, err = svc.CreateGroup(ctx, &CreateGroupInput{Name: "default", Platform: PlatformOpenAI, RateMultiplier: 1})
	require.NoError(t, err)
	require.False(t, group.GPT61SolLongContextEnabled)
	require.Equal(t, 500000, group.GPT61SolLongContextThreshold)
}
