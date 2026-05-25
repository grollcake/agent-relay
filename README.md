# Agent Relay

> **🤖 AI 에이전트에게**: 당신이 AI 에이전트이고 사용자로부터 이 저장소를 본인 프로젝트에 반영해달라는 요청을 받았다면, [`BOOTSTRAP.md`](BOOTSTRAP.md)를 그대로 따라 부트스트랩을 수행하세요. 기존 `AGENTS.md`는 덮어쓰지 말고 짧은 Agent Relay 섹션만 추가하세요. 전체 `bootstrap/AGENTS.md`는 대상 프로젝트에 `AGENTS.md`가 없을 때만 복사합니다. 이미 `.agent-relay/`가 있으면 새 설치가 아니라 업데이트 요청인지 확인하고, 업데이트 요청이면 `.agent-relay/VERSION`을 읽은 뒤 기본 upstream `https://github.com/grollcake/agent-relay`와 비교하세요. `bootstrap/` 트리는 목적지 경로를 그대로 미러링한 복사 원본입니다.

Agent Relay는 **PM / Planner / Executor 에이전트 팀**이 역할을 나누고, 기록 기반으로 작업을 이어가기 위한 파일 기반 협업 규약입니다.

이 저장소는 **사양과 템플릿을 모아둔 가이드 저장소**입니다. 사용자 프로젝트에 Agent Relay를 도입하면, 그 프로젝트에 `.agent-relay/`가 생깁니다.

---

## 목적

1. PM / Planner / Executor 에이전트 팀의 역할 분담으로 속도와 완전성을 높인다.
2. PM은 Planner/Executor에게 백그라운드로 위임하여 항상 사용자에 대응 가능한 상태로 유지한다.
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

에이전트는 대상 프로젝트의 `.agent-relay/VERSION`을 읽고 기본 upstream `https://github.com/grollcake/agent-relay`와 비교합니다. `AGENTS.md`는 프로젝트 고유 지침을 보존하고 Agent Relay 섹션만 최신 `bootstrap/AGENTS.md`와 비교해 갱신합니다. `LESSON-LEARNED.md`는 안내 문서이므로 로컬 수정이 없을 때만 갱신하고, `GUIDANCE.md`, `lesson-learned/`, `relay.log`, `runs/`는 프로젝트별 상태이므로 덮어쓰지 않습니다.

---

## 목표 부트스트랩 트리

아래 구조는 `runs/` 중심 흐름에 맞춘 현재 `bootstrap/` 배포 템플릿 구조입니다.

```text
bootstrap/
├── AGENTS.md                          # 프로젝트 루트로 복사 (필수)
├── CLAUDE.md                          # Claude Code에서 실행 중이면 생성, 이미 있으면 머지 (선택)
├── .codex/
│   └── instructions.md                # 이미 있으면 머지 (선택)
├── .cursor/
│   └── rules/
│       └── agent-relay.mdc            # .cursor/ 디렉토리가 있으면 신규 (선택)
└── .agent-relay/                      # 프로젝트 루트로 복사 (필수)
    ├── PROTOCOL.md                    # Agent Relay 규칙(필수)
    ├── VERSION                        # 설치 버전
    ├── GUIDANCE.md                    # 장기 지침/제약 누적
    ├── LESSON-LEARNED.md              # 완료 작업에서 얻은 해결 지식 기록 안내
    ├── relay.log                      # 추가 전용 이벤트 로그
    ├── lesson-learned/                # 완료 작업에서 얻은 해결 지식 기록
    │   └── .gitkeep
    ├── runs/                          # 라운드 산출물(PLAN/RUN/REVIEW/DONE)
    │   └── .gitkeep
    └── templates/                     # 산출물·누적 문서 양식
        ├── guidance.md
        ├── lesson-learned.md
        ├── plan.md
        ├── run.md
        ├── review.md
        └── done.md
```

`bootstrap/` 안의 모든 경로는 **목적지 경로 그대로**입니다. 도구별 지시 파일(`CLAUDE.md`/`.codex/`)은 현재 실행 중이거나 사용 흔적이 있는 도구에만 생성 또는 머지되고, 안 쓰면 생성되지 않습니다.

---

## 에이전트 팀 구성

**PM / Planner / Executor 에이전트 팀**을 구성하고 이 프로토콜을 적용한다. 리더는 **PM (허브)**이며 사용자 소통, 분류, 범위/위험 결정, 위임, 결과 해석, 최종 보고를 담당한다. Planner와 Executor는 PM을 통해서만 통신한다.

| 역할 | 책임 |
| --- | --- |
| **PM (허브)** | 작업을 라우팅하고 증거가 요청을 충족하는지 판단한다. 단순 중계자가 아니다. |
| **Planner** | `PLAN`을 작성하고, 구현이 계획과 일치하는지 검토한다. 발견은 `blocker`(반드시 수정) 또는 `nit`(비차단)으로 표시한다. |
| **Executor** | `PLAN`을 구현·검증하고, 범위를 넓히지 않은 채 모호함은 PM에게 되돌린다. |

Planner와 Executor는 **PM을 통해서만** 통신한다.

**강제 선행 규칙:** PM, Planner, Executor는 기록이 필요한 작업에 착수하기 전에 반드시 `.agent-relay/GUIDANCE.md`, `.agent-relay/LESSON-LEARNED.md`, `.agent-relay/lesson-learned/`를 읽어 현재 적용할 지침과 이전 해결 지식을 확인한다. 명백한 기록 제외 요청은 이 확인 없이 응답할 수 있지만, 파일 변경·조사·설계 판단·프로젝트 지침 의존 답변으로 넘어가면 먼저 이 확인을 완료해야 한다.

---

## 작업 분류

PM은 먼저 요청이 명백한 기록 제외 대상인지 가볍게 판단한다. 기록이 필요한 요청이면 필수 지침·교훈 확인을 마친 뒤 `Trivial` 또는 `Standard`로 분류한다.

| 분류 | 일반적 범위 | 처리 방식 |
| --- | --- | --- |
| 기록 제외 | 단순 질문 답변, 짧은 설명, 브레인스토밍 | 응답만 하고 이벤트를 남기지 않음 |
| `Trivial` | 사소한 텍스트/설정 변경, 명백한 국소 편집 | PM이 직접 처리하고 `REQUEST → RUN_DONE` 기록 |
| `Standard` | 다중 파일 구현, 설계 판단, 구현 검증이 필요한 작업 | Planner → Executor → Planner 검토 |

---

## 백그라운드 위임

`Standard` 작업은 가능한 한 백그라운드로 위임한다. PM은 위임 후에도 사용자 응답을 계속 담당한다. 완료를 기다리기 위한 폴링이나 sleep은 하지 않는다. 백그라운드 실행이 어렵다면 같은 단계와 산출물 전달 방식을 순차적으로 따른다.

---

## 이벤트 타임라인

`relay.log`는 작업 이력을 추가-전용 이벤트 타임라인으로 남긴다.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

- `timestamp`는 KST 기준 `YYYY-MM-DDTHH:MM:SS` 형식으로 기록한다.
- `task-id`는 무작위 소문자 영문 4글자를 쓴다.
- 이벤트는 `REQUEST`, `PLAN`, `RUN_ST`, `RUN_ED`, `REVIEW`, `DONE`, `RUN_DONE`만 쓴다.
- PM 직접 처리 흐름은 `REQUEST → RUN_DONE`이다.
- 표준 처리 흐름은 `REQUEST → PLAN → RUN_ST → RUN_ED → REVIEW → DONE`이다.
- `role` 주변 공백은 정렬용이며 의미가 없다.
- `RUN_ST`는 Executor가 `RUN-<NN>` 착수 시 기록한다. `path` 없음.
- `RUN_ED`는 Executor가 `RUN-<NN>.md` 작성 후 기록한다. `path` 필수.

---

## 산출물

모든 라운드 산출물은 `.agent-relay/runs/`에 작성하며, 작업당 하나의 안정된 `<date>-<slug>` 키를 쓴다. **이전 라운드를 덮어쓰지 않는다.**

| 산출물 | 경로 | 작성자 | 의미 |
| --- | --- | --- | --- |
| Plan | `.agent-relay/runs/<YYYYMMDD>-<SLUG>-PLAN.md` | Planner | 계획과 성공 기준 |
| Submission | `.agent-relay/runs/<YYYYMMDD>-<SLUG>-RUN-<NN>.md` | Executor | 라운드 `<NN>`의 변경/검증/리스크 |
| Review | `.agent-relay/runs/<YYYYMMDD>-<SLUG>-REVIEW-<NN>.md` | Planner | 같은 라운드 발견 |
| Acceptance | `.agent-relay/runs/<YYYYMMDD>-<SLUG>-DONE.md` | Planner | `blocker` 없이 검토를 통과한 수락 결과 |

- `<NN>`은 `01`부터 시작한다.
- `<SLUG>`는 PM이 정한 소문자 kebab-case 작업 키를 쓴다.
- `task-id`는 `relay.log` 이벤트 식별자이고, `<SLUG>`는 `.agent-relay/runs/` 산출물 파일 키다.
- 같은 작업의 모든 라운드 산출물은 같은 `<YYYYMMDD>-<SLUG>` 키를 쓴다.
- 산출물은 `.agent-relay/templates/plan.md`, `run.md`, `review.md`, `done.md` 형식을 따른다.
- Executor는 절대 `DONE`을 쓰지 않는다.
- Planner는 검토에 `blocker`가 없을 때 `DONE` 산출물을 작성한다. `nit`는 `DONE`에 기록할 수 있다.
- 사용자가 명시적으로 승인하기 전에는 PM이 `DONE` 이벤트를 기록하여 작업을 종료할 수 없다.
- 각 `RUN`은 변경 파일, 변경 요약, 테스트/검증, 미해결 리스크를 기록한다.
- 산출물 작성과 `relay.log` 이벤트 추가는 별개의 필수 작업이다. 각 단계는 산출물 작성과 해당 이벤트 추가가 모두 끝나야 완료된 것으로 본다.
- 산출물 작성자가 해당 이벤트를 추가한다. Planner는 `PLAN`/`REVIEW`, Executor는 `RUN_ST`/`RUN_ED`, PM은 `REQUEST`/`RUN_DONE`과 사용자 승인 후 최종 `DONE`을 기록한다.
- PM은 다음 단계 위임 전에 직전 단계 이벤트가 `relay.log`에 추가됐는지 확인한다.

---

## 파이프라인

1. PM이 요청을 분류한다.
2. 기록 제외 대상이면 응답만 하고 이벤트를 남기지 않는다.
3. `Trivial`이면 PM이 직접 처리하고 `REQUEST → RUN_DONE` 이벤트 흐름으로 작업을 닫는다.
4. `Standard`이면 Planner가 `PLAN`을 작성한다.
5. PM이 Executor에게 위임한다.
6. Executor는 착수 시 `RUN_ST`를 기록하고, `PLAN`·성공 기준·범위에 따라 구현한 뒤 `RUN-01`을 쓰고 `RUN_ED`를 기록한다.
7. Planner는 해당 `RUN` 경로를 받아 같은 번호의 `REVIEW`를 쓴다.
8. `blocker`가 없으면 Planner가 `DONE` 산출물을 쓴다. PM은 결과·nit·리스크와 `DONE` 산출물 경로를 사용자에게 보고하고 승인을 요청한다.
9. 사용자가 명시적으로 승인한 뒤에만 PM이 `DONE` 이벤트를 기록하여 작업을 닫는다.
10. `DONE` 승인을 받은 PM은 해당 세션에서 발생한 착오, 해결 방법, 사용자 의견을 종합하여 `.agent-relay/GUIDANCE.md` 수정안 또는 `.agent-relay/lesson-learned/` 추가안을 사용자에게 제안한다. 사용자가 수락한 항목만 기록한다.
11. `blocker`가 있으면 Executor가 다음 라운드에서 `RUN_ST` → `RUN-<NN>` → `RUN_ED`를, Planner가 다음 `REVIEW`를 쓴다.
12. `REVIEW-03`까지도 `blocker`가 남으면 PM은 상태를 보고하고 사용자에게 **재시도 / 계획 수정 / 부분 수락 / 중단** 중 선택을 요청한다.

`Trivial` 작업은 사용자 완료 승인 없이 `RUN_DONE`으로 닫을 수 있다. 다만 장기 지침이나 재사용 가능한 교훈이 생겼다면 PM은 사용자에게 기록안을 제안하고, 사용자가 수락한 항목만 `GUIDANCE.md` 또는 `lesson-learned/`에 추가한다.

---

## 위임과 보고

### 컨텍스트 관리

PM, Planner, Executor는 한 작업 안에서 컨텍스트 교체 없이 연속 사용하는 것을 전제로 한다.

PM은 컨텍스트가 불필요하게 커지지 않도록 결과를 받을 때 산출물 **경로 + 최소 결정 정보**만 보관한다: 한 줄 결과, 한 줄 검증 상태, 해당 시 `blocker` 건수/요약, 잔존 리스크나 사용자 결정 요구. 결정적 모호함이 없는 한 전체 산출물을 PM 컨텍스트에 적재하지 않는다.

Planner/Executor는 가능하면 같은 컨텍스트를 유지하되, 다음 중 하나라도 발생하면 PM이 사용자에게 **교체 여부**를 물어야 한다.

- 한 인스턴스에 후속 메시지가 5개 이상 누적
- 작업 주제가 명백히 바뀜
- 응답이 눈에 띄게 느려지거나 이전 컨텍스트를 혼동

### 위임 시 필수 필드

Planner/Executor에게 보내는 모든 프롬프트는 자기완결적이어야 하며 다음을 포함한다.

- 목표(goal)
- 관련 파일 또는 조사 범위
- 산출물 타입과 정확한 경로
- 성공 기준과 검증 방법
- 범위 외 작업 금지 명시
- 불명확한 사항은 추정하지 말고 PM에게 되돌릴 것
- 단계에 필요한 입력 산출물 경로

### 보고 시 필수 필드

Planner가 PM에게 보고할 때는 다음만 간결히 포함한다.

- `PLAN`/`REVIEW`/`DONE` 산출물 경로
- 판단 결과
- `blocker` 수와 요약
- `nit` 요약
- 사용자 결정 필요 여부

Executor가 PM에게 보고할 때는 다음만 간결히 포함한다.

- `RUN` 산출물 경로
- 변경 요약
- 검증 결과
- 미해결 리스크
- 범위 밖으로 넘긴 사항

## GUIDANCE 사용 기준

`.agent-relay/GUIDANCE.md`는 세션이 바뀌어도 계속 지켜야 하는 장기 지침과 프로젝트 제약을 단일 파일로 누적하는 문서다. `DONE` 승인 이후 PM이 장기적으로 유지할 사용자 의견이나 제약의 반영안을 제안하고, 사용자가 수락한 경우에만 갱신한다.

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

`.agent-relay/LESSON-LEARNED.md`는 해결 지식 기록의 목적과 작성 방식을 설명하는 안내 문서다. 완료된 작업에서 새롭게 알게 된 착오, 오해하기 쉬운 전제, 다시 활용할 수 있는 해결 방법은 `.agent-relay/lesson-learned/`에 누적한다. `DONE` 승인 이후 PM이 기록안을 제안하고, 사용자가 수락한 경우에만 추가한다.

기록 파일은 `.agent-relay/templates/lesson-learned.md` 형식을 따르며 `.agent-relay/lesson-learned/<YYYYMMDD>-<slug>.md` 이름으로 저장한다.

기록할 내용:

- 잘못 가정했던 원인이나 실패한 접근과 그 이유
- 문제를 해결한 핵심 방법과 재현 가능한 확인 방법
- 다음 작업에서 같은 착오를 피하기 위해 확인할 사항

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
AGENTS.md와 .agent-relay/PROTOCOL.md 기준으로 진행해.
새 세션이면 Agent Relay의 읽기 순서를 먼저 따라줘.
```

---

## 더 읽기

- [`BOOTSTRAP.md`](BOOTSTRAP.md) — 에이전트용 부트스트랩 절차 (단계별, 머지/중단 분기 포함)
- [`PROTOCOL-GUIDE.md`](PROTOCOL-GUIDE.md) — 한국어 가이드 (운영 원칙, 부트스트랩 절차, 파일별 역할 정리)
- [`bootstrap/.agent-relay/PROTOCOL.md`](bootstrap/.agent-relay/PROTOCOL.md) — 영문 정식 사양 (핵심지침, 이벤트 타입, 기타 규칙)
