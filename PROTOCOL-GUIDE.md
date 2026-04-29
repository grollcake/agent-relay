# Agent Relay 한국어 가이드

이 문서는 Agent Relay를 프로젝트에 도입하거나 운영할 때 읽는 한국어 해설서입니다.
정식 배포 파일은 `bootstrap/` 아래에 있으며, 이 문서는 그 파일들의 의도와 사용 방법을 설명합니다.

## 1. Agent Relay란

Agent Relay는 서로 다른 AI 코딩 에이전트가 같은 프로젝트에서 작업을 이어받기 위한 파일 기반 규약입니다.

지원하려는 상황은 다음과 같습니다.

- Claude Code, Codex CLI, Cursor 같은 서로 다른 도구가 같은 프로젝트를 다룰 때
- 한 에이전트가 시작한 작업을 다른 에이전트나 다른 세션이 이어받을 때
- 세션이 끊긴 뒤 최근 작업 맥락을 빠르게 복구해야 할 때
- 분석, 구현, 테스트, 리뷰처럼 역할이 나뉘는 흐름에서 최소한의 인수인계가 필요할 때

Agent Relay는 오케스트레이터가 아닙니다. 에이전트를 실행하거나 병렬로 관리하지 않습니다.
역할은 단순합니다. 프로젝트 안에 작은 작업 기억과 인수인계 표준을 남기는 것입니다.

## 2. 핵심지침

Agent Relay의 핵심지침은 매 세션에서 에이전트가 실제로 따라야 하는 작업 흐름입니다.

1. 작업에 합류하거나 재개할 때는 `.agent-relay/PROTOCOL.md`의 Agent Relay 읽기 순서를 따릅니다.
2. 세션 에이전트 이름을 하나 정하고 `relay.log`에서 일관되게 씁니다.
3. 의미 있는 작업은 `.agent-relay/relay.log`에 기록합니다.
4. 다음 에이전트가 이슈와 `relay.log`만 보고 5분 안에 이어받기 어려울 때만 handoff를 만듭니다.
5. handoff를 만들 때는 긴 맥락을 파일로 분리하고 `relay.log`에 `path=`로 연결합니다.
6. 오래 유지될 사용자 지침, 제약, 선호, 관례, 보안 규칙, 금지사항만 `.agent-relay/GUIDANCE.md`에 기록합니다.
7. `.agent-relay/`에 비밀정보, 자격증명, 고객정보, 민감한 운영정보를 저장하지 않습니다.

Git 정책처럼 설치와 저장소 운영에 관한 규칙은 핵심지침 번호에 넣지 않고 별도 기타 규칙으로 둡니다.

## 3. 권장 파일 구조

사용자 프로젝트에 도입된 뒤의 권장 구조입니다.

```text
project-root/
├── AGENTS.md
├── CLAUDE.md                  # 선택
├── .codex/
│   └── instructions.md         # 선택
├── .cursor/
│   └── rules/
│       └── agent-relay.mdc     # 선택
└── .agent-relay/
    ├── PROTOCOL.md
    ├── VERSION
    ├── INDEX.md
    ├── GUIDANCE.md             # 장기 지침/제약 누적
    ├── relay.log
    ├── handoff/
    │   └── .gitkeep
    └── templates/
        ├── handoff.md
        └── task-done.md
```

실제 복사 원본은 이 저장소의 `bootstrap/` 트리입니다. `bootstrap/` 안의 경로는 목적지 프로젝트 기준 경로와 일치합니다.

## 4. 파일별 역할

| 파일 | 필수 여부 | 역할 |
|---|---:|---|
| `AGENTS.md` | 필수 | 모든 에이전트가 읽는 공통 작업 지침입니다. |
| `.agent-relay/PROTOCOL.md` | 필수 | Agent Relay의 최소 규칙입니다. |
| `.agent-relay/VERSION` | 필수 | 설치 버전입니다. 업데이트 시 기본 upstream과 비교하는 기준으로 씁니다. |
| `.agent-relay/INDEX.md` | 권장 | 프로젝트 식별 정보, 중요 파일, 열린 handoff를 빠르게 찾는 지도입니다. |
| `.agent-relay/GUIDANCE.md` | 누적 관리 | 세션을 넘어 유지할 사용자 지침, 제약, 금지사항을 담는 문서입니다. |
| `.agent-relay/relay.log` | 필수 | 추가 전용 작업 이벤트 로그입니다. |
| `.agent-relay/handoff/` | 필수 | 긴 인수인계 문서를 두는 디렉토리입니다. |
| `.agent-relay/templates/` | 권장 | 로그와 인수인계 작성 형식 예시입니다. |

## 5. 합류할 때 읽는 순서

새 세션, 새 에이전트, 오래된 작업 재개, handoff 인수 상황에서는 먼저 `AGENTS.md`를 읽습니다.
그 다음 이번 세션의 에이전트 이름을 정합니다.

형식:

```text
<에이전트>(<LLM 또는 버전>[, <역할>])
```

예:

```text
Codex(GPT-5.5)
Claude Code(Claude Sonnet 4.5)
Cursor(GPT-5.5, Reviewer)
```

LLM 또는 버전을 모르면 괄호 없이 에이전트 이름만 씁니다. 예: `Codex`, `Claude Code`, `Cursor`.

사람이 직접 수행한 작업에만 `Human`을 씁니다.

전체 순서는 다음과 같습니다.

```text
1. AGENTS.md
2. 세션 에이전트 이름 결정
3. README.md
4. .agent-relay/GUIDANCE.md
5. .agent-relay/PROTOCOL.md
6. .agent-relay/INDEX.md
7. .agent-relay/relay.log의 마지막 50줄
8. 참조된 handoff 파일이 있으면 읽음
```

같은 세션에서 연속 작업 중이라면 매 사용자 메시지마다 `relay.log`를 다시 읽지 않습니다.
세션 시작 때 읽고, 의미 있는 작업을 마칠 때 기록하는 흐름이면 충분합니다.

## 6. 기록해야 하는 작업

다음 작업은 `relay.log`에 결과를 남기는 것이 좋습니다.

- 코드 변경
- 문서 변경
- 설정 변경
- 테스트 실행
- 디버깅 시도
- 아키텍처 또는 설계 판단
- 이슈 상태 변경
- 다음 에이전트가 알아야 할 맥락 생성

단순 질의응답, 짧은 설명, 브레인스토밍은 보통 기록하지 않습니다.
사용자가 맥락 보존을 명시적으로 요청한 경우에는 예외로 기록할 수 있습니다.

## 7. relay.log 형식

`relay.log`는 추가 전용입니다. 기존 줄을 수정하지 않습니다.
정정이 필요하면 새 `CORRECTION` 이벤트를 추가합니다.

권장 이벤트 타입은 다음과 같습니다.

- `SESSION_START`
- `TASK_BEGIN`
- `TASK_DONE`
- `HANDOFF`
- `HANDOFF_RECEIVED`
- `HANDOFF_CLOSED`
- `CORRECTION`

의미 있는 작업을 시작하고, 세션이 끊기면 복구가 필요할 수 있는 경우 `TASK_BEGIN`을 남깁니다.

```text
<timestamp> | TASK_BEGIN | agent=<agent> | task=<MMDD-xxx> | summary=<task summary>
```

가장 자주 쓰는 `TASK_DONE` 형식은 다음과 같습니다.

```text
<timestamp> | TASK_DONE  | agent=<agent> | task=<MMDD-xxx> | result=<success|partial|failed|blocked> | files=<summary> | tests=<summary>
```

`agent`는 `<에이전트>(<LLM 또는 버전>[, <역할>])` 형식으로 적습니다. 예: `Codex(GPT-5.5)`, `Cursor(GPT-5.5, Reviewer)`. LLM 또는 버전을 모르면 `Codex`처럼 에이전트 이름만 씁니다. 사람이 직접 수행한 작업에만 `Human`을 씁니다.
같은 세션에서는 `agent=`에 들어가는 에이전트 이름을 일관되게 유지합니다.

`task`는 `TASK_BEGIN`, `TASK_DONE`, 관련 `HANDOFF`를 연결하는 짧은 키입니다. `TASK_BEGIN` 시점의 로컬 월일과 무작위 소문자 3글자로 만듭니다. 예: `0428-qmx`. 최근 릴레이 맥락에서 같은 `task=`가 아직 닫히지 않았으면 다시 생성합니다.

`TASK_DONE`은 작업 시도가 닫혔다는 뜻이고, `result`는 그 결과가 성공인지, 부분 완료인지, 실패인지, 차단 상태인지를 기록합니다.

예시:

```text
2026-04-28T16:20:00+09:00 | TASK_BEGIN | agent=Codex(GPT-5.5) | task=0428-qmx | summary="Align relay event names"
2026-04-28T16:40:00+09:00 | TASK_DONE  | agent=Codex(GPT-5.5) | task=0428-qmx | result=success | files=src/cube/state.ts, tests/cube-state.test.ts | tests=npm test passed
```

긴 설명은 `relay.log`에 직접 넣지 말고 `.agent-relay/handoff/` 아래 문서로 분리한 뒤 로그에서는 경로만 참조합니다.

handoff 파일을 만들 때는 `relay.log`에 `HANDOFF` 이벤트를 추가하고 `path=<handoff-file>`로 연결합니다.

예시:

```text
2026-04-28T16:40:00+09:00 | HANDOFF | agent=Codex(GPT-5.5) | task=0428-qmx | result=blocked | path=.agent-relay/handoff/0428-qmx--auth-refresh.md | summary="Auth refresh behavior needs follow-up analysis"
```

## 8. handoff 작성 기준

handoff는 모든 작업마다 쓰는 문서가 아닙니다.
다음 작업자가 `relay.log`와 관련 이슈만 보고 5분 안에 이어받기 어려울 때 작성합니다.

handoff가 필요한 경우:

- 다른 에이전트나 다른 역할이 이어서 작업해야 하는 경우
- 결과가 `partial`, `failed`, `blocked`인 경우
- 테스트, 리뷰, 분석을 다른 역할이 이어받아야 하는 경우
- 중요한 맥락이 로그 한 줄에 담기 어려운 경우

handoff가 필요하지 않은 경우:

- 작업이 완료되었고 다음 소유자가 없는 경우
- `relay.log` 한 줄로 충분히 설명되는 경우
- 별도 인수인계 없이 다음 작업자가 바로 진행할 수 있는 경우

파일 이름 형식:

```text
.agent-relay/handoff/<task>--<topic>.md
```

예시:

```text
.agent-relay/handoff/0428-qmx--cube-rotation.md
```

## 9. GUIDANCE 사용 기준

`GUIDANCE.md`는 기본 템플릿으로 복사됩니다.
부트스트랩 직후 억지로 채우는 문서가 아닙니다.

이 파일은 진행 상태, 현재 작업, 임시 계획, 다음 작업을 기록하는 곳이 아닙니다. 그런 정보는 `relay.log`와 handoff에 둡니다.

대신 Agent Relay를 사용하는 동안 사용자가 지속적으로 주는 장기 지침을 누적합니다.

기록 대상:

- 오래 유지되는 프로젝트 맥락
- 사용자 선호와 상시 지시
- 기술·제품·운영 제약
- 하지 말아야 할 것
- 보안·개인정보 규칙
- 코드/테스트/브랜치/작업 관례

처음에는 비어 있거나 자리표시자가 남아 있어도 됩니다. 이후 장기 지침이 생길 때 갱신합니다.

갱신 시점:

- 사용자가 앞으로도 지켜야 할 지침을 말했을 때
- 사용자가 하지 말아야 할 일을 지정했을 때
- 보안, 개인정보, 운영상 제약이 추가되었을 때
- 코드 스타일, 테스트 방식, 브랜치 방식 같은 반복 적용 규칙이 생겼을 때

갱신하지 않는 경우:

- 현재 작업 진행 상황
- 임시 계획
- 디버깅 중간 메모
- 한 번만 적용되는 요청

기록 예시:

- "새 런타임 의존성은 사용자 확인 없이 추가하지 말 것"
- "공개 API는 하위 호환성을 유지할 것"
- "이 저장소에서는 pnpm을 사용할 것"

기록하지 않는 예시:

- "이 테스트를 고쳐줘"
- "먼저 다른 구현을 시도해봐"
- "현재 의심 원인은 race condition"

## 10. 부트스트랩 절차

대상 프로젝트에 Agent Relay를 적용할 때는 `BOOTSTRAP.md`를 기준으로 진행합니다.
핵심 절차는 다음과 같습니다.

1. 부트스트랩 세션의 에이전트 이름을 정합니다.
2. `.agent-relay/`가 이미 있으면 아무 파일도 바꾸지 않고 중단 후 보고합니다.
3. 대상 프로젝트의 `AGENTS.md`, 도구별 지시 파일 존재 여부를 확인합니다.
4. `AGENTS.md`가 없으면 `bootstrap/AGENTS.md`를 복사합니다.
5. `AGENTS.md`가 이미 있으면 기존 내용을 보존하고 `bootstrap/AGENTS.md`의 `Agent Relay` 섹션만 병합합니다.
6. 도구별 지시 파일을 필요에 따라 병합합니다.
7. `bootstrap/.agent-relay/`를 복사합니다.
8. `relay.log`의 자리표시자 timestamp와 `agent=` 값을 현재 세션 정보로 바꿉니다.
9. `.agent-relay/VERSION`에 설치 버전을 기록합니다.
10. `INDEX.md`의 자리표시자를 프로젝트 정보로 채웁니다.
11. `GUIDANCE.md`는 장기 지침과 제약을 누적하는 문서라고 안내합니다.
12. 비밀정보가 `.agent-relay/`에 들어가지 않았는지 확인합니다.
13. Git 저장소라면 `.agent-relay/`가 커밋 대상인지 확인합니다.

## 11. 업데이트 절차

사용자가 "agent-relay 최신화해줘", "업데이트해줘", "sync 해줘"처럼 요청하면 새 부트스트랩이 아니라 업데이트 절차를 수행합니다.

1. `.agent-relay/VERSION`을 읽어 현재 설치 버전을 확인합니다.
2. 기본 upstream `https://github.com/grollcake/agent-relay`의 최신 `main`을 임시 위치에 가져와 현재 프로젝트와 비교합니다.
3. `AGENTS.md`는 최신 `bootstrap/AGENTS.md`의 `Agent Relay` 섹션과 비교해 현재 파일의 Agent Relay 섹션만 교체하거나 보강합니다.
4. `.agent-relay/PROTOCOL.md`와 `.agent-relay/templates/`는 로컬 수정이 없거나 안전히 구분될 때 최신 upstream으로 갱신합니다.
5. `.agent-relay/GUIDANCE.md`, `.agent-relay/relay.log`, `.agent-relay/handoff/`는 덮어쓰지 않습니다.
6. `.agent-relay/INDEX.md`는 프로젝트별 지도이므로 덮어쓰지 않고, 새 권장 섹션이 있으면 사용자 확인 후 보강합니다.
7. 업데이트가 성공하면 `.agent-relay/VERSION`을 최신 upstream의 `VERSION` 값으로 갱신하고 `relay.log`에 `TASK_DONE`을 추가합니다.

`AGENTS.md`에 Agent Relay 지시와 프로젝트 고유 지시가 섞여 있어 자동 분리가 어렵다면 파일을 바꾸지 않고 충돌로 보고합니다.

## 12. 도구별 지시 파일

도구별 파일은 선택입니다. 사용자가 해당 도구를 쓰는 흔적이 있을 때만 추가하거나 병합합니다.

| 조건 | 처리 |
|---|---|
| `CLAUDE.md`가 이미 있음 | 기존 내용을 보존하고 `bootstrap/CLAUDE.md`의 Agent Relay 포인터를 병합합니다. |
| `.codex/instructions.md`가 이미 있음 | 기존 내용을 보존하고 `bootstrap/.codex/instructions.md`의 포인터를 병합합니다. |
| `.cursor/` 디렉토리가 있음 | `.cursor/rules/agent-relay.mdc`를 추가합니다. 같은 파일이 있으면 병합합니다. |

해당 파일이나 디렉토리가 없으면 새로 만들지 않는 것이 기본입니다.
예외는 `.cursor/` 디렉토리가 이미 있는 경우의 Cursor rule 파일입니다.

## 13. 보안 규칙

`.agent-relay/`에는 다음 정보를 저장하지 않습니다.

- API 키
- 토큰
- 비밀번호
- 개인 키와 자격 증명
- 고객 데이터와 개인정보
- 민감한 내부 정보
- 운영 비밀

이 규칙은 `relay.log`, `INDEX.md`, `GUIDANCE.md`, handoff 문서 모두에 적용됩니다.

## 14. 기타 규칙: Git 정책

Git 저장소에서는 `.agent-relay/`를 커밋합니다.

Agent Relay는 에이전트 간 작업 인수인계를 위한 프로젝트 자산입니다. 릴레이 규칙, 프로젝트 지도, 작업 로그, handoff 문서는 다음 에이전트가 같은 저장소에서 이어받을 수 있어야 합니다.

따라서 `.agent-relay/`를 `.gitignore`에 추가하지 않습니다. 이미 무시 규칙이 있다면 제거합니다.

비밀정보는 커밋 여부와 무관하게 `.agent-relay/`에 저장하지 않습니다.

## 15. 사용자에게 줄 수 있는 짧은 호출문

매번 긴 규칙을 붙일 필요는 없습니다.

```text
Agent Relay 규칙에 따라 진행해.
```

더 명확히 지시하려면 다음처럼 말합니다.

```text
AGENTS.md와 .agent-relay/PROTOCOL.md 기준으로 진행해.
새 세션이면 Agent Relay의 읽기 순서를 먼저 따라줘.
```

영어 환경에서는 다음처럼 줄일 수 있습니다.

```text
Follow Agent Relay. If this is a new or resumed session, follow the Agent Relay read order before acting.
```

## 16. 배포 파일 위치

정식 템플릿과 원문 규칙은 다음 파일을 기준으로 합니다.

| 파일 | 용도 |
|---|---|
| `BOOTSTRAP.md` | 에이전트가 대상 프로젝트에 Agent Relay를 설치할 때 따르는 절차 |
| `bootstrap/AGENTS.md` | 대상 프로젝트 루트에 둘 공통 지시문 |
| `bootstrap/.agent-relay/PROTOCOL.md` | Agent Relay 정식 최소 규칙 |
| `bootstrap/.agent-relay/VERSION` | 설치 버전 템플릿 |
| `bootstrap/.agent-relay/INDEX.md` | 프로젝트 지도 템플릿 |
| `bootstrap/.agent-relay/GUIDANCE.md` | 장기 지침/제약 누적 템플릿 |
| `bootstrap/.agent-relay/relay.log` | 초기 로그 템플릿 |
| `bootstrap/.agent-relay/templates/task-done.md` | `TASK_DONE` 로그 형식 |
| `bootstrap/.agent-relay/templates/handoff.md` | handoff 문서 형식 |

## 17. 정리

Agent Relay는 큰 프로세스 관리 체계가 아니라, 에이전트가 작업을 끊기지 않게 넘기기 위한 최소 릴레이 습관입니다.

핵심 판단 기준은 하나입니다.

```text
다음 에이전트가 5분 안에 이어받을 수 있는가?
```

가능하면 `relay.log` 한 줄이면 충분합니다.
불가능하면 handoff 문서를 남깁니다.
