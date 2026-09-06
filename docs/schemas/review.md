# Review Schema

The review schema defines a human-in-the-loop checkpoint for agent-produced
outputs that must be confirmed, edited, or rejected by a human before they are
trusted or allowed to proceed.

**Schema**: [`schema/review/review.schema.json`](https://github.com/plexusone/multi-agent-spec/blob/main/schema/review/review.schema.json)

## Overview

A `ReviewBatch` is a set of `ReviewItem`s produced by one agent run. Each item
carries the agent's proposed value, its confidence, and why it was flagged; a
human then confirms it, replaces it with an edited value, or rejects it.

This is distinct from a [Loop](loop.md)'s `escalation: human`, which pauses a
retry loop once `max_attempts` is exhausted. A review checkpoint is a per-item
gate usable anywhere an agent's output needs a human decision before
acceptance — for example:

- Confirming or correcting an LLM-proposed gold-label candidate before it
  counts as ground truth for evaluation
- Approving a flagged finding before it is published
- Gating an action before it takes effect

## Review Status

| Status | Meaning |
|--------|---------|
| `pending` | Proposed, not yet reviewed (default) |
| `confirmed` | A human accepted the proposed value as-is |
| `edited` | A human replaced the proposed value with their own |
| `rejected` | A human rejected the item; no final value |

## Schema Definition

```json
{
  "id": "domain-classify-batch-001",
  "source": "domain-classifier",
  "created_at": "2026-09-06T00:00:00Z",
  "items": [
    {
      "id": "aiid-1621",
      "proposed_value": "LLM",
      "reason": "Data-exposure framing is genuinely ambiguous",
      "confidence": 0.4,
      "status": "pending"
    },
    {
      "id": "aiid-1664",
      "proposed_value": "Identity",
      "confidence": 0.6,
      "status": "confirmed",
      "final_value": "Identity",
      "reviewer": "jwang",
      "reviewed_at": "2026-09-06T01:00:00Z"
    }
  ]
}
```

## Fields

### ReviewBatch

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Batch identifier |
| `source` | string | No | What produced the batch (agent name, rubric id) |
| `created_at` | string (date-time) | No | When the batch was produced |
| `items` | ReviewItem[] | Yes | The individual review checkpoints |

### ReviewItem

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Identifies the subject being reviewed |
| `proposed_value` | any | Yes | What the agent proposed |
| `reason` | string | No | Why this item was produced or flagged |
| `confidence` | number | No | Agent confidence, `0`–`1` |
| `status` | ReviewStatus | No | Current review state (default `pending`) |
| `final_value` | any | No | Human-confirmed or human-edited value |
| `reviewer` | string | No | Who made the decision |
| `notes` | string | No | Reviewer remarks, especially on rejection |
| `reviewed_at` | string (date-time) | No | When the decision was made |

## Go SDK

```go
import multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"

batch := multiagentspec.NewReviewBatch("batch-001", "domain-classifier", time.Now())
batch.AddItem(multiagentspec.ReviewItem{
    ID:            "aiid-1621",
    ProposedValue: "LLM",
    Confidence:    0.4,
    Reason:        "ambiguous",
})

// A human confirms, edits, or rejects each item.
batch.Items[0].Confirm("jwang", time.Now())

if batch.AllResolved() {
    // proceed
}
```

`ReviewBatch.Validate()` checks batch/item ids are present and unique and that
confidence values fall in `[0, 1]`.
