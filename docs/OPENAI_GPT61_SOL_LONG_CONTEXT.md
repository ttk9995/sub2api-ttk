# GPT-6.1 Sol subscription long-context billing

Billing policy updated on October 21, 2026.

## Official announcement and scope

On October 12, 2026 at 07:45 in Asia/Shanghai (October 11 at 23:45 UTC),
Tibo (@thsottiaux) announced that GPT-6.1 Sol's default context window in
Codex and Work increased from 272k to 500k. The maximum context window
remains 1M. See the [announcement](https://x.com/thsottiaux/status/2109430343346581527).

His [follow-up](https://x.com/thsottiaux/status/2109430894604026039) at
07:48 in Asia/Shanghai explicitly says that long-context fees no longer
apply when GPT-6.1 Sol draws usage from the subscription. The user-provided
screenshot includes both posts. X blocks automated retrieval; the matching
post text and timestamps were also checked against a
[third-party archive](https://codex-reset.com/tibo/x).

500k is the default context window, not a replacement surcharge threshold.
The subscription exemption is not limited to requests below 500k.

## Deployment policy for API Key accounts

The operator reported on October 21 that many API Key upstreams have not
adopted the 500k policy. GPT-6.1 Sol therefore defaults to a 272,000-token
long-context billing threshold again. Both Go fallback cards and the
bundled catalog use 272k; a remote default card with the previous 500k
threshold is corrected to 272k as well.

Each OpenAI or composite group can opt into its own token threshold using
`gpt61_sol_long_context_enabled` and `gpt61_sol_long_context_threshold`.
The override is off by default, with 500000 prefilled as an editable value.
Turning it off preserves the saved value and restores the default ladder.
Valid thresholds are integers from 1 through 2147483647. The threshold
counts ordinary input, cache reads, and cache creation together; output
tokens do not count toward it. The surcharge applies strictly above the
threshold to the whole request, retaining input/cache 2x and output 1.5x.

The override only changes GPT-6.1 Sol's existing model ladder, including
recognized aliases. It does not enable long-context billing by itself.
Group/account billing controls, OAuth/SetupToken exemptions, disabled
ladders, administrator-defined channel intervals, and other models keep
their behavior. Flat group/channel prices retain their configured rates
and use the selected model threshold when their model ladder is active.
Account cost statistics continue to follow account/channel cost policies
independently of the group's customer billing threshold.

The group fields are persisted by migration
`247_group_gpt61_sol_long_context.sql`, copied when duplicating groups,
and included in API key authentication snapshots. Snapshot version 25
refreshes older cached group data. Saving a group invalidates its existing
authentication cache. Model Plaza derives its displayed tiers from the
same group policy used for billing.

## Sub2API behavior

- The exemption uses the actual upstream OpenAI OAuth or SetupToken
  credential, including the parent credential of an account shadow. The
  client's Sub2API API key does not determine the upstream billing method.
- It matches the final model sent upstream, including recognized model
  aliases and compact spellings. A billing label or response-model label
  alone does not grant the exemption.
- Model long-context multipliers stay disabled for eligible subscription
  requests even when account or group long-context billing is enabled.
  Response-model repricing and Free Fast recalculation retain this policy.
- API Key upstreams use the 272k default or the enabled group threshold
  described above. Other models, administrator-defined pricing intervals, and
  service-tier prices retain their existing behavior.
- Ordinary input, cached input, cache creation, and output keep their base
  prices when the model's subscription exemption applies.

Regression tests cover default and custom threshold boundaries, cached
input, service tiers, OAuth/SetupToken exemptions, gateway usage recording,
model aliases, shared-card immutability, explicit intervals, Model Plaza
tiers, group update semantics, and authentication cache round-trips.
These are local unit tests using stub repositories; they do not measure
a live account's upstream usage allowance.
