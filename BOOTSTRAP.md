# Bootstrap Guide

이 문서는 **AI 에이전트가 사용자 프로젝트에 Agent Relay를 적용할 때 따라야 하는 절차**입니다.
사람이 직접 따라할 수도 있지만, 기본은 에이전트가 이 문서를 읽고 자동 수행합니다.

저장소 개요와 사용자용 한 줄 프롬프트는 [`README.md`](README.md)를 참조하세요.

---

## 사전 확인

부트스트랩을 시작하기 전, 사용자 프로젝트(이하 `<project>`)의 다음 상태를 확인합니다.

- `<project>/AGENTS.md` 존재 여부
- `<project>/.agent-relay/` 디렉토리 존재 여부
- `<project>/CLAUDE.md`, `<project>/.codex/instructions.md` 존재 여부 (있으면 머지 대상)
- `<project>/.cursor/` 디렉토리 존재 여부 (있으면 룰 파일 신규 추가 대상)
- `<project>/README.md` 존재 여부 (없어도 진행 가능, 단계 6 참고)

---

## 부트스트랩 절차

### 0. 세션 에이전트 이름 결정

부트스트랩을 수행하는 주체 이름을 먼저 정합니다.

- 형식: `<에이전트>(<LLM 또는 버전>[, <역할>])`
- 예: `Codex(GPT-5.5)`, `Claude Code(Claude Sonnet 4.5)`, `Cursor(GPT-5.5, Reviewer)`
- LLM 또는 버전을 모르면 괄호를 쓰지 않고 에이전트 이름만 사용합니다. 예: `Codex`, `Claude Code`, `Cursor`
- 사람이 직접 수행하면 `Human`을 사용합니다.

이 이름은 `relay.log`의 `agent=` 값에 사용합니다.

### 1. `.agent-relay/` 디렉토리 처리

`<project>/.agent-relay/` 존재 여부에 따라 분기합니다.

- **이미 있으면**: **부트스트랩을 중단**하고 사용자에게 보고합니다.
  - 이미 Agent Relay가 적용된 프로젝트일 가능성이 큽니다.
  - 사용자 확인 없이 기존 `AGENTS.md`, `relay.log`, `handoff/`, `INDEX.md`를 덮어쓰거나 머지하지 않습니다.
  - 사용자가 명시적으로 "재초기화" 또는 "특정 파일만 갱신"을 요청한 경우에만, 해당 범위로 한정해 진행합니다.
- **없으면**: 이후 단계로 진행합니다.

### 2. `AGENTS.md` 처리

`<project>/AGENTS.md` 존재 여부에 따라 분기합니다.

- **없으면**: `bootstrap/AGENTS.md`를 그대로 `<project>/AGENTS.md`로 복사합니다.
- **이미 있으면**: 덮어쓰지 말고 **머지**합니다.
  - 기존 내용 전체를 보존합니다.
  - `bootstrap/AGENTS.md`에서 다음 5개 섹션을 추가합니다.
    - `Core Rules`
    - `When Starting a New Session`
    - `During Work`
    - `After Meaningful Work`
    - `Handoff`
  - 동일한 주제의 섹션이 이미 있으면, 누락된 항목만 보강합니다.
  - 기존 프로젝트의 규칙·역할·지시문을 임의로 삭제하거나 재작성하지 않습니다.
  - 머지 결과 상단에 짧은 안내 한 줄을 둡니다. 예: `This project also follows Agent Relay (see .agent-relay/PROTOCOL.md).`

### 3. 도구별 지시 파일 처리 (선택)

대상 프로젝트가 사용하는 AI 도구에 따라 다음 파일을 처리합니다. 기본 원칙은 **있으면 머지/추가, 없으면 생성하지 않음**입니다 (`.cursor/`만 예외).

| 트리거 | 동작 | 원본 |
|---|---|---|
| `<project>/CLAUDE.md` 존재 | 기존 내용 보존 + Agent Relay 포인터 머지 | `bootstrap/CLAUDE.md` |
| `<project>/.codex/instructions.md` 존재 | 기존 내용 보존 + Agent Relay 포인터 머지 | `bootstrap/.codex/instructions.md` |
| `<project>/.cursor/` 디렉토리 존재 | `.cursor/rules/agent-relay.mdc` 신규 추가 (충돌 없음) | `bootstrap/.cursor/rules/agent-relay.mdc` |

머지 규칙:

- 기존 내용을 보존하고 Agent Relay 포인터(보통 5~10줄)만 추가합니다.
- 동일한 포인터가 이미 있으면 중복 추가하지 않습니다.
- 해당 파일/디렉토리가 없으면 그 도구를 사용자가 안 쓴다는 뜻이므로 새로 생성하지 않습니다.
- `.cursor/rules/agent-relay.mdc`는 동일 이름의 파일이 이미 있을 때만 머지하고, 그 외에는 새로 추가합니다.

### 4. `.agent-relay/` 디렉토리 복사

`bootstrap/.agent-relay/`를 통째로 `<project>/.agent-relay/`로 복사합니다.

### 5. `relay.log` placeholder 치환

`<project>/.agent-relay/relay.log`의 첫 줄 placeholder timestamp를 현재 시각으로 치환합니다.

- 형식: ISO-8601, 타임존 포함 (예: `2026-04-28T19:45:00+09:00`)
- 예시 변환:
  - 변환 전: `<YYYY-MM-DDTHH:MM:SS+09:00> | SESSION_START | agent=Human | event=PROJECT_BOOTSTRAPPED | summary="Agent Relay initialized"`
  - 변환 후: `2026-04-28T19:45:00+09:00 | SESSION_START | agent=Codex(GPT-5.5) | event=PROJECT_BOOTSTRAPPED | summary="Agent Relay initialized"`
- `agent=Human`은 사용자가 직접 부트스트랩한 경우에만 사용합니다.
- 에이전트가 자동 부트스트랩한 경우 0단계에서 정한 세션 에이전트 이름을 사용합니다. 예: `agent=Codex(GPT-5.5)`.

### 6. `INDEX.md` placeholder 채우기

`<project>/.agent-relay/INDEX.md`의 placeholder를 현재 프로젝트 정보로 채웁니다.

| Placeholder | 채우는 값 |
|---|---|
| `<project-name>` | `<project>` 디렉토리명 또는 프로젝트의 실제 이름 |
| `<repo-name>` | `git config --get remote.origin.url`에서 추출하거나, 없으면 디렉토리명 |
| `Primary Purpose` | 프로젝트의 장기 목적을 `README.md`에서 추출하거나 사용자에게 한 줄 요청 |

### 7. `GUIDANCE.md` 안내

`GUIDANCE.md`는 기본 템플릿으로 복사됩니다. 부트스트랩 직후 억지로 채우지 않습니다.

이 파일은 진행 상태나 현재 작업을 기록하는 곳이 아닙니다. Agent Relay를 사용하는 동안 사용자가 지속적으로 주는 장기 지침, 제약, 선호, 금지사항을 누적하는 곳입니다.

기록 대상:

- 오래 유지되는 프로젝트 맥락
- 사용자 선호와 상시 지시
- 기술·제품·운영 제약
- 하지 말아야 할 것
- 보안·개인정보 규칙
- 코드/테스트/브랜치/작업 관례

### 8. 보안 점검

다음 정보는 `<project>/.agent-relay/` 어디에도 저장하지 않습니다.

- API 키, 토큰, 비밀번호
- 자격 증명, 사적 키
- 고객 데이터, 개인정보
- 민감한 내부 정보, 운영 비밀

부트스트랩 직후 `relay.log`, `INDEX.md`, `GUIDANCE.md`에 위 정보가 들어가지 않았는지 한 번 점검합니다.

### 9. Git 정책 확인

`<project>`가 Git 저장소라면 `.agent-relay/`는 커밋 대상입니다.

Agent Relay는 에이전트 간 작업 인수인계를 위한 프로젝트 자산입니다. `.agent-relay/`를 `.gitignore`에 추가하지 않습니다.

이미 `.gitignore`에 `.agent-relay/`가 있다면 제거해야 합니다. 단, 보안 규칙에 따라 비밀정보는 `.agent-relay/`에 저장하지 않습니다.

---

## 완료 보고 형식

부트스트랩이 끝나면 사용자에게 다음 형식으로 보고합니다.

```text
Agent Relay 부트스트랩 완료.

생성/변경한 파일:
- AGENTS.md (신규 / 머지)
- (선택) CLAUDE.md / .codex/instructions.md 머지
- (선택) .cursor/rules/agent-relay.mdc 신규 추가
- .agent-relay/PROTOCOL.md
- .agent-relay/INDEX.md
- .agent-relay/relay.log (timestamp 치환됨)
- .agent-relay/handoff/
- .agent-relay/templates/work-end.md
- .agent-relay/templates/handoff.md
- .agent-relay/GUIDANCE.md (장기 지침/제약 누적 템플릿)
- (필요 시) .gitignore에서 .agent-relay/ 무시 규칙 제거

다음 단계:
- INDEX.md의 Project / Important Files 확인
- 이후 사용자 장기 지침이 생기면 GUIDANCE.md에 누적
- .agent-relay/가 Git 추적 대상인지 확인
- 첫 의미 있는 작업 후 relay.log에 WORK_END 한 줄 남기기
```

---

## 중단·실패 보고 형식

`.agent-relay/`가 이미 있어 진행을 중단했거나, 머지 중 충돌이 발생한 경우 다음과 같이 보고합니다.

```text
Agent Relay 부트스트랩을 중단했습니다.

이유:
- <project>/.agent-relay/ 가 이미 존재
- 기존 relay.log 마지막 이벤트: <timestamp> | <event_type> | ...

확인이 필요한 사항:
- 재초기화를 원하시면 .agent-relay/를 백업 후 알려주세요.
- 일부 파일만 갱신하려면 어떤 파일을 갱신할지 지정해주세요.
```

---

## 참고

- 사양 전체: [`bootstrap/.agent-relay/PROTOCOL.md`](bootstrap/.agent-relay/PROTOCOL.md)
- 한국어 가이드: [`PROTOCOL-GUIDE.md`](PROTOCOL-GUIDE.md)
