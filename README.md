# Agent Relay

> **🤖 AI 에이전트에게**: 당신이 AI 에이전트이고 사용자로부터 이 저장소를 본인 프로젝트에 반영해달라는 요청을 받았다면, [`BOOTSTRAP.md`](BOOTSTRAP.md)를 그대로 따라 부트스트랩을 수행하세요. 기존 `AGENTS.md`는 덮어쓰지 말고 `<agent-relay-rules>...</agent-relay-rules>` 블록만 추가하세요. 전체 `bootstrap/AGENTS.md`는 대상 프로젝트에 `AGENTS.md`가 없을 때만 복사합니다. 이미 `.agent-relay/`가 있으면 새 설치가 아니라 업데이트 요청인지 확인하고, 업데이트 요청이면 `.agent-relay/VERSION`을 읽은 뒤 기본 upstream `https://github.com/grollcake/agent-relay`와 비교하세요. `bootstrap/` 트리는 목적지 경로를 그대로 미러링한 복사 원본입니다.

Agent Relay는 **LeadAI / PlanAI / ExecAI 에이전트 팀**이 역할을 나누고, 기록 기반으로 작업을 이어가기 위한 파일 기반 협업 규약입니다.

이 저장소는 **사양과 템플릿을 모아둔 가이드 저장소**입니다. 사용자 프로젝트에 Agent Relay를 도입하면, 그 프로젝트에 `.agent-relay/`가 생깁니다.

---

## 목적

1. LeadAI / PlanAI / ExecAI 에이전트 팀의 역할 분담으로 속도와 완전성을 높인다.
2. LeadAI는 PlanAI/ExecAI에게 백그라운드로 위임한 뒤 즉시 짧은 상태를 사용자에게 반환하고, 항상 새 사용자 요구에 대응 가능한 상태를 유지한다.
3. 작업 이력을 구조화된 이벤트 타임라인으로 남긴다.
4. 이벤트 타임라인과 산출물을 기반으로 작업 맥락을 전달하고 검토하여 갑자기 중단된 경우에도 연속 작업이 가능하게 한다.
5. 코딩 에이전트가 바뀐 경우(예: Claude Code → Codex)에도 작업 연속성을 유지한다.
6. 세션이 바뀌어도 반복 적용할 사용자 지침과 프로젝트 제약을 누적한다.
7. 완료된 작업에서 얻은 착오, 해결 방법, 재사용 가능한 확인 절차를 다음 작업에 활용한다.

---

## 이 저장소 사용법

### AI 에이전트에게 한 줄로 시키기

대상 프로젝트 디렉토리에서 본인 에이전트(Claude Code 등)에게 다음과 같이 요청하세요.

```text
github.com/grollcake/agent-relay 를 내 프로젝트에 반영해줘
```

### 업데이트하기

이미 Agent Relay가 적용된 프로젝트에서는 다음처럼 요청할 수 있습니다.

```text
agent-relay 최신화해줘
```

에이전트는 대상 프로젝트의 `.agent-relay/VERSION`을 읽고 기본 upstream `https://github.com/grollcake/agent-relay`와 비교합니다. `AGENTS.md`는 프로젝트 고유 지침을 보존하고 `<agent-relay-rules>...</agent-relay-rules>` 블록만 최신 `bootstrap/AGENTS.md`와 비교해 갱신합니다. `LESSON-LEARNED.md`, `GUIDANCE.md`, `lesson-learned/`, `relay.log`, `runs/`는 프로젝트별 상태이므로 덮어쓰지 않습니다. 업데이트 완료 후 LeadAI는 `relay.log`에 `REQUEST → RUN_DONE`을 추가합니다. `summary`에 이전·이후 `VERSION`을 포함합니다.

---

## 목표 부트스트랩 트리

아래 구조는 `runs/` 중심 흐름에 맞춘 현재 `bootstrap/` 배포 템플릿 구조입니다.

```text
bootstrap/
├── AGENTS.md                          # 프로젝트 루트로 복사 (기본 공통 지시문)
├── CLAUDE.md                          # Claude Code용 동일 Agent Relay 블록 (선택)
└── .agent-relay/                      # 프로젝트 루트로 복사 (필수)
    ├── PROTOCOL.md                    # Agent Relay 규칙(필수)
    ├── VERSION                        # 설치 버전
    ├── GUIDANCE.md                    # 장기 지침/제약 누적
    ├── LESSON-LEARNED.md              # lesson-learned/ 실제 기록 인덱스
    ├── relay.log                      # 추가 전용 이벤트 로그
    ├── protocol-guard                       # relay.log 이벤트 추가와 단계 전이 검증 CLI
    ├── lesson-learned/                # 완료 작업에서 얻은 해결 지식 기록
    │   └── .gitkeep
    ├── runs/                          # 라운드 산출물(PLAN/RUN/REVIEW/CLOSE)
    │   └── .gitkeep
    └── templates/                     # 산출물·누적 문서 양식
        ├── guidance.md
        ├── lesson-learned.md
        ├── plan.md
        ├── run.md
        ├── review.md
        └── close.md
```

`bootstrap/` 안의 모든 경로는 **목적지 경로 그대로**입니다. `CLAUDE.md`의 `<agent-relay-rules>...</agent-relay-rules>` 블록은 `AGENTS.md`의 동명 블록과 동일하게 유지하므로 Claude 사용자는 `AGENTS.md`를 제거해도 같은 규칙을 유지할 수 있습니다. `CLAUDE.md`는 Claude Code에서 실행 중이거나 이미 존재할 때만 생성 또는 머지됩니다.

---

## 에이전트 팀 구성

**LeadAI / PlanAI / ExecAI 에이전트 팀**을 구성하고 이 프로토콜을 적용한다. **LeadAI (허브)**는 사용자 소통, 분류, 범위/위험 결정, 위임, 결과 해석, 최종 보고를 담당한다. LeadAI는 PlanAI/ExecAI에게 백그라운드로 작업을 위임한 뒤 즉시 짧은 상태를 사용자에게 반환하고, 위임 중에도 새 사용자 요구에 대응 가능한 상태를 유지한다. PlanAI와 ExecAI는 LeadAI를 통해서만 통신한다.

| 역할 | 책임 |
| --- | --- |
| **LeadAI (허브)** | 작업을 라우팅하고 증거가 요청을 충족하는지 판단한다. 단순 중계자가 아니다. |
| **PlanAI** | `PLAN`을 작성하고, 구현이 계획과 일치하는지 검토 증거를 정리한다. `REVIEW`는 승인이 아니며 발견은 `blocker`(반드시 수정) 또는 `nit`(비차단)으로 표시한다. |
| **ExecAI** | `PLAN`을 구현·검증하고, 범위를 넓히지 않은 채 모호함은 LeadAI에게 되돌린다. |

PlanAI와 ExecAI는 **LeadAI를 통해서만** 통신한다.

**강제 선행 규칙:** 기록이 필요한 각 단계의 담당 역할은 단계 시작 시 `.agent-relay/GUIDANCE.md`와 `.agent-relay/LESSON-LEARNED.md` 인덱스를 읽고, 현재 범위와 `Applies When` 또는 `Trigger / Symptom`이 맞는 개별 기록만 `.agent-relay/lesson-learned/`에서 읽는다. LeadAI의 초기 선별만 의존하지 않고, PlanAI와 ExecAI도 자신의 확장된 단계 범위에 맞춰 다시 선별한다. 명백한 기록 제외 요청은 이 확인 없이 응답할 수 있다.

---

## 작업 분류

세션 시작 시 LeadAI는 이번 Agent Relay 세션에서 Git 브랜치 전략을 사용할지 묻는다: 항상 브랜치 사용, 브랜치 사용 안 함, 작업마다 확인.

LeadAI는 먼저 요청이 명백한 기록 제외 대상인지 가볍게 판단한다. 기록이 필요한 요청이면 필수 지침·교훈 확인을 마친 뒤 `Direct` 또는 `Standard`로 분류한다. Agent Relay **부트스트랩**과 **업데이트**(`.agent-relay/`·Agent Relay 지시 파일 동기화)는 기록 제외가 아니며, LeadAI가 직접 수행하면 `Direct`로 `REQUEST → RUN_DONE`을 기록한다.

| 분류 | 일반적 범위 | 처리 방식 |
| --- | --- | --- |
| 기록 제외 | 단순 질문 답변, 짧은 설명, 브레인스토밍 | 응답만 하고 이벤트를 남기지 않음 |
| `Direct` | 사소한 텍스트/설정 변경, 명백한 국소 편집, Agent Relay 부트스트랩·업데이트 | LeadAI가 직접 처리하고 `REQUEST → RUN_DONE` 기록 |
| `Standard` | 다중 파일 구현, 설계 판단, 구현 검증이 필요한 작업 | 세션 브랜치 전략 적용 → PlanAI → ExecAI → PlanAI 검토 → 승인 후 필요 시 자동 병합 |

---

## 백그라운드 위임

`Standard` 작업은 분류 직후 세션 시작 때 정한 Git 브랜치 전략을 따른다. 전용 작업 브랜치를 쓰는 경우 현재 브랜치를 기준 브랜치로 기억하고 작업 브랜치를 만든 뒤, 그 브랜치에서 `REQUEST`부터 기록한다. 브랜치를 쓰지 않는 전략이면 현재 브랜치에서 기록과 변경을 진행한다. 전용 작업 브랜치를 쓴 경우 승인 후 `CLOSE`을 기록해 승인 상태를 커밋한 다음 기준 브랜치로 자동 병합한다. 위임은 가능한 한 백그라운드로 수행하며, LeadAI는 위임 직후 사용자에게 짧은 상태를 반환하고 완료 대기·폴링·sleep으로 사용자 응답을 막지 않는다.

---

## 이벤트 타임라인

`relay.log`는 작업 이력을 추가-전용 이벤트 타임라인으로 남긴다.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

- `timestamp`는 로컬 시스템 시간 기준 `YYYY-MM-DDTHH:MM:SS` 형식으로 기록한다.
- `task-id`는 무작위 소문자 영문 4글자를 쓴다.
- LeadAI는 `REQUEST` 기록 시 `task-id` 하나를 정하고, 같은 Standard 작업의 `PLANNED`/`EXECUTED`/`REVIEW`/`FEEDBACK`/`CLOSE`까지 재사용한다. 새 `REQUEST`마다 새 `task-id`를 쓴다.
- 이벤트는 `REQUEST`, `PLANNED`, `EXECUTED`, `REVIEW`, `FEEDBACK`, `CLOSE`, `RUN_DONE`만 쓴다.
- LeadAI 직접 처리 흐름은 `REQUEST → RUN_DONE`이다.
- 표준 처리 흐름은 `REQUEST` → `PLANNED` → `EXECUTED` → `REVIEW` → `CLOSE`이다.
- `Standard`의 `REQUEST`는 세션 브랜치 전략을 적용한 뒤 기록한다. 전용 작업 브랜치를 쓰는 경우 승인 전에는 기준 브랜치에 해당 작업의 이벤트나 변경을 기록하지 않는다.
- `FEEDBACK`은 사용자가 `CLOSE` 승인 전 피드백·결함을 알려줄 때 LeadAI가 기록한다. 같은 `task-id`와 산출물 파일 키를 유지한다.
- `FEEDBACK` 후 LeadAI는 **현재 런에 추가**할지 **새로운 런**으로 돌릴지 사용자에게 묻는다. 명백한 결함이면 사용자 확인 없이 현재 런에 추가한다.
- **현재 런에 추가**: 마지막 `RUN-<NN>` 범위와 기존 `PLAN` 안에서 `EXECUTED` → `REVIEW-<NN>`을 다시 진행한다. `CLOSE` 이벤트 승인 전이면 `RUN-<NN>.md` 갱신을 허용한다.
- **새로운 런**: 다음 `RUN-<NN+1>`로 `EXECUTED` → `REVIEW-<NN+1>`을 진행한다.
- `role` 주변 공백은 정렬용이며 의미가 없다.
- `event`는 8자 폭, `role`은 6자 폭으로 왼쪽 정렬하고 부족한 자리는 공백으로 채운다.
- `EXECUTED`는 ExecAI가 `RUN-<NN>.md`를 작성한 뒤 LeadAI가 기록한다. `path` 필수.
- 한 라운드 `<NN>`은 하나의 `EXECUTED`와 그에 대응하는 `REVIEW`로 식별한다. `blocker`로 다음 라운드를 돌릴 때 같은 `task-id`에 새 `EXECUTED`/`REVIEW`를 추가한다.

---

## 산출물

모든 라운드 산출물은 `.agent-relay/runs/`에 작성하며, 작업당 하나의 안정된 `<YYYYMMDD>-<HHMM>-<slug>` 키를 쓴다. **이전 라운드를 덮어쓰지 않는다.**

| 산출물 | 경로 | 작성자 | 의미 |
| --- | --- | --- | --- |
| Plan | `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-PLAN.md` | PlanAI | 계획과 성공 기준 |
| Submission | `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-RUN-<NN>.md` | ExecAI | 라운드 `<NN>`의 변경/검증/리스크 |
| Review | `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-REVIEW-<NN>.md` | PlanAI | 같은 라운드 발견 |
| Closure | `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-CLOSE.md` | LeadAI | 사용자 승인 이후의 최종 종료 기록 |

- `<NN>`은 `01`부터 시작한다.
- 이전 라운드를 덮어쓰지 않는다. 예외: `FEEDBACK` 후 현재 런에 추가할 때, `CLOSE` 이벤트 승인 전이면 `RUN-<NN>.md` 갱신을 허용한다.
- `<SLUG>`는 LeadAI가 정한 소문자 kebab-case 작업 키를 쓴다.
- `<YYYYMMDD>`와 `<HHMM>`은 LeadAI가 `REQUEST`를 기록할 때의 로컬 시스템 날짜·시분(24시간, 구분자 없음)을 쓴다. 예: `20260526-1430-diary-write`.
- `task-id`는 `relay.log` 이벤트 식별자이고, `<YYYYMMDD>-<HHMM>-<SLUG>`는 `.agent-relay/runs/` 산출물 파일 키다. 같은 Standard 작업에서는 `task-id` 하나와 파일 키 하나를 함께 쓴다.
- 같은 작업의 모든 라운드 산출물은 같은 `<YYYYMMDD>-<HHMM>-<SLUG>` 키를 쓴다.
- 산출물은 `.agent-relay/templates/plan.md`, `run.md`, `review.md`, `close.md` 형식을 따른다.
- ExecAI는 절대 `CLOSE`을 쓰지 않는다.
- `REVIEW`는 사용자 결정을 위한 증거이지 승인이 아니다. PlanAI는 완료를 승인하거나 `CLOSE` 산출물을 작성하지 않는다.
- 사용자가 명시적으로 승인한 뒤에만 LeadAI가 `CLOSE` 산출물을 작성하고 `CLOSE` 이벤트를 기록하여 작업을 종료할 수 있다.
- 어느 `REVIEW` 이후든 사용자의 `FEEDBACK`은 정상 파이프라인 단계이며, 완료된 작업의 예외적 되돌림으로 취급하지 않는다.
- 각 `RUN`은 변경 파일, 변경 요약, 테스트/검증, 미해결 리스크를 기록한다.
- 산출물 작성과 `relay.log` 이벤트 추가는 별개의 필수 작업이다. PlanAI와 ExecAI는 산출물 완료와 제안 summary를 LeadAI에게 통지할 뿐 `relay.log`를 쓰거나 기록 완료를 주장하지 않는다. 각 단계는 LeadAI가 해당 이벤트를 추가하고 확인한 뒤에만 완료된 것으로 본다.
- 모든 `relay.log` 이벤트는 LeadAI가 추가한다. 가능하면 직접 `echo >> relay.log` 대신 `.agent-relay/protocol-guard append ...`를 사용한다.
- LeadAI는 다음 단계 위임 전에 `.agent-relay/protocol-guard gate ...` 또는 `tail -50 .agent-relay/relay.log`로 직전 단계 이벤트가 `relay.log`에 추가됐는지 확인한다. 확인하지 못하면 다음 단계 위임을 중단하고 LeadAI 소유의 로그 추가 또는 수정 작업을 완료한다.

필수 게이트:

| 다음 단계 | 확인할 직전 이벤트 |
| --- | --- |
| ExecAI 위임 | `PLANNED` |
| PlanAI 검토 위임 | `EXECUTED` |
| 사용자 승인 요청 | `REVIEW` |
| 최종 `CLOSE` 이벤트 | 명시적 사용자 승인 |

---

## 파이프라인

1. LeadAI가 요청을 분류한다.
2. 기록 제외 대상이면 응답만 하고 이벤트를 남기지 않는다.
3. `Direct`이면 LeadAI가 직접 처리하고 `REQUEST → RUN_DONE` 이벤트 흐름으로 작업을 닫는다.
4. `Standard`이면 LeadAI가 세션 Git 브랜치 전략을 적용한 뒤 `REQUEST`를 기록한다.
5. PlanAI가 상단 `LeadAI Brief`를 포함한 `PLAN`을 작성한다.
6. LeadAI는 기본적으로 `LeadAI Brief`만 읽고 완전성을 확인한 뒤, 그 안의 `ExecAI Prompt`로 ExecAI에게 위임한다.
7. ExecAI는 `PLAN`·성공 기준·범위에 따라 구현한 뒤 `RUN-01`을 쓰고, LeadAI가 `EXECUTED`를 기록한다.
8. LeadAI가 `EXECUTED` 이벤트를 확인한 뒤 PlanAI에게 검토를 위임하고, PlanAI는 해당 `RUN` 경로를 받아 같은 번호의 `REVIEW`를 쓴다.
9. LeadAI가 `REVIEW` 이벤트를 확인한다. `blocker`가 없으면 사용자 결정 준비가 된 것이지 승인이 아니다. LeadAI는 결과·검증, 존재하는 조치 대상 nit·리스크, `REVIEW` 산출물 경로를 사용자에게 보고하고 승인을 요청한다.
10. 사용자가 명시적으로 승인한 뒤에만 LeadAI가 `CLOSE` 산출물을 작성하고 `CLOSE` 이벤트를 기록한 뒤 승인된 상태를 필요에 따라 커밋한다. 전용 작업 브랜치를 사용했다면 기준 브랜치로 자동 병합한다. 커밋 또는 병합에 문제가 있으면 강제하지 않고 blocker로 보고한다.
11. 사용자가 승인 대신 피드백·결함을 알려주면 LeadAI가 작업 브랜치에 `FEEDBACK`을 기록한다. 명백한 결함이면 현재 런에 추가하고, 그렇지 않으면 **현재 런에 추가 / 새로운 런** 중 사용자 선택을 받는다. 이후 ExecAI 작업부터 다시 진행한다.
12. `CLOSE` 승인을 받은 LeadAI는 해당 세션에서 발생한 착오, 해결 방법, 사용자 의견을 종합하여 `.agent-relay/GUIDANCE.md` 수정안 또는 `.agent-relay/lesson-learned/` 추가안을 사용자에게 제안한다. 사용자가 수락한 항목만 기록한다.
13. `blocker`가 있으면 LeadAI가 다음 라운드를 ExecAI에게 위임하고, ExecAI가 `RUN-<NN>`을 쓰면 LeadAI가 `EXECUTED`를 기록한 뒤 PlanAI가 다음 `REVIEW`를 쓴다. `REVIEW-03` 전까지 사용자 승인 없이 진행한다.
14. `REVIEW-03`까지도 `blocker`가 남으면 LeadAI는 상태를 보고하고 사용자에게 **재시도 / 계획 수정 / 부분 수락 / 중단** 중 선택을 요청한다.

Standard 작업에서 사용자 개입이 필요한 경우는 `CLOSE` 최종 승인, `CLOSE` 승인 전 피드백·결함(`FEEDBACK`)과 `FEEDBACK` 후 현재 런·새 런 선택, `REVIEW-03` 이후에도 `blocker`가 남는 경우, LeadAI가 사용자 결정이 필요하다고 판단한 경우뿐이다. 전용 작업 브랜치를 쓴 경우 승인이 끝나면 병합은 자동으로 진행하며 별도 확인을 받지 않는다.

`Direct` 작업은 사용자 완료 승인 없이 `RUN_DONE`으로 닫을 수 있다. 다만 장기 지침이나 재사용 가능한 교훈이 생겼다면 LeadAI는 사용자에게 기록안을 제안하고, 사용자가 수락한 항목만 `GUIDANCE.md` 또는 `lesson-learned/`에 추가한다.

---

## 위임과 보고

### 컨텍스트 관리

LeadAI, PlanAI, ExecAI는 한 작업 안에서 컨텍스트 교체 없이 연속 사용하는 것을 전제로 한다.

LeadAI는 컨텍스트가 불필요하게 커지지 않도록 결과를 받을 때 산출물 **경로 + 최소 결정 정보**만 보관한다: 한 줄 결과, 한 줄 검증 상태, 해당 시 `blocker` 건수/요약, 잔존 리스크나 사용자 결정 요구. PlanAI는 각 `PLAN` 상단에 목표, 범위, 성공 기준, 리스크, 필수 확인, 최소 `ExecAI Prompt`를 담은 `LeadAI Brief`를 작성한다. LeadAI는 ExecAI 위임 전 기본적으로 이 브리프만 읽고, 브리프가 누락·불완전·모순·고위험이거나 사용자 결정에 상세 검토가 필요한 경우에만 `PLAN` 전문을 읽는다.

PlanAI/ExecAI는 가능하면 같은 컨텍스트를 유지하되, 다음 중 하나라도 발생하면 LeadAI가 사용자에게 **교체 여부**를 물어야 한다.

- 한 인스턴스에 후속 메시지가 5개 이상 누적
- 작업 주제가 명백히 바뀜
- 응답이 눈에 띄게 느려지거나 이전 컨텍스트를 혼동

### 위임 시 필수 필드

위임 프롬프트에는 명시 요구사항, 읽을/쓸 경로, LeadAI가 추가할 이벤트만 전달한다. 성공 기준·검증·리스크는 PlanAI가 `PLAN`에서 정의하고 ExecAI와 검토자는 이를 따른다. 하위 AI는 범위를 임의로 넓히거나 `relay.log`를 쓰지 않고, 완료와 제안 summary만 LeadAI에게 알린다.

### 보고 시 필수 필드

PlanAI가 LeadAI에게 보고할 때는 다음만 간결히 포함한다.

- `PLAN`/`REVIEW` 산출물 경로
- LeadAI가 append할 이벤트명과 제안 summary
- 판단 결과
- `blocker` 수와 요약
- `nit` 요약
- 사용자 결정 필요 여부

ExecAI가 LeadAI에게 보고할 때는 다음만 간결히 포함한다.

- `RUN` 산출물 경로
- LeadAI가 append할 이벤트명과 제안 summary
- 변경 요약
- 검증 결과
- 미해결 리스크
- 범위 밖으로 넘긴 사항

### 사용자에게 보여 주는 보고

사용자 대상 보고는 기본적으로 짧게 작성한다. `Direct` 작업은 결과, 핵심
변경 범위, 검증만 1~3문장으로 알린다. 생성·보존 파일 전체 목록, 프로토콜
진행 설명, 비어 있는 리스크/다음 단계 섹션은 사용자가 요청하거나 조치가
필요할 때만 포함한다.

`Standard` 작업의 승인 요청도 결과, 검증 상태, 조치가 필요한
blocker/리스크, `REVIEW` 경로만 우선 보여 준다. 상세
변경과 증거는 요청받지 않는 한 산출물에 둔다.

## GUIDANCE 사용 기준

`.agent-relay/GUIDANCE.md`는 세션이 바뀌어도 계속 지켜야 하는 장기 지침과 프로젝트 제약을 단일 파일로 누적하는 문서다. `CLOSE` 승인 이후 LeadAI가 장기적으로 유지할 사용자 의견이나 제약의 반영안을 제안하고, 사용자가 수락한 경우에만 갱신한다.

기록할 내용:

- 사용자 선호와 상시 지시
- 기술·제품·운영 제약
- 보안·개인정보 규칙과 금지사항
- 코드, 테스트, 브랜치 등 반복 적용할 작업 관례

기록하지 않을 내용:

- 현재 작업의 진행 상태
- 임시 계획이나 디버깅 중간 메모
- 한 번만 적용되는 요청

진행 상태와 작업 결과는 `relay.log`와 `.agent-relay/runs/` 산출물에 기록한다.

## LESSON-LEARNED 사용 기준

`.agent-relay/LESSON-LEARNED.md`는 `.agent-relay/lesson-learned/`에 저장된 실제 기록의 검색 인덱스다. 인덱스에는 각 기록의 `Applies When`, `Trigger / Symptom`, 파일 경로, 해결 요약을 남긴다. 기록이 필요한 각 단계의 담당 역할은 인덱스 전체를 훑은 뒤 현재 범위에 맞는 상세 기록만 읽는다. 완료된 작업에서 새롭게 알게 된 착오, 오해하기 쉬운 전제, 다시 활용할 수 있는 해결 방법은 `.agent-relay/lesson-learned/`에 누적하고, 수락된 기록은 인덱스에 함께 추가한다. `CLOSE` 승인 이후 LeadAI가 기록안을 제안하고, 사용자가 수락한 경우에만 추가한다.

기록 파일은 `.agent-relay/templates/lesson-learned.md` 형식을 따르며 `.agent-relay/lesson-learned/<YYYYMMDD>-<trigger-or-symptom>.md` 이름으로 저장한다.

기록할 내용:

- 잘못 가정했던 원인이나 실패한 접근과 그 이유
- 문제를 해결한 핵심 방법과 재현 가능한 확인 방법
- 다음 작업에서 같은 착오를 피하기 위해 확인할 사항
- 적용되는 작업 범위·환경·도구(`Applies When`)와 검색 가능한 현상(`Trigger / Symptom`)

`GUIDANCE.md`는 앞으로 지켜야 하는 지침과 제약을 담고, `lesson-learned/`는 완료된 작업에서 얻은 해결 지식을 담는다. 단순 작업 결과나 진행 기록은 `relay.log`와 `.agent-relay/runs/` 산출물에 둔다.

## 안전·금지

- `.agent-relay/`에 비밀정보, 자격증명, 고객정보, 민감한 운영정보를 저장하지 않는다.

자세한 사양은 [`bootstrap/.agent-relay/PROTOCOL.md`](bootstrap/.agent-relay/PROTOCOL.md), 한국어 가이드는 [`PROTOCOL-GUIDE.md`](PROTOCOL-GUIDE.md)를 참고하세요.

---

## 매 프롬프트마다 쓰기 좋은 한 줄

```text
Agent Relay 규칙에 따라 진행해.
```

조금 더 명확하게:

```text
AGENTS.md 또는 CLAUDE.md와 .agent-relay/PROTOCOL.md 기준으로 진행해.
새 세션이면 Agent Relay의 읽기 순서를 먼저 따라줘.
```

---

## 더 읽기

- [`BOOTSTRAP.md`](BOOTSTRAP.md) — 에이전트용 부트스트랩 절차 (단계별, 머지/중단 분기 포함)
- [`PROTOCOL-GUIDE.md`](PROTOCOL-GUIDE.md) — 한국어 가이드 (운영 원칙, 부트스트랩 절차, 파일별 역할 정리)
- [`bootstrap/.agent-relay/PROTOCOL.md`](bootstrap/.agent-relay/PROTOCOL.md) — 영문 정식 사양 (핵심지침, 이벤트 타입, 기타 규칙)
