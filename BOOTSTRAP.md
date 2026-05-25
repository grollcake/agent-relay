# Bootstrap

이 문서는 **AI 에이전트가 사용자 프로젝트에 Agent Relay를 적용할 때 따라야 하는 절차**입니다.
사람이 직접 따라할 수도 있지만, 기본은 에이전트가 이 문서를 읽고 자동 수행합니다.

저장소 개요와 사용자용 한 줄 프롬프트는 [`README.md`](README.md)를 참조하세요.

---

## 사전 확인

부트스트랩을 시작하기 전, 사용자 프로젝트(이하 `<project>`)의 다음 상태를 확인합니다.

- `<project>/AGENTS.md` 존재 여부
- `<project>/.agent-relay/` 디렉토리 존재 여부
- `<project>/.agent-relay/VERSION` 존재 여부 (이미 설치된 경우 업데이트 기준)
- `<project>/CLAUDE.md`, `<project>/.codex/instructions.md` 존재 여부 (있으면 머지 대상)
- `<project>/.cursor/` 디렉토리 존재 여부 (있으면 룰 파일 신규 추가 대상)
- `<project>/README.md` 존재 여부 (없어도 진행 가능)

---

## 부트스트랩 절차

### 1. `.agent-relay/` 디렉토리 처리

`<project>/.agent-relay/` 존재 여부에 따라 분기합니다.

- **이미 있으면**: **부트스트랩을 중단**하고 사용자에게 보고합니다.
  - 이미 Agent Relay가 적용된 프로젝트일 가능성이 큽니다.
  - 사용자가 "최신화", "업데이트", "sync"를 요청한 경우에는 새 설치가 아니라 아래 "업데이트 절차"를 따릅니다.
  - 사용자 확인 없이 기존 `AGENTS.md`, `relay.log`, `runs/`, `GUIDANCE.md`, `lesson-learned/`를 덮어쓰거나 머지하지 않습니다.
  - 사용자가 명시적으로 "재초기화" 또는 "특정 파일만 갱신"을 요청한 경우에만, 해당 범위로 한정해 진행합니다.
- **없으면**: 이후 단계로 진행합니다.

### 2. `AGENTS.md` 처리

`<project>/AGENTS.md` 존재 여부에 따라 분기합니다.

- **없으면**: `bootstrap/AGENTS.md`를 그대로 `<project>/AGENTS.md`로 복사합니다.
- **이미 있으면**: 덮어쓰지 말고 **머지**합니다.
  - 기존 내용 전체를 보존합니다.
  - 기존 파일에는 `bootstrap/AGENTS.md`의 `## Agent Relay` 섹션만 추가합니다. 이 섹션의 원본은 `bootstrap/AGENTS.md`이며, `BOOTSTRAP.md`에는 전문을 중복 유지하지 않습니다.
  - 동일한 Agent Relay 섹션이나 포인터가 이미 있으면 중복 추가하지 않고 누락된 문장만 보강합니다.
  - 기존 프로젝트의 규칙·역할·지시문을 임의로 삭제하거나 재작성하지 않습니다.

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

### 5. 초기 `relay.log` 기록

복사된 `<project>/.agent-relay/relay.log`에는 부트스트랩 작업을 기록하기 위한 두 줄이 들어 있습니다. 이 두 줄의 placeholder를 실제 값으로 바꿉니다.

1. 현재 KST 시간을 `YYYY-MM-DDTHH:MM:SS` 형식으로 적습니다.
2. 무작위 소문자 영문 4글자로 `task-id` 하나를 만듭니다.
3. `REQUEST` 줄과 `RUN_DONE` 줄에 같은 `task-id`를 넣습니다.
4. 두 줄의 timestamp를 실제 부트스트랩 시작/완료 시각으로 바꿉니다. 같은 시각을 써도 됩니다.

예시:

```text
2026-04-28T19:45:00 | abcd | REQUEST  | PM | Bootstrap Agent Relay
2026-04-28T19:46:00 | abcd | RUN_DONE | PM | Agent Relay initialized
```

### 6. `GUIDANCE.md`와 `LESSON-LEARNED.md` 안내

`GUIDANCE.md`는 기본 템플릿으로 복사됩니다. 부트스트랩 직후 억지로 채우지 않습니다.

`LESSON-LEARNED.md`는 완료된 작업에서 얻은 재사용 가능한 해결 지식을 `.agent-relay/lesson-learned/`에 기록하는 기준을 설명합니다. 기록은 작업 완료 승인 이후 사용자가 수락한 항목만 추가합니다.

### 7. Git 정책 확인

`<project>`가 Git 저장소라면 `.agent-relay/`는 커밋 대상입니다. 이미 `.gitignore`에 `.agent-relay/`가 있다면 제거해야 합니다.

### 8. 설치 후 프로토콜 확인

부트스트랩을 마친 에이전트는 대상 프로젝트의 `<project>/.agent-relay/PROTOCOL.md`를 읽고, 이후 작업부터 그 프로토콜을 따릅니다.

---

## 완료 보고 형식

부트스트랩이 끝나면 사용자에게 다음 형식으로 보고합니다.

```text
Agent Relay 부트스트랩 완료.

적용 버전: x.x

생성/변경한 파일:
- AGENTS.md (신규 / 머지)
- (선택) CLAUDE.md / .codex/instructions.md 머지
- (선택) .cursor/rules/agent-relay.mdc 신규 추가
- .agent-relay/
- .agent-relay/PROTOCOL.md
- .agent-relay/VERSION
- .agent-relay/GUIDANCE.md
- .agent-relay/LESSON-LEARNED.md
- .agent-relay/relay.log
- .agent-relay/lesson-learned/
- .agent-relay/runs/
- .agent-relay/templates/
- (필요 시) .gitignore에서 .agent-relay/ 무시 규칙 제거

다음 단계:
- .agent-relay/PROTOCOL.md를 읽고 이후 작업부터 그 분류와 이벤트 흐름에 따라 기록
- .agent-relay/가 Git 추적 대상인지 확인
```

---

## 업데이트 절차

사용자가 이미 Agent Relay가 적용된 프로젝트에서 "agent-relay 최신화해줘", "업데이트해줘", "sync 해줘"처럼 요청하면 새 부트스트랩을 하지 않고 다음 절차를 따릅니다.

### 1. 설치 버전 읽기

`<project>/.agent-relay/VERSION`을 읽습니다.

- 없으면 설치 버전을 알 수 없으므로 중단하고 사용자에게 보고합니다.
- 있으면 그 값을 현재 설치 버전으로 보고, 기본 upstream `https://github.com/grollcake/agent-relay`의 `main` 브랜치와 비교합니다.

### 2. upstream 확보와 비교

기본 upstream의 최신 `main`을 임시 위치에 가져와 현재 프로젝트와 비교합니다.

- 현재 프로젝트 파일과 최신 upstream을 직접 비교하되, 프로젝트별 내용은 보수적으로 보존합니다.

### 3. 파일별 업데이트 정책

| 대상 | 업데이트 방식 |
|---|---|
| `AGENTS.md` | 최신 `bootstrap/AGENTS.md`의 `Agent Relay` 섹션과 비교해 현재 파일의 Agent Relay 섹션만 교체 또는 보강 |
| `.agent-relay/PROTOCOL.md` | 로컬 수정이 없거나 안전히 구분되면 최신 upstream으로 갱신 |
| `.agent-relay/templates/` | 템플릿 파일은 최신 upstream과 비교해 갱신 |
| `.agent-relay/VERSION` | 성공 후 최신 upstream의 `VERSION` 값으로 갱신 |
| `CLAUDE.md`, `.codex/instructions.md`, `.cursor/rules/agent-relay.mdc` | 존재하는 경우 Agent Relay 지시만 비교해 보강 |
| `.agent-relay/GUIDANCE.md` | 덮어쓰지 않음 |
| `.agent-relay/LESSON-LEARNED.md` | 해결 지식 기록 안내 문서이므로 로컬 수정이 없거나 안전히 구분될 때만 갱신 |
| `.agent-relay/relay.log` | 기존 줄을 수정하지 않음. 이전 버전의 다른 형식도 보존하고, 새 이벤트부터 최신 `REQUEST`, `PLAN`, `RUN`, `REVIEW`, `DONE`, `RUN_DONE` 형식을 사용 |
| `.agent-relay/runs/` | 덮어쓰지 않음 |
| `.agent-relay/lesson-learned/` | 덮어쓰지 않음 |

### 4. `AGENTS.md` 머지 규칙

`AGENTS.md`는 프로젝트별 지시가 가장 많이 섞이는 파일이므로 다음 원칙을 따릅니다.

- 기존 프로젝트 규칙, 역할, 금지사항, 빌드/테스트 지시를 삭제하지 않습니다.
- `## Agent Relay` 섹션이 있으면 그 섹션만 최신 `bootstrap/AGENTS.md`의 `## Agent Relay` 섹션과 비교합니다.
- Agent Relay 포인터만 있고 섹션이 없으면 최신 섹션을 추가하되 중복 문장은 제거합니다.
- Agent Relay 섹션 안에 프로젝트 고유 지시가 섞여 있어 자동 분리가 어렵다면 파일을 바꾸지 않고 충돌로 보고합니다.

### 5. 완료 보고

업데이트가 끝나면 변경 파일, 보존한 파일, 충돌 또는 수동 확인이 필요한 파일, 업데이트 전후 버전을 보고합니다.

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
