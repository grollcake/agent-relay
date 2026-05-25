# Agent Relay 한국어 가이드

이 문서는 Agent Relay를 프로젝트에 도입하거나 운영할 때 읽는 한국어 해설서입니다.
정식 배포 파일은 `bootstrap/` 아래에 있으며, 이 문서는 그 파일들의 의도와 사용 방법을 설명합니다.

## 1. Agent Relay란

Agent Relay는 **Leader / Planner / Runner 에이전트 팀**이 역할을 나누고, 기록 기반으로 작업을 이어가기 위한 파일 기반 협업 규약입니다. 프로젝트 안에 작업 분류, 이벤트 타임라인, 라운드 산출물, 작업 맥락 전달 표준을 남기는 것입니다.

## 2. 에이전트 팀 구성

저장소 수준 지시가 다른 절차를 지정하지 않는 한, 표준 구현 작업에는 먼저 **Leader / Planner / Runner 에이전트 팀**을 구성하고 Leader / Planner / Runner 프로토콜을 적용합니다. **Leader (허브)**는 사용자 소통, 분류, 범위/위험 결정, 위임, 결과 해석, 최종 보고를 담당합니다. Planner와 Runner는 Leader를 통해서만 통신합니다.

| 역할 | 책임 |
| --- | --- |
| **Leader (허브)** | 작업을 라우팅하고 증거가 요청을 충족하는지 판단합니다. 단순 중계자가 아닙니다. |
| **Planner** | `PLAN`을 작성하고, 구현이 계획과 일치하는지 검토합니다. 발견은 `blocker`(반드시 수정) 또는 `nit`(비차단)으로 표시합니다. |
| **Runner** | `PLAN`을 구현하고 검증합니다. 모호함은 범위를 넓히지 않은 채 Leader에게 되돌립니다. |

Planner와 Runner는 **Leader를 통해서만** 통신합니다. 사용 도구의 능력에 따라 배정된 멤버를 병렬 또는 순차로 실행할 수 있지만, 기록 없이 단일 에이전트 작업으로 축소해서는 안 됩니다.

**강제 선행 규칙:** Leader, Planner, Runner는 기록이 필요한 작업에 착수하기 전에 반드시 `.agent-relay/GUIDANCE.md`, `.agent-relay/LESSON-LEARNED.md`, `.agent-relay/lesson-learned/`를 읽어 현재 적용할 지침과 이전 해결 지식을 확인합니다. 명백한 기록 제외 요청은 이 확인 없이 응답할 수 있지만, 파일 변경·조사·설계 판단·프로젝트 지침 의존 답변으로 넘어가면 먼저 이 확인을 완료해야 합니다.

## 3. 작업 분류

Leader는 먼저 요청이 명백한 기록 제외 대상인지 가볍게 판단합니다. 기록이 필요한 요청이면 필수 지침·교훈 확인을 마친 뒤 `Trivial` 또는 `Standard`로 분류합니다. Agent Relay **부트스트랩**과 **업데이트**(`.agent-relay/`·Agent Relay 지시 파일 동기화)는 기록 제외가 아니며, Leader가 직접 수행하면 `Trivial`로 `REQUEST → RUN_DONE`을 기록합니다.

| 분류 | 일반적 범위 | 처리 방식 |
| --- | --- | --- |
| 기록 제외 | 단순 질문 답변, 짧은 설명, 브레인스토밍 | 응답만 하고 이벤트를 남기지 않음 |
| `Trivial` | 사소한 텍스트/설정 변경, 명백한 국소 편집, Agent Relay 부트스트랩·업데이트 | Leader가 직접 처리하고 `REQUEST → RUN_DONE` 기록 |
| `Standard` | 다중 파일 구현, 설계 판단, 구현 검증이 필요한 작업 | Planner → Runner → Planner 검토 |

## 4. 백그라운드 위임

`Standard` 작업은 가능한 한 백그라운드로 위임합니다. Leader는 위임 후에도 사용자 응답을 계속 담당합니다. 완료를 기다리기 위한 폴링이나 sleep은 하지 않습니다. 백그라운드 실행이 어렵다면 같은 단계와 산출물 전달 방식을 순차적으로 따릅니다.

## 5. 이벤트 타임라인

`relay.log`는 추가 전용 이벤트 타임라인입니다. 기존 줄을 수정하지 않습니다.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

- `timestamp`는 KST 기준 `YYYY-MM-DDTHH:MM:SS` 형식으로 기록합니다.
- `task-id`는 무작위 소문자 영문 4글자를 씁니다.
- Leader는 `REQUEST` 기록 시 `task-id` 하나를 정하고, 같은 Standard 작업의 `PLAN`/`RUN_ST`/`RUN_ED`/`REVIEW`/`FEEDBACK`/`DONE`까지 재사용합니다. 새 `REQUEST`마다 새 `task-id`를 씁니다.
- 이벤트는 `REQUEST`, `FEEDBACK`, `PLAN`, `RUN_ST`, `RUN_ED`, `REVIEW`, `DONE`, `RUN_DONE`만 씁니다.
- Leader 직접 처리 흐름은 `REQUEST → RUN_DONE`입니다.
- 표준 처리 흐름은 `REQUEST` → `PLAN` → `RUN_ST` → `RUN_ED` → `REVIEW` → `DONE`입니다.
- `FEEDBACK`은 사용자가 `DONE` 승인 전 피드백·결함을 알려줄 때 Leader가 기록합니다. 같은 `task-id`와 산출물 파일 키를 유지합니다.
- `FEEDBACK` 후 Leader는 **현재 런에 추가**할지 **새로운 런**으로 돌릴지 사용자에게 묻습니다. 명백한 결함이면 사용자 확인 없이 현재 런에 추가합니다.
- **현재 런에 추가**: 마지막 `RUN-<NN>` 범위와 기존 `PLAN` 안에서 `RUN_ST` → `RUN_ED` → `REVIEW-<NN>`을 다시 진행합니다. `DONE` 이벤트 승인 전이면 `RUN-<NN>.md` 갱신을 허용합니다.
- **새로운 런**: 다음 `RUN-<NN+1>`로 `RUN_ST` → `RUN_ED` → `REVIEW-<NN+1>`을 진행합니다.
- `role` 주변 공백은 정렬용이며 의미가 없습니다.
- `event`와 `role`은 각각 8자 폭으로 왼쪽 정렬하고 부족한 자리는 공백으로 채웁니다.
- `RUN_ST`는 Leader가 Runner 위임 시 기록합니다. `path` 없음. `summary`에 `RUN-<NN>` 번호를 포함합니다.
- `RUN_ED`는 Runner가 `RUN-<NN>.md` 작성 후 기록합니다. `path` 필수.
- 한 라운드 `<NN>`은 `RUN_ST` → `RUN_ED` 한 쌍입니다. `RUN_ST`는 `path`가 없으므로 바로 다음 `RUN_ED`와 짝입니다. `blocker`로 다음 라운드를 돌릴 때 같은 `task-id`에 `RUN_ST`/`RUN_ED`를 다시 추가합니다.
- 긴 설명은 `relay.log`에 직접 넣지 말고 `.agent-relay/runs/`의 라운드 산출물로 분리합니다.

예시:

```text
2026-05-25T20:40:00 | qmxz | REQUEST  | Leader  | Fix typo in README
2026-05-25T20:41:00 | qmxz | RUN_DONE | Leader  | Fixed typo directly
2026-05-25T20:50:00 | abcd | REQUEST  | Leader  | Update protocol docs
2026-05-25T20:55:00 | abcd | PLAN     | Planner  | Plan written | .agent-relay/runs/20260525-2055-docs-PLAN.md
2026-05-25T20:56:00 | abcd | RUN_ST   | Leader  | RUN-01 delegated
2026-05-25T21:10:00 | abcd | RUN_ED   | Runner  | Changes submitted | .agent-relay/runs/20260525-2055-docs-RUN-01.md
2026-05-25T21:15:00 | abcd | REVIEW   | Planner  | Accepted | .agent-relay/runs/20260525-2055-docs-REVIEW-01.md
2026-05-25T21:16:00 | abcd | DONE     | Leader  | Completed
```

사용자가 `DONE` 승인 전 결함을 알려준 경우 — 명백한 결함, 현재 런에 추가(같은 `task-id`):

```text
2026-05-26T10:20:00 | abcd | FEEDBACK | Leader  | User reported missing validation
2026-05-26T10:21:00 | abcd | RUN_ST   | Leader  | RUN-01 retry
2026-05-26T10:35:00 | abcd | RUN_ED   | Runner  | Fix submitted | .agent-relay/runs/20260525-2055-docs-RUN-01.md
2026-05-26T10:40:00 | abcd | REVIEW   | Planner  | Accepted | .agent-relay/runs/20260525-2055-docs-REVIEW-01.md
```

피드백 후 새로운 런을 선택한 경우:

```text
2026-05-26T11:00:00 | abcd | FEEDBACK | Leader  | User requested scope change
2026-05-26T11:05:00 | abcd | RUN_ST   | Leader  | RUN-02 delegated
2026-05-26T11:20:00 | abcd | RUN_ED   | Runner  | Changes submitted | .agent-relay/runs/20260525-2055-docs-RUN-02.md
2026-05-26T11:25:00 | abcd | REVIEW   | Planner  | Accepted | .agent-relay/runs/20260525-2055-docs-REVIEW-02.md
```

`blocker`로 RUN-02가 필요한 경우(같은 `task-id`):

```text
2026-05-25T21:20:00 | abcd | RUN_ST   | Leader  | RUN-02 delegated
2026-05-25T21:35:00 | abcd | RUN_ED   | Runner  | Changes submitted | .agent-relay/runs/20260525-2055-docs-RUN-02.md
2026-05-25T21:40:00 | abcd | REVIEW   | Planner  | Accepted | .agent-relay/runs/20260525-2055-docs-REVIEW-02.md
```

## 6. 산출물

모든 라운드 산출물은 `.agent-relay/runs/`에 작성하며, 작업당 하나의 안정된 `<YYYYMMDD>-<HHMM>-<slug>` 키를 씁니다. **이전 라운드를 덮어쓰지 않습니다.**

| 산출물 | 경로 | 작성자 | 의미 |
| --- | --- | --- | --- |
| Plan | `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-PLAN.md` | Planner | 계획과 성공 기준 |
| Submission | `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-RUN-<NN>.md` | Runner | 라운드 `<NN>`의 변경/검증/리스크 |
| Review | `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-REVIEW-<NN>.md` | Planner | 같은 라운드 발견 |
| Acceptance | `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-DONE.md` | Planner | `blocker` 없이 검토를 통과한 수락 결과 |

- `<NN>`은 `01`부터 시작합니다.
- 이전 라운드를 덮어쓰지 않습니다. 예외: `FEEDBACK` 후 현재 런에 추가할 때, `DONE` 이벤트 승인 전이면 `RUN-<NN>.md` 갱신을 허용합니다.
- `<SLUG>`는 Leader가 정한 소문자 kebab-case 작업 키를 씁니다.
- `<YYYYMMDD>`와 `<HHMM>`은 Leader가 `REQUEST`를 기록할 때의 KST 날짜·시분(24시간, 구분자 없음)을 씁니다. 예: `20260526-1430-diary-write`.
- `task-id`는 `relay.log` 이벤트 식별자이고, `<YYYYMMDD>-<HHMM>-<SLUG>`는 `.agent-relay/runs/` 산출물 파일 키입니다. 같은 Standard 작업에서는 `task-id` 하나와 파일 키 하나를 함께 씁니다.
- 같은 작업의 모든 라운드 산출물은 같은 `<YYYYMMDD>-<HHMM>-<SLUG>` 키를 씁니다.
- 산출물은 `.agent-relay/templates/plan.md`, `run.md`, `review.md`, `done.md` 형식을 따릅니다.
- Runner는 절대 `DONE`을 쓰지 않습니다.
- Planner는 검토에 `blocker`가 없을 때 `DONE` 산출물을 작성합니다. `nit`는 `DONE`에 기록할 수 있습니다.
- **사용자가 명시적으로 승인하기 전에는 Leader가 `DONE` 이벤트를 기록하여 작업을 종료할 수 없습니다.**
- 각 `RUN`은 변경 파일, 변경 요약, 테스트/검증, 미해결 리스크를 기록합니다.
- `relay.log`는 `REQUEST`, `FEEDBACK`, `PLAN`, `RUN_ST`, `RUN_ED`, `REVIEW`, `DONE`, `RUN_DONE` 이벤트를 추가-전용으로 남기고 `path`로 산출물을 가리킵니다.
- 산출물 작성과 `relay.log` 이벤트 추가는 별개의 필수 작업입니다. 각 단계는 산출물 작성과 해당 이벤트 추가가 모두 끝나야 완료된 것으로 봅니다.
- 산출물 작성자가 해당 이벤트를 추가합니다. Planner는 `PLAN`/`REVIEW`, Runner는 `RUN_ED`, Leader는 `REQUEST`/`FEEDBACK`/`RUN_ST`/`RUN_DONE`과 사용자 승인 후 최종 `DONE`을 기록합니다.
- Leader는 다음 단계 위임 전에 직전 단계 이벤트가 `relay.log`에 추가됐는지 확인합니다.

## 7. 파이프라인

1. Leader가 요청을 분류합니다.
2. 기록 제외 대상이면 응답만 하고 이벤트를 남기지 않습니다.
3. `Trivial`이면 Leader가 직접 처리하고 `REQUEST → RUN_DONE` 이벤트 흐름으로 작업을 닫습니다.
4. `Standard`이면 Planner가 `PLAN`을 작성합니다.
5. Leader가 Runner에게 위임하고 `RUN_ST`를 기록합니다.
6. Runner는 `PLAN`·성공 기준·범위에 따라 구현한 뒤 `RUN-01`을 쓰고 `RUN_ED`를 기록합니다.
7. Planner는 해당 `RUN` 경로를 받아 같은 번호의 `REVIEW`를 씁니다.
8. `blocker`가 없으면 Planner가 `DONE` 산출물을 씁니다. Leader는 결과·nit·리스크와 `DONE` 산출물 경로를 사용자에게 보고하고 승인을 요청합니다.
9. 사용자가 명시적으로 승인한 뒤에만 Leader가 `DONE` 이벤트를 기록하여 작업을 닫습니다.
10. 사용자가 승인 대신 피드백·결함을 알려주면 Leader가 `FEEDBACK`을 기록합니다. 명백한 결함이면 현재 런에 추가하고, 그렇지 않으면 **현재 런에 추가 / 새로운 런** 중 사용자 선택을 받습니다. 이후 5번(`RUN_ST`)부터 다시 진행합니다.
11. `DONE` 승인을 받은 Leader는 해당 세션에서 발생한 착오, 해결 방법, 사용자 의견을 종합하여 `.agent-relay/GUIDANCE.md` 수정안 또는 `.agent-relay/lesson-learned/` 추가안을 사용자에게 제안합니다. 사용자가 수락한 항목만 기록합니다.
12. `blocker`가 있으면 Leader가 `RUN_ST`로 다음 라운드를 위임하고, Runner가 `RUN-<NN>` → `RUN_ED`를, Planner가 다음 `REVIEW`를 씁니다. `REVIEW-03` 전까지 사용자 승인 없이 진행합니다.
13. `REVIEW-03`까지도 `blocker`가 남으면 Leader는 상태를 보고하고 사용자에게 **재시도 / 계획 수정 / 부분 수락 / 중단** 중 선택을 요청합니다.

Standard 작업에서 사용자 개입이 필요한 경우는 `DONE` 최종 승인, `DONE` 승인 전 피드백·결함(`FEEDBACK`)과 `FEEDBACK` 후 현재 런·새 런 선택, `REVIEW-03` 이후에도 `blocker`가 남는 경우, Leader가 사용자 결정이 필요하다고 판단한 경우뿐입니다. `FEEDBACK` 이후 `RUN`/`REVIEW` 재진행과 중간 라운드 반복은 Leader가 진행합니다.

`Trivial` 작업은 사용자 완료 승인 없이 `RUN_DONE`으로 닫을 수 있습니다. 다만 장기 지침이나 재사용 가능한 교훈이 생겼다면 Leader는 사용자에게 기록안을 제안하고, 사용자가 수락한 항목만 `GUIDANCE.md` 또는 `lesson-learned/`에 추가합니다.

## 8. 위임과 보고

### 컨텍스트 관리

Leader, Planner, Runner는 한 작업 안에서 컨텍스트 교체 없이 연속 사용하는 것을 전제로 합니다.

Leader는 컨텍스트가 불필요하게 커지지 않도록 위임 결과를 받을 때 산출물 **경로 + 최소 결정 정보**만 보관합니다.

- 한 줄 결과
- 한 줄 검증 상태
- 해당 시 `blocker` 건수/요약
- 잔존 리스크 또는 사용자 결정 요구

결정적 모호함이나 사용자 결정이 필요할 때를 제외하고 전체 산출물을 Leader 컨텍스트에 적재하지 않습니다.

Planner/Runner는 가능하면 같은 컨텍스트를 유지하되, 다음 중 하나라도 발생하면 Leader가 사용자에게 **교체 여부**를 물어야 합니다.

- 한 인스턴스에 후속 메시지가 5개 이상 누적
- 작업 주제가 명백히 바뀜
- 응답이 눈에 띄게 느려지거나 이전 컨텍스트를 혼동

### 위임 시 필수 필드

Planner/Runner에게 보내는 모든 프롬프트는 자기완결적이어야 하며 다음을 포함합니다.

- 목표(goal)
- 관련 파일 또는 조사 범위
- 산출물 타입과 정확한 경로
- 성공 기준과 검증 방법
- 범위 외 작업 금지 명시
- 불명확한 사항은 추정하지 말고 Leader에게 되돌릴 것
- 단계에 필요한 입력 산출물 경로

### 보고 시 필수 필드

Planner가 Leader에게 보고할 때는 다음만 간결히 포함합니다.

- `PLAN`/`REVIEW`/`DONE` 산출물 경로
- 판단 결과
- `blocker` 수와 요약
- `nit` 요약
- 사용자 결정 필요 여부

Runner가 Leader에게 보고할 때는 다음만 간결히 포함합니다.

- `RUN` 산출물 경로
- 변경 요약
- 검증 결과
- 미해결 리스크
- 범위 밖으로 넘긴 사항

## 9. 목표 파일 구조

Agent Relay를 사용자 프로젝트에 도입했을 때의 `runs/` 중심 파일 구조입니다.

```text
project-root/
├── AGENTS.md
├── CLAUDE.md                  # 선택, Claude Code에서 실행 중이면 생성
├── .codex/
│   └── instructions.md         # 선택
├── .cursor/
│   └── rules/
│       └── agent-relay.mdc     # 선택
└── .agent-relay/
    ├── PROTOCOL.md
    ├── VERSION
    ├── GUIDANCE.md             # 장기 지침/제약 누적
    ├── LESSON-LEARNED.md       # 완료 작업에서 얻은 해결 지식 기록 안내
    ├── relay.log
    ├── lesson-learned/         # 완료 작업에서 얻은 해결 지식 기록
    │   └── .gitkeep
    ├── runs/                   # 라운드 산출물(PLAN/RUN/REVIEW/DONE)
    │   └── .gitkeep
    └── templates/
        ├── guidance.md
        ├── lesson-learned.md
        ├── plan.md
        ├── run.md
        ├── review.md
        └── done.md
```

## 10. 파일별 역할

| 파일 | 필수 여부 | 역할 |
|---|---:|---|
| `AGENTS.md` | 필수 | 모든 에이전트가 읽는 공통 작업 지침입니다. |
| `.agent-relay/PROTOCOL.md` | 필수 | Agent Relay의 최소 규칙입니다. |
| `.agent-relay/VERSION` | 필수 | 설치 버전입니다. 업데이트 시 기본 upstream과 비교하는 기준으로 씁니다. |
| `.agent-relay/GUIDANCE.md` | 누적 관리 | 세션을 넘어 유지할 사용자 지침, 제약, 금지사항을 담는 문서입니다. |
| `.agent-relay/LESSON-LEARNED.md` | 필수 | 완료 작업에서 얻은 해결 지식 기록의 목적과 작성 방식을 설명합니다. |
| `.agent-relay/relay.log` | 필수 | 작업 이벤트 로그입니다 (추가 전용) |
| `.agent-relay/lesson-learned/` | 누적 관리 | 완료된 작업에서 얻은 재사용 가능한 해결 지식이 쌓이는 디렉토리입니다. |
| `.agent-relay/runs/` | 목표 필수 | `Standard` 작업의 라운드 산출물이 쌓이는 디렉토리입니다. |
| `.agent-relay/templates/` | 목표 권장 | 산출물·누적 문서 작성 형식 예시입니다. |

## 11. 합류할 때 읽는 순서

새 세션을 시작하면 AI 에이전트는 `AGENTS.md`를 우선 읽습니다. `AGENTS.md`에서 Agent Relay 안내를 확인하면 `.agent-relay/PROTOCOL.md`를 읽고, 그 규칙에 따라 프로젝트 맥락과 진행 중인 작업을 확인합니다.

Agent Relay에 합류할 때의 읽기 순서는 다음과 같습니다.

```text
1. AGENTS.md에서 Agent Relay 안내 확인
2. .agent-relay/PROTOCOL.md
3. .agent-relay/GUIDANCE.md
4. .agent-relay/LESSON-LEARNED.md
5. .agent-relay/lesson-learned/의 기존 기록
6. .agent-relay/relay.log의 마지막 50줄
7. 진행 중인 라운드가 있으면 .agent-relay/runs/의 최신 PLAN/RUN/REVIEW 읽기
```

같은 세션에서 연속 작업 중이라면 매 사용자 메시지마다 `relay.log`를 다시 읽지 않습니다. 다만 기록이 필요한 새 요청에 착수할 때는 Leader, Planner, Runner 모두 `GUIDANCE.md`, `LESSON-LEARNED.md`, `lesson-learned/`를 반드시 다시 읽습니다.

## 12. 기록해야 하는 작업

다음 작업은 `relay.log`에 결과를 남기는 것이 좋습니다.

- 코드 변경
- 문서 변경
- 설정 변경
- 테스트 실행
- 디버깅 시도
- 아키텍처 또는 설계 판단
- 이슈 상태 변경
- 다음 역할이나 후속 세션이 알아야 할 맥락 생성

`Standard`로 분류된 작업은 라운드 산출물(`PLAN`/`RUN`/`REVIEW`/`DONE`)이 본문 기록이고, `relay.log`는 그 산출물을 가리키는 이벤트 인덱스 역할을 합니다.

단순 질의응답, 짧은 설명, 브레인스토밍은 보통 기록하지 않습니다.
사용자가 맥락 보존을 명시적으로 요청한 경우에는 예외로 기록할 수 있습니다.

## 13. GUIDANCE 사용 기준

`GUIDANCE.md`는 기본 템플릿으로 복사됩니다.
부트스트랩 직후 억지로 채우는 문서가 아닙니다.

이 파일은 진행 상태, 현재 작업, 임시 계획, 다음 작업을 기록하는 곳이 아닙니다. 그런 정보는 `relay.log`와 라운드 산출물에 둡니다.

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

## 14. 부트스트랩 절차

대상 프로젝트에 Agent Relay를 적용할 때의 실제 절차는 `BOOTSTRAP.md`를
기준으로 진행합니다. 아래는 `README.md`의 `runs/` 중심 목표 구조를
기준으로 한 핵심 절차입니다.

1. 부트스트랩 이후 작업을 맡을 `Leader`, `Planner`, `Runner` 팀을 구성합니다.
2. `.agent-relay/`가 이미 있으면 아무 파일도 바꾸지 않고 중단 후 보고합니다.
3. 대상 프로젝트의 `AGENTS.md`, 도구별 지시 파일 존재 여부를 확인합니다.
4. `AGENTS.md`가 없으면 `bootstrap/AGENTS.md`를 복사합니다.
5. `AGENTS.md`가 이미 있으면 기존 내용을 보존하고 `bootstrap/AGENTS.md`의 `Agent Relay` 섹션만 병합합니다.
6. 도구별 지시 파일을 필요에 따라 병합합니다.
7. `bootstrap/.agent-relay/`를 복사합니다. `lesson-learned/`와 `runs/`는 빈 디렉토리(`.gitkeep`)로 생성합니다.
8. `relay.log`의 자리표시자 timestamp와 이벤트 줄을 현재 작업 정보에 맞게 바꿉니다.
9. `.agent-relay/VERSION`에 설치 버전을 기록합니다.
10. `GUIDANCE.md`와 `LESSON-LEARNED.md`의 용도를 안내합니다.
11. 비밀정보가 `.agent-relay/`에 들어가지 않았는지 확인합니다.
12. Git 저장소라면 `.agent-relay/`가 커밋 대상인지 확인합니다.

## 15. 업데이트 절차

사용자가 "agent-relay 최신화해줘", "업데이트해줘", "sync 해줘"처럼 요청하면 새 부트스트랩이 아니라 업데이트 절차를 수행합니다.

1. `.agent-relay/VERSION`을 읽어 현재 설치 버전을 확인합니다.
2. 기본 upstream `https://github.com/grollcake/agent-relay`의 최신 `main`을 임시 위치에 가져와 현재 프로젝트와 비교합니다.
3. `AGENTS.md`는 최신 `bootstrap/AGENTS.md`의 `Agent Relay` 섹션과 비교해 현재 파일의 Agent Relay 섹션만 교체하거나 보강합니다.
4. `.agent-relay/PROTOCOL.md`와 `.agent-relay/templates/`는 로컬 수정이 없거나 안전히 구분될 때 최신 upstream으로 갱신합니다.
5. `.agent-relay/GUIDANCE.md`, `.agent-relay/lesson-learned/`, `.agent-relay/relay.log`, `.agent-relay/runs/`는 덮어쓰지 않습니다.
6. `.agent-relay/LESSON-LEARNED.md`는 안내 문서이므로 로컬 수정이 없거나 안전히 구분될 때만 갱신합니다.
7. 업데이트가 성공하면 `.agent-relay/VERSION`을 최신 upstream의 `VERSION` 값으로 갱신합니다.
8. Leader는 `relay.log`에 `REQUEST → RUN_DONE`을 추가합니다. 메타 작업이라 기록을 생략하지 않습니다. 보통 `Trivial`이며, `summary`에 이전·이후 `VERSION`을 포함합니다. 범위가 `Standard`에 해당할 때만 `REQUEST → PLAN → RUN_ST → RUN_ED → REVIEW → DONE`을 씁니다.

이전 버전의 `relay.log`가 `agent=`, `task=`, `TASK_BEGIN` 같은 형식을 사용하더라도 기존 줄은 수정하지 않습니다. 새 버전 적용 후 추가하는 이벤트부터 새 형식을 사용합니다.

`AGENTS.md`에 Agent Relay 지시와 프로젝트 고유 지시가 섞여 있어 자동 분리가 어렵다면 파일을 바꾸지 않고 충돌로 보고합니다.

## 16. 도구별 지시 파일

도구별 파일은 선택입니다. 현재 부트스트랩을 실행 중인 에이전트 도구이거나 사용자가 해당 도구를 쓰는 흔적이 있을 때만 추가하거나 병합합니다.

| 조건 | 처리 |
|---|---|
| Claude Code에서 실행 중이거나 `CLAUDE.md`가 이미 있음 | 없으면 생성하고, 있으면 기존 내용을 보존한 채 `bootstrap/CLAUDE.md`의 Agent Relay 포인터를 병합합니다. |
| `.codex/instructions.md`가 이미 있음 | 기존 내용을 보존하고 `bootstrap/.codex/instructions.md`의 포인터를 병합합니다. |
| `.cursor/` 디렉토리가 있음 | `.cursor/rules/agent-relay.mdc`를 추가합니다. 같은 파일이 있으면 병합합니다. |

해당 파일이나 디렉토리가 없으면 새로 만들지 않는 것이 기본입니다.
예외는 `.cursor/` 디렉토리가 이미 있는 경우의 Cursor rule 파일입니다.

## 17. 보안 규칙

`.agent-relay/`에는 다음 정보를 저장하지 않습니다.

- API 키
- 토큰
- 비밀번호
- 개인 키와 자격 증명
- 고객 데이터와 개인정보
- 민감한 내부 정보
- 운영 비밀

이 규칙은 `relay.log`, `GUIDANCE.md`, 라운드 산출물 모두에 적용됩니다.

## 18. 기타 규칙: Git 정책

Git 저장소에서는 `.agent-relay/`를 커밋합니다.

Agent Relay는 역할 및 세션 간 작업 인수인계를 위한 프로젝트 자산입니다. 릴레이 규칙, 프로젝트 지도, 작업 로그, 라운드 산출물은 다음 역할 또는 후속 세션이 같은 저장소에서 이어받을 수 있어야 합니다.

따라서 `.agent-relay/`를 `.gitignore`에 추가하지 않습니다. 이미 무시 규칙이 있다면 제거합니다.

비밀정보는 커밋 여부와 무관하게 `.agent-relay/`에 저장하지 않습니다.

## 19. 사용자에게 줄 수 있는 짧은 호출문

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

## 20. 배포 파일과 목표 템플릿 위치

현재 존재하는 배포 파일과 `runs/` 중심 목표 템플릿은 다음과 같습니다.

| 파일 | 용도 |
|---|---|
| `BOOTSTRAP.md` | 에이전트가 대상 프로젝트에 Agent Relay를 설치할 때 따르는 절차 |
| `bootstrap/AGENTS.md` | 대상 프로젝트 루트에 둘 공통 지시문 |
| `bootstrap/.agent-relay/PROTOCOL.md` | Agent Relay 정식 최소 규칙 |
| `bootstrap/.agent-relay/VERSION` | 설치 버전 템플릿 |
| `bootstrap/.agent-relay/GUIDANCE.md` | 장기 지침/제약 누적 템플릿 |
| `bootstrap/.agent-relay/LESSON-LEARNED.md` | 완료 작업에서 얻은 해결 지식 기록 안내 |
| `bootstrap/.agent-relay/relay.log` | 초기 로그 템플릿 |
| `bootstrap/.agent-relay/lesson-learned/` | 완료 작업에서 얻은 해결 지식 기록 디렉토리 |
| `bootstrap/.agent-relay/templates/guidance.md` | `GUIDANCE.md` 누적 지침 형식 |
| `bootstrap/.agent-relay/templates/lesson-learned.md` | 해결 지식 기록 형식 |
| `bootstrap/.agent-relay/templates/plan.md` | `PLAN` 산출물 형식 |
| `bootstrap/.agent-relay/templates/run.md` | `RUN-NN` 산출물 형식 |
| `bootstrap/.agent-relay/templates/review.md` | `REVIEW-NN` 산출물 형식 |
| `bootstrap/.agent-relay/templates/done.md` | `DONE` 산출물 형식 |

## 21. 정리

Agent Relay는 Leader / Planner / Runner 에이전트 팀이 역할을 나누고, 이벤트 타임라인과 산출물을 기반으로 작업을 이어가기 위한 파일 기반 협업 규약입니다.

핵심 판단 기준은 세 개입니다.

```text
1) 이 요청은 기록 제외인가, Trivial인가, Standard인가?
2) Standard라면 PLAN → RUN_ST → RUN_ED → REVIEW → DONE 라운드를 지키고 있는가?
3) 다음 역할 또는 후속 세션이 이벤트 타임라인과 산출물로 이어받을 수 있는가?
```

기록 제외 대상이면 응답만 하고 이벤트를 남기지 않습니다.
`Trivial`이면 Leader가 직접 처리하고 `REQUEST → RUN_DONE` 흐름이면 충분합니다.
`Standard`이면 PLAN → RUN_ST → RUN_ED → REVIEW를 거치고, `blocker`가 없으면 Planner가 `DONE` 산출물을 작성합니다. Leader는 사용자가 명시적으로 승인한 경우에만 `DONE` 이벤트를 기록하여 작업을 닫습니다.
