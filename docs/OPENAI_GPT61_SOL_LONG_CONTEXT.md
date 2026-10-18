# GPT-6.1 Sol subscription long-context billing

Verified on October 19, 2026.

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

The operator confirmed on October 19, 2026 that all upstreams in this
deployment use subscriptions, including upstreams registered as API Key
accounts. There is no official API business in this deployment. API Key
is therefore an integration method, not evidence of official API billing.

At the operator's request, the GPT-6.1 Sol default price card for API Key
accounts now applies no model long-context surcharge through 500,000 input
tokens inclusive. Above 500,000, existing input/cache multipliers (2x),
output multipliers (1.5x), and account/group controls continue to apply.
This 500k threshold is deployment policy; the announcement itself removes
subscription long-context fees without specifying a new fee threshold.

The bundled catalog and both Go fallback price cards use 500k. The billing
policy also migrates GPT-6.1 Sol's legacy 272k default threshold to 500k
after loading a remote catalog, so stale remote data cannot restore early
surcharges. Disabled ladders, other explicit thresholds, channel/group
custom prices, and other models retain their existing behavior.

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
- API Key upstreams use the deployment's 500k default threshold described
  above. Other models, administrator-defined pricing intervals, and
  service-tier prices retain their existing behavior.
- Ordinary input, cached input, cache creation, and output keep their base
  prices when the model's subscription exemption applies.

OAuth/SetupToken requests were already exempt. The regression update
extends boundary coverage around 272k and 500k, plus 900k, for OAuth,
SetupToken, and API Key upstreams, with and without cached input. Price-card
tests cover bundled, stale remote, and fallback sources, recognized model
aliases, shared-card immutability, and preserved custom pricing. These are
local unit tests using stub repositories; they do not measure a live
account's upstream usage allowance.
