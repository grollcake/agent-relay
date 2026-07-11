# Agent Relay

Agent Relay는 AI 에이전트가 계획·구현·검토 역할을 나누고, 파일로 작업 맥락을 이어가는 협업 규약이다. 에이전트나 세션이 바뀌어도 기록을 바탕으로 작업을 계속할 수 있게 한다.

> **AI 에이전트에게**: 이 저장소를 다른 프로젝트에 반영해달라는 요청을 받으면 [부트스트랩 절차](BOOTSTRAP.md)를 따른다. 기존 프로젝트 지침과 Agent Relay 상태는 덮어쓰지 않는다.

## 핵심 개념

| 역할 | 책임 |
| --- | --- |
| Director | 사용자 소통, 작업 분류, 위임, 상태 기록, 최종 보고 |
| Planner | 계획 작성과 구현 결과 검토 |
| Executor | 계획에 따른 구현과 검증 |

Director는 Planner와 Executor에게 작업을 위임하고 사용자 응답 창구로 남는다. Planner와 Executor는 Director를 통해서만 작업을 주고받는다.

작업은 파일에 기록한다.

- `relay.log`: 작업 이벤트 타임라인
- `runs/`: 작업별 계획, 실행, 검토, 종료 기록
- `GUIDANCE.md`: 세션이 바뀌어도 유지할 지침과 제약
- `lesson-learned/`: 완료된 작업에서 얻은 재사용 가능한 지식

## 설치와 업데이트

대상 프로젝트에서 에이전트에게 다음처럼 요청한다.

```text
github.com/grollcake/agent-relay 를 내 프로젝트에 반영해줘
```

이미 설치된 경우에는 다음처럼 요청한다.

```text
agent-relay 최신화해줘
```

설치와 업데이트의 파일별 처리 규칙은 [BOOTSTRAP.md](BOOTSTRAP.md)에 있다.

### 지원 환경

- macOS: amd64, arm64
- Linux: amd64, arm64
- Windows: amd64, arm64; 운영 셸은 Git for Windows의 Git Bash

Agent Relay는 Go로 컴파일된 단일 네이티브 바이너리이며 설치된 프로젝트에는
Go 런타임이 필요하지 않다. Windows 바이너리도 네이티브로 동작하지만 저장소와
Git 작업의 기준 셸은 Git Bash로 고정한다. 소스 빌드와 테스트에는 Go 1.26 이상이
필요하다.

## 작업 흐름

간단한 설명이나 질의응답은 기록하지 않는다. 기록할 작업은 범위에 따라 나눈다.

| 분류 | 흐름 |
| --- | --- |
| Direct | Director가 직접 처리하고 `REQUEST → RUN_DONE`을 기록한다. |
| Standard | `REQUEST → PLANNED → EXECUTED → REVIEW → CLOSE`로 진행한다. |

Standard 작업의 기본 흐름은 다음과 같다.

1. Director가 요청을 분류한다.
2. Planner가 계획과 성공 기준을 작성한다.
3. Executor가 계획을 구현하고 검증한다.
4. Planner가 구현 결과를 검토한다.
5. Director가 결과를 사용자에게 보고하고 승인을 요청한다.
6. 사용자가 명시적으로 승인하면 Director가 작업을 종료한다.

`REVIEW`는 승인에 필요한 증거이지 승인이 아니다. 사용자 피드백이나 blocker가 있으면 Executor와 Planner 단계를 반복한다.

## 이벤트와 산출물

`relay.log`는 다음 형식의 추가 전용 이벤트 기록이다.

```text
<timestamp> | <task-id> | <event> | <role> | <summary> | <path?>
```

Standard 작업의 산출물은 `.agent-relay/runs/`에 같은 `<KEY>`로 저장한다.

| 파일 | 작성자 | 내용 |
| --- | --- | --- |
| `<KEY>-PLAN.md` | Planner | 계획과 성공 기준 |
| `<KEY>-RUN-<NN>.md` | Executor | 변경과 검증 결과 |
| `<KEY>-REVIEW-<NN>.md` | Planner | 검토 결과 |
| `<KEY>-CLOSE.md` | Director | 사용자 승인 기록 |

이전 라운드는 보존한다. 자세한 이벤트 전이와 파일 규칙은 [프로토콜](bootstrap/.agent-relay/PROTOCOL.md)을 따른다.

## 설치 후 구조

```text
.agent-relay/
├── PROTOCOL.md
├── DIRECTOR.md
├── PLANNER.md
├── EXECUTOR.md
├── GUIDANCE.md
├── LESSON-LEARNED.md
├── bin/
│   └── agent-relay[.exe]
├── lesson-learned/
├── relay.log
├── runs/
└── templates/
```

`.agent-relay/`는 프로젝트 인수인계 데이터이므로 Git에 포함한다.

## 개발과 바이너리 빌드

```text
go test ./...
go run ./cmd/build-release
```

`build-release`는 `CGO_ENABLED=0`으로 지원 플랫폼 6종을 빌드하고
`bootstrap/.agent-relay/bin/SHA256SUMS`를 갱신한다. 바이너리와 체크섬은 함께
커밋하며 부트스트랩과 업데이트에서 무결성을 검증한다.

## 운영 원칙

- Director만 `relay.log`와 작업 상태를 변경한다.
- Executor는 계획 범위를 임의로 넓히지 않는다.
- Standard 작업은 사용자의 명시적 승인 전에는 종료하지 않는다.
- 사용자에게 보고된 결함은 증거를 확보한 뒤 수정하고 스모크 테스트를 남긴다.
- `.agent-relay/`에 비밀정보, 자격증명, 개인정보, 민감한 운영정보를 저장하지 않는다.
- 커밋 전이나 업데이트 후에는 `.agent-relay/bin/agent-relay[.exe] lint`로 기록 상태를 검증한다.
- Agent Relay 자체 변경은 `go test ./...`로 회귀 테스트한다.

## 더 읽기

- [BOOTSTRAP.md](BOOTSTRAP.md) — 설치와 업데이트 절차
- [PROTOCOL-GUIDE.md](PROTOCOL-GUIDE.md) — 한국어 상세 가이드
- [PROTOCOL.md](bootstrap/.agent-relay/PROTOCOL.md) — 공통 규칙
- [DIRECTOR.md](bootstrap/.agent-relay/DIRECTOR.md), [PLANNER.md](bootstrap/.agent-relay/PLANNER.md), [EXECUTOR.md](bootstrap/.agent-relay/EXECUTOR.md) — 역할별 규칙
