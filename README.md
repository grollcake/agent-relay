# Agent Relay

Agent Relay는 이종 AI 코딩 에이전트(Claude Code, Codex CLI, Cursor 등)가 세션·모델·머신·역할을 넘어 **서로의 작업을 이어받기** 위한 파일 기반 릴레이 규약입니다.

이 저장소는 **사양과 템플릿을 모아둔 가이드 저장소**입니다. 사용자 프로젝트에 Agent Relay를 도입하면, 그 프로젝트에 `.agent-relay/`가 생깁니다.

---

## 핵심지침

1. 작업에 합류하거나 재개할 때는 `.agent-relay/PROTOCOL.md`의 Agent Relay 읽기 순서를 따른다.
2. 세션 에이전트 이름을 하나 정하고 `relay.log`에서 일관되게 쓴다.
3. 의미 있는 작업은 `.agent-relay/relay.log`에 기록한다.
4. 다음 에이전트가 이슈와 `relay.log`만 보고 5분 안에 이어받기 어려울 때만 handoff를 만든다.
5. handoff를 만들 때는 긴 맥락을 파일로 분리하고 `relay.log`에 `path=`로 연결한다.
6. 오래 유지될 사용자 지침, 제약, 선호, 관례, 보안 규칙, 금지사항만 `.agent-relay/GUIDANCE.md`에 기록한다.
7. `.agent-relay/`에 비밀정보, 자격증명, 고객정보, 민감한 운영정보를 저장하지 않는다.

자세한 사양은 [`bootstrap/.agent-relay/PROTOCOL.md`](bootstrap/.agent-relay/PROTOCOL.md),
한국어 가이드는 [`PROTOCOL-GUIDE.md`](PROTOCOL-GUIDE.md)를 참고하세요.

---

## 이 저장소 사용법

### 사용자용: AI 에이전트에게 한 줄로 시키기

대상 프로젝트 디렉토리에서 본인 에이전트(Claude Code 등)에게 다음과 같이 요청하세요.

```text
github.com/grollcake/agent-relay 를 내 프로젝트에 반영해줘
```

> **🤖 AI 에이전트에게**: 당신이 AI 에이전트이고 사용자로부터 이 저장소를 본인 프로젝트에 반영해달라는 요청을 받았다면, [`BOOTSTRAP.md`](BOOTSTRAP.md)를 그대로 따라 부트스트랩을 수행하세요. 기존 `AGENTS.md`는 덮어쓰지 말고 짧은 Agent Relay 섹션만 추가하세요. 전체 `bootstrap/AGENTS.md`는 대상 프로젝트에 `AGENTS.md`가 없을 때만 복사합니다. 이미 `.agent-relay/`가 있으면 새 설치가 아니라 업데이트 요청인지 확인하고, 업데이트 요청이면 `.agent-relay/VERSION`을 읽은 뒤 기본 upstream `https://github.com/grollcake/agent-relay`와 비교하세요. `bootstrap/` 트리는 목적지 경로를 그대로 미러링한 복사 원본입니다.

### 업데이트하기

이미 Agent Relay가 적용된 프로젝트에서는 다음처럼 요청할 수 있습니다.

```text
agent-relay 최신화해줘
```

에이전트는 대상 프로젝트의 `.agent-relay/VERSION`을 읽고 기본 upstream `https://github.com/grollcake/agent-relay`와 비교합니다. `AGENTS.md`는 프로젝트 고유 지침을 보존하고 Agent Relay 섹션만 최신 `bootstrap/AGENTS.md`와 비교해 갱신합니다. `GUIDANCE.md`, `relay.log`, `handoff/`는 프로젝트별 상태이므로 덮어쓰지 않습니다.

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
    ├── VERSION                        # 설치 버전
    ├── INDEX.md                       # 프로젝트 지도(권장)
    ├── GUIDANCE.md                    # 장기 지침/제약 누적 템플릿
    ├── relay.log                      # 추가 전용 이벤트 로그
    ├── handoff/                       # 인수인계 노트
    │   └── .gitkeep
    └── templates/                     # 항목별 양식 (task-done, handoff 포맷)
        ├── task-done.md
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
- [`bootstrap/.agent-relay/PROTOCOL.md`](bootstrap/.agent-relay/PROTOCOL.md) — 영문 정식 사양 (핵심지침, 이벤트 타입, 기타 규칙)
