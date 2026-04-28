# Agent Relay

Agent Relay는 이종 AI 코딩 에이전트(Claude Code, Codex CLI, Cursor 등)가 세션·모델·머신·역할을 넘어 **서로의 작업을 이어받기** 위한 파일 기반 릴레이 규약입니다.

이 저장소는 **사양과 템플릿을 모아둔 가이드 저장소**입니다. 사용자 프로젝트에 Agent Relay를 도입하면, 그 프로젝트에 `.agent-relay/`가 생깁니다.

---

## 핵심 3원칙

1. **Read when joining** — 합류/재개 시 최근 맥락을 읽는다.
2. **Record when finishing** — 의미 있는 작업 후 `relay.log`에 기록한다.
3. **Handoff only when needed** — 다음 에이전트가 5분 안에 못 이어받을 때만 handoff를 쓴다.

자세한 사양은 [`bootstrap/.agent-relay/PROTOCOL.md`](bootstrap/.agent-relay/PROTOCOL.md),
한국어 가이드는 [`PROTOCOL-GUIDE.md`](PROTOCOL-GUIDE.md)를 참고하세요.

---

## 이 저장소 사용법

### 사용자용: AI 에이전트에게 한 줄로 시키기

대상 프로젝트 디렉토리에서 본인 에이전트(Claude Code 등)에게 다음과 같이 요청하세요.

```text
github.com/grollcake/agent-relay 를 내 프로젝트에 반영해줘
```

> **🤖 AI 에이전트에게**: 당신이 AI 에이전트이고 사용자로부터 이 저장소를 본인 프로젝트에 반영해달라는 요청을 받았다면, [`BOOTSTRAP.md`](BOOTSTRAP.md)를 그대로 따라 부트스트랩을 수행하세요. 기존 `AGENTS.md`는 덮어쓰지 말고 Agent Relay 섹션만 머지하고, 이미 `.agent-relay/`가 있으면 진행을 멈추고 사용자에게 보고하세요. `bootstrap/` 트리는 목적지 경로를 그대로 미러링한 복사 원본입니다.

---

## 부트스트랩 트리

```text
bootstrap/
├── AGENTS.md                          # 프로젝트 루트로 복사 (필수)
├── CLAUDE.md                          # CLAUDE.md가 이미 있으면 머지 (선택)
├── .codex/
│   └── instructions.md                # 이미 있으면 머지 (선택)
├── .cursor/
│   └── rules/
│       └── agent-relay.mdc            # .cursor/ 디렉토리가 있으면 신규 (선택)
└── .agent-relay/                      # 프로젝트 루트로 복사 (필수)
    ├── PROTOCOL.md                    # Agent Relay 규칙(필수)
    ├── INDEX.md                       # 프로젝트 지도(권장)
    ├── GUIDANCE.md                    # 장기 지침/제약 누적 템플릿
    ├── relay.log                      # 추가 전용 이벤트 로그
    ├── handoff/                       # 인수인계 노트
    │   └── .gitkeep
    └── templates/                     # 항목별 양식 (work-end, handoff 포맷)
        ├── work-end.md
        └── handoff.md
```

`bootstrap/` 안의 모든 경로는 **목적지 경로 그대로**입니다. 도구별 지시 파일(`CLAUDE.md`/`.codex/`)은 사용자가 그 도구를 쓰는 경우에만 머지되고, 안 쓰면 생성되지 않습니다.

---

## 매 프롬프트마다 쓰기 좋은 한 줄

```text
Agent Relay 규칙에 따라 진행해.
```

조금 더 명확하게:

```text
AGENTS.md와 .agent-relay/PROTOCOL.md 기준으로 진행해.
새 세션이면 Agent Relay의 읽기 순서를 먼저 따라줘.
```

---

## 더 읽기

- [`BOOTSTRAP.md`](BOOTSTRAP.md) — 에이전트용 부트스트랩 절차 (단계별, 머지/중단 분기 포함)
- [`PROTOCOL-GUIDE.md`](PROTOCOL-GUIDE.md) — 한국어 가이드 (운영 원칙, 부트스트랩 절차, 파일별 역할 정리)
- [`bootstrap/.agent-relay/PROTOCOL.md`](bootstrap/.agent-relay/PROTOCOL.md) — 영문 정식 사양 (Rule 1~3, 이벤트 타입, 안전 규칙)
