# Agent Benchmark Architecture Research: Patterns for UT-Bench 3D Evaluation

## Overview

Comparative analysis of 5 open-source agent benchmark projects, focusing on architectural patterns borrowable for UT-Bench's **Agent × Model × Skill** 3D evaluation design.

---

## 1. SWE-bench — Evaluation Harness Architecture

### Core Pattern: Containerized Patch-and-Test Pipeline

SWE-bench uses a deterministic, binary evaluation pipeline:

```
Issue Text → Model generates Patch → Apply Patch in Docker → Run Test Suite → Pass/Fail
```

### Key Architectural Patterns

| Pattern | Implementation | Problem Solved |
|---------|---------------|----------------|
| **Layered Docker Image** | Base → Environment → Instance images with tiered caching (`--cache_level`) | Reproducibility + speed; common layers reused across runs |
| **Binary Resolution Grading** | Instance resolved iff ALL `FAIL_TO_PASS` tests pass AND ALL `PASS_TO_PASS` tests still pass | Eliminates ambiguity; no partial credit noise |
| **Test Category Bifurcation** | `FAIL_TO_PASS` (bug reproduction) vs `PASS_TO_PASS` (regression guard) | Separates "did it fix the bug?" from "did it break anything?" |
| **Gold Standard Validation** | `--predictions_path gold` runs the harness against known-correct patches | Verifies the evaluation infrastructure itself is sound |
| **Dataset Variant Abstraction** | Same harness evaluates Full/Lite/Verified/Multimodal via `--dataset_name` | One harness, multiple evaluation dimensions |
| **Parallel Worker Pool** | `--max_workers` with Docker container-per-instance | Horizontal scaling of evaluation |

### Scoring

- **Primary metric**: Resolution Rate = (Resolved / Submitted) × 100%
- **No MAP metric** — SWE-bench does not use Mean Average Precision. The user's "MAP" reference likely refers to the general evaluation protocol, not a specific metric.
- **No partial credit** — each instance is binary resolved/not-resolved
- **Cost & step tracking** — Leaderboard also tracks cost and step count as first-class dimensions

### What UT-Bench Can Borrow

- **Test category bifurcation** maps directly to UT-Bench's skill dimension: "did the unit test get generated?" (analogous to FAIL_TO_PASS) + "is the generated test itself correct?" (analogous to PASS_TO_PASS)
- **Layered caching** for expensive evaluation environments
- **Gold standard validation** pattern for benchmark integrity checks
- **Multi-variant evaluation** via same harness maps to Agent × Model × Skill matrix

---

## 2. SWE-agent — Agent-Computer Interface (ACI)

### Core Pattern: Interface-Centric Agent Design

SWE-agent's central insight: **the interface between agent and environment matters more than the agent logic itself**. The ACI is a first-class, configurable abstraction.

```
Agent ←→ ACI (YAML-configured tools + format) ←→ Computer Environment
```

### Key Architectural Patterns

| Pattern | Implementation | Problem Solved |
|---------|---------------|----------------|
| **ACI as Configurable Abstraction** | Single YAML file governs: available tools, parse function, prompt templates, history processors | Swap entire agent behavior without code changes |
| **Tool Bundles** | Directories of command definitions loaded by path; composable via YAML list | Mix-and-match tool sets for different tasks |
| **Parse Function Polymorphism** | `function_calling` / `thought_action` / `json` / `xml` — same tools, different action formats | Agent-agnostic: works with any model's output style |
| **Guardrail-by-Design** | Linter auto-runs on edits (blocks syntax errors); file viewer shows 100 lines/turn; search lists files only (not context) | Prevents common agent failure modes at the interface level |
| **Jinja2 Template System** | System, instance, next-step, and no-output templates all configurable | Full control over agent perception |
| **History Processor Chain** | Pluggable processors modify conversation history before LLM calls | Context window management, caching optimization |

### Trajectory Format

- **First-class artifact** stored in dedicated `trajectories/` directory
- Captures full action-observation sequence
- Used for both debugging and training data generation
- v1.1.0 released "10s of thousands of training trajectories"

### What UT-Bench Can Borrow

- **YAML-configured evaluation interface** — UT-Bench could define evaluation protocols as YAML configs, allowing different skill dimensions to have different tool sets and scoring criteria without code changes
- **Guardrail-by-design** — Build quality checks into the evaluation interface rather than post-hoc filtering
- **Parse function polymorphism** — Same evaluation harness can score outputs from agents with different output formats
- **Trajectory as first-class artifact** — Store full execution traces for multi-dimensional analysis

---

## 3. OpenHands — EventStream Architecture & Runtime Abstraction

### Core Pattern: Event-Sourced Agent Orchestration

OpenHands V1 (SDK) uses immutable event-sourcing as the central architectural pattern:

```
ConversationState (single source of truth)
  └── Append-only EventLog
       ├── LLMConvertibleEvents (visible to LLM)
       │    ├── MessageEvent, ActionEvent, ObservationEvent
       │    └── SystemPromptEvent, CondensationSummaryEvent
       └── Internal Events (not LLM-visible)
            ├── ConversationStateUpdateEvent, PauseEvent
            └── CondensationRequest, Condensation
```

### Key Architectural Patterns

| Pattern | Implementation | Problem Solved |
|---------|---------------|----------------|
| **Event-Sourced State** | All state = immutable events + metadata snapshot; recovery via replay from `start_id` | Crash recovery, audit trail, replay debugging |
| **Action-Observation Causality** | Every Observation carries a `cause` field linking back to originating Action | Precise action-effect pairing; pending action resolution |
| **Workspace Abstraction** | `BaseWorkspace` ABC → `LocalWorkspace` (subprocess) / `RemoteWorkspace` (HTTP to agent-server) | Seamless local-to-containerized migration with same API |
| **Opt-In Sandboxing** | V1 unified agent+tool in single process by default; containerize only when needed | Eliminates V0's dual-process divergence; faster local dev |
| **Conversation Factory** | `Conversation(agent, workspace)` resolves to local or remote execution | Same code path for prototyping and production |
| **Hierarchical Delegation** | Sub-agents run as independent conversations sharing parent metrics | Isolated subtask execution with unified cost tracking |
| **Stateless Agent + Callback** | Agents are immutable specs; emit events via `on_event()` callback | Security interleaving, incremental execution, real-time streaming |
| **Dual-Path Persistence** | Metadata → `base_state.json` (rewritten); Events → individual JSON files (append-only) | Efficient: metadata fast to load, events never rewritten |

### Sandbox Isolation (V1 Design)

- **Default**: No sandbox — agent and tools in same process (fast prototyping)
- **Production**: Docker containers with dedicated filesystem, network policies, resource limits
- **Key insight**: V0 forced universal sandboxing → dual-process state divergence → V1 made sandboxing opt-in, eliminating the problem

### What UT-Bench Can Borrow

- **Event-sourced evaluation** — Store every evaluation step as immutable events; enables replay, audit, and multi-dimensional post-hoc analysis across Agent × Model × Skill
- **Action-Observation causality** — Link each score back to the specific agent action that produced it; essential for skill-level debugging
- **Workspace abstraction** — Same evaluation logic runs locally (dev) or containerized (CI); swap `LocalWorkspace` ↔ `RemoteWorkspace`
- **Dual-path persistence** — Fast metadata access + append-only event log; ideal for large-scale 3D evaluation result storage
- **Hierarchical delegation pattern** — Sub-evaluations (per-skill) can run as independent conversations sharing parent (per-agent) metrics → natural 3D structure

---

## 4. DeepEval — Metric System Design

### Core Pattern: Pytest-Native Composable Metrics

DeepEval treats LLM evaluation as a **testing discipline** with pluggable, self-contained metrics:

```
BaseMetric (abstract)
  └── measure(test_case) → sets .score (0-1) + .reason (text)
       ├── LLM-as-Judge metrics (uses LLM to evaluate)
       ├── Statistical metrics (deterministic)
       └── Local NLP model metrics (no API calls)
```

### Key Architectural Patterns

| Pattern | Implementation | Problem Solved |
|---------|---------------|----------------|
| **Strategy Pattern for Metrics** | Every metric implements `.measure()` → `.score` + `.reason`; interchangeable without test code changes | Swap evaluation strategies freely |
| **Uniform Scoring Contract** | All metrics normalize to 0–1 with threshold-based pass/fail + human-readable reason | Cross-metric comparison and composition |
| **Metric Composition** | Pass a list of metrics to `assert_test()` or `evaluate()`; all must pass (conjunctive) | Multi-dimensional evaluation in one call |
| **Aggregation via Composite Metrics** | RAGAS metric averages 4 sub-metrics; custom composites possible | Collapse multiple dimensions into interpretable scores |
| **Three Evaluation Backends** | LLM-as-Judge / Statistical / Local NLP — same metric concept, different computation | Flexibility in cost vs accuracy tradeoffs |
| **Decorator-Based Tracing** | `@observe(metrics=[...])` wraps any function; creates hierarchical span tree | Zero-code-change observability; component-level evaluation |
| **Singleton Test Run Manager** | `global_test_run_manager` accumulates results across multiple `evaluate()` calls | Consistent state across a full evaluation session |
| **Disk-Based Caching** | `global_test_run_cache_manager` prevents redundant LLM API calls | Cost reduction for repeated evaluations |
| **Dual Execution Mode** | `assert_test()` (pytest/CI) vs `evaluate()` (notebook/exploratory) | Same metrics in automated and interactive contexts |

### Multi-Dimensional Evaluation Taxonomy

| Category | Metrics | UT-Bench Analogy |
|----------|---------|-------------------|
| **Agentic** | Task Completion, Tool Correctness, Step Efficiency, Plan Adherence | Agent dimension (how well does the agent perform?) |
| **RAG** | Answer Relevancy, Faithfulness, Contextual Recall/Precision | Skill dimension (how well is the specific skill executed?) |
| **Multi-Turn** | Knowledge Retention, Turn Relevancy, Role Adherence | Cross-skill dimension (coherence across skills) |
| **Custom** | G-Eval, DAG | Model dimension (model-specific quality assessment) |

### What UT-Bench Can Borrow

- **Strategy Pattern** — Each skill's evaluation is a metric with `.measure()` → `.score` + `.reason`; swap evaluation strategies per skill without changing framework
- **Uniform scoring contract** — All dimensions (Agent, Model, Skill) normalize to same scale, enabling matrix operations
- **Metric composition for 3D** — Agent-level metrics + Model-level metrics + Skill-level metrics composed conjunctively or via aggregation
- **Three-backend pattern** — Same metric concept evaluated via LLM-as-Judge, static analysis, or runtime execution → maps to UT-Bench's need for multiple evaluation methods
- **Decorator-based tracing** — Attach evaluation to specific agent actions/skills without code intrusion
- **Session-level result accumulation** — Singleton manager collects results across the full Agent × Model × Skill matrix

---

## 5. mini-swe-agent — Minimal Design Philosophy

### Core Pattern: Agent-Environment-Model with Radical Simplification

```
DefaultAgent(~100 LOC) + LitellmModel + LocalEnvironment → 74% on SWE-bench Verified
```

### Key Architectural Patterns

| Pattern | Implementation | Problem Solved |
|---------|---------------|----------------|
| **Bash-Only Interface** | No tool-calling interface; agent uses bash directly | Works with literally any model; eliminates scaffold-specific overfitting |
| **Stateless Subprocess Execution** | `subprocess.run` per action; no persistent shell | Maximum stability; trivial sandboxing (swap subprocess → docker exec) |
| **Linear History = Prompt** | Trajectory IS the message sequence; no transformation | Debugging is trivial; fine-tuning data is ready-made |
| **Three-Component Decomposition** | Agent + Environment + Model as clean separable classes | Each dimension independently swappable |
| **Trivial Sandbox Swapping** | Since every action is independent subprocess.run, swap to Docker/Podman/Singularity with one change | Environment-agnostic evaluation |

### What This Reveals About UT-Bench

mini-swe-agent proves that **complex scaffolding is not necessary for strong results**. The model quality dominates. This has a direct implication for UT-Bench's 3D design:

- **Agent dimension**: A simple baseline agent (like mini) should be included to isolate model effects from scaffold effects
- **Model dimension**: The simple scaffold ensures model quality is what's being measured, not scaffold-specific optimizations
- **Skill dimension**: Different skills may need different levels of scaffolding complexity — some skills may work with bash-only, others need specialized tools

---

## Cross-Project Pattern Synthesis for UT-Bench 3D Evaluation

### Pattern 1: Layered Abstraction (from SWE-bench + OpenHands)

```
Evaluation Layer (scoring logic, metric composition)
  └── Harness Layer (test execution, sandbox management)
       └── Instance Layer (per-task setup, patch application)
```

**For UT-Bench**: Replace "Instance" with "EvaluationCube" — each (Agent, Model, Skill) triple gets its own isolated evaluation context.

### Pattern 2: Configurable Evaluation Protocol (from SWE-agent)

```yaml
evaluation:
  dimensions:
    - name: agent
      variants: [mini, swe-agent, openhands]
    - name: model
      variants: [gpt-4o, claude-3.7, glm-5]
    - name: skill
      variants: [unit-test-gen, bug-fix, refactoring]
  metrics:
    - type: execution  # FAIL_TO_PASS / PASS_TO_PASS
      threshold: 1.0
    - type: llm_judge  # DeepEval-style
      threshold: 0.7
    - type: static     # lint, coverage
      threshold: 0.8
```

### Pattern 3: Event-Sourced 3D Matrix (from OpenHands)

```
For each (Agent_i, Model_j, Skill_k):
  → Create ConversationState
  → Run evaluation as event stream
  → Action-Observation pairs tagged with [agent, model, skill]
  → Immutable event log enables post-hoc slicing along any dimension
```

### Pattern 4: Composable Metric System (from DeepEval)

```
BaseMetric
  ├── AgentMetric (measures agent capability dimension)
  │    ├── TaskCompletionMetric
  │    ├── ToolCorrectnessMetric
  │    └── StepEfficiencyMetric
  ├── ModelMetric (measures model quality dimension)
  │    ├── CodeQualityMetric
  │    └── ReasoningMetric
  └── SkillMetric (measures skill-specific dimension)
       ├── UnitTestCoverageMetric
       ├── UnitTestCorrectnessMetric
       └── UnitTestRelevanceMetric

3DAggregateMetric = f(AgentMetric, ModelMetric, SkillMetric)
```

### Pattern 5: Sandbox Strategy Spectrum (from all projects)

| Approach | Source | When to Use in UT-Bench |
|----------|--------|------------------------|
| No sandbox (in-process) | OpenHands V1 default | Fast iteration, model-only evaluation |
| Docker per evaluation | SWE-bench | Full reproducibility, CI/CD |
| `subprocess.run` swappable | mini-swe-agent | Flexible; swap to Docker for production |
| Opt-in containerization | OpenHands V1 | Development speed + production safety |

**For UT-Bench**: Start with subprocess-based (mini-swe-agent style) for speed; graduate to Docker for reproducible CI runs. Same evaluation logic, different sandbox backend.

### Pattern 6: Trajectory as Observable Artifact (from SWE-agent + OpenHands)

- **SWE-agent**: Trajectory = prompt sequence (linear, inspectable)
- **OpenHands**: Trajectory = event stream (structured, replayable, causal links)
- **DeepEval**: Trajectory = span tree (decorator-based, hierarchical)

**For UT-Bench**: Use OpenHands-style event streams with DeepEval-style metric attachment:
- Each evaluation step emits an event
- Events are tagged with [agent, model, skill] dimensions
- Metrics attach to events (not the other way around)
- Full matrix reconstruction from any slice

---

## Recommended Architecture for UT-Bench

```
┌─────────────────────────────────────────────────────┐
│                  UT-Bench Framework                  │
├──────────────┬──────────────┬───────────────────────┤
│  Dimension   │  Dimension   │  Dimension            │
│  Agent       │  Model       │  Skill                │
│  (how)       │  (who)       │  (what)               │
├──────────────┴──────────────┴───────────────────────┤
│           Composable Metric System                    │
│  (DeepEval Strategy Pattern + Uniform Scoring)       │
├─────────────────────────────────────────────────────┤
│           Event-Sourced Evaluation Engine             │
│  (OpenHands EventStream + Causal Action-Observation) │
├─────────────────────────────────────────────────────┤
│           Configurable Harness (YAML-driven)          │
│  (SWE-agent ACI config + SWE-bench layered Docker)   │
├─────────────────────────────────────────────────────┤
│           Pluggable Sandbox Runtime                   │
│  (mini-swe-agent subprocess → Docker promotion)      │
└─────────────────────────────────────────────────────┘
```
