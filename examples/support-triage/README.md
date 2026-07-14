# support-triage — one agent, pure prompt

The smallest useful Kastor module: a single agent that classifies an inbound
support ticket into a category, a priority, and a one-line summary. No tools,
no orchestration — just a typed IO contract and a prompt, compiled to a
runnable LangGraph project.

| File | Purpose |
|------|---------|
| `kastor.hcl` | the model and the LangGraph codegen target |
| `triage.agent` | typed inputs (`subject`, `body`, optional `customer_tier`) and outputs (`category`, `priority`, `summary`) |
| `triage_system.prompt` | the classification rubric; requires exactly the agent's inputs |

## Validate and build

From the repository root:

```sh
kastor validate examples/support-triage/
kastor build examples/support-triage/
```

Build writes a self-contained Python project to
`examples/support-triage/gen/langgraph/` (generated output is reproducible
and never committed).

## Run the generated agent

Requires Python 3.11+ and an OpenAI API key (the module's `model "fast"` is
`openai` / `gpt-4o-mini` — swap the provider in `kastor.hcl` and rebuild to
use another vendor).

```sh
cd examples/support-triage/gen/langgraph
python -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt

export OPENAI_API_KEY=sk-...
python main.py triage --inputs '{
  "subject": "Charged twice this month",
  "body": "I upgraded to the annual plan on Friday and my card shows two charges of $228. I need one refunded before my statement closes.",
  "customer_tier": "pro"
}'
```

The agent returns its declared outputs as structured JSON:

```json
{
  "category": "billing",
  "priority": "high",
  "summary": "Customer sees two $228 charges after upgrading to the annual plan and needs one refunded."
}
```

`customer_tier` is optional — omit it and the prompt treats the customer as
free tier.

## See validation catch a mistake

The prompt's variables and the agent's IO contract are checked against each
other at compile time. Rename an input the prompt depends on:

```sh
sed -i '' 's/input "subject"/input "subject_line"/' examples/support-triage/triage.agent
kastor validate examples/support-triage/
```

```
triage.agent: agent.triage: system_prompt prompt.triage_system: variable "subject" is not an input or output of the agent
kastor: validation failed: 1 error
```

Undo the edit (`git checkout -- examples/support-triage/triage.agent`) and
`kastor validate` is green again — the same loop an AI writing this module
uses to self-correct.
