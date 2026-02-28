# 메일 템플릿 ↔ 다국어 로케일 매핑

모든 HTML 메일은 `layout.html.tmpl`을 사용하며, 레이아웃의 헤더/푸터 번역은 **passcode** 로케일에서 가져옵니다.

## 설정으로 덮어쓰기 (mail_template.yaml, 재빌드 불필요)

**config.yaml과 독립된** `mail_template.yaml`에서 **제품명·푸터·저작권**과 **본문·제목** 메시지를 로케일별로 설정할 수 있습니다.

- `config.yaml`의 `service.mail_template_file`에 경로 지정 (예: `config/mail_template.yaml`, 프로세스 CWD 기준).
- 비우면 passcode.*.yaml, security-notifications.*.yaml 로케일 번역 사용.
- 파일 형식: 최상위 키는 언어 코드(en, ko, ja, zh). 각 언어 아래에 `product_name`, `footer_sent_by`, `copyright` 및 본문/제목용 키(예: `login_text`, `ttl_text`, `mfa_create_text`, `subject_mfa_create` 등) 정의. `{{ .ServiceName }}`, `{{ .Code }}` 등 치환 가능.

예: `backend/config/mail_template.yaml` 참고.

## 템플릿별 번역 키

| 템플릿 | 로케일 파일 | 본문(바디) 키 | 제목 키 (Go에서 사용) |
|--------|--------------|----------------|------------------------|
| layout.html.tmpl | passcode | mail_product_name, mail_footer_sent_by, mail_copyright | — |
| login | passcode | login_text, ttl_text | email_subject_login |
| recovery | passcode | recovery_text, ttl_text | subject_recovery |
| email_verification | passcode | email_verification_text, ttl_text | subject_email_verification |
| email_login_attempted | passcode | email_login_attempted_text | subject_email_login_attempted |
| email_registration_attempted | passcode | email_registration_attempted_text | subject_email_registration_attempted |
| mfa_create | security-notifications | mfa_create_text | subject_mfa_create |
| mfa_delete | security-notifications | mfa_delete_text | subject_mfa_delete |
| email_create | security-notifications | email_create_text | subject_email_create |
| email_delete | security-notifications | email_delete_text | subject_email_delete |
| passkey_create | security-notifications | passkey_create_text | subject_passkey_create |
| password_update | security-notifications | password_update_text | subject_password_update |
| primary_email_update | security-notifications | primary_email_update_text | subject_primary_email_update |

## 지원 로케일 (ezauth 기준)

- **passcode**: `passcode.en.yaml`, `passcode.ko.yaml`, `passcode.ja.yaml`, `passcode.zh-CN.yaml`, `passcode.zh.yaml`  
  → en, ko, ja, zh-CN, zh 모두 동일 키 세트 필요 (로그인/복구/인증/레이아웃).
- **security-notifications**: `security-notifications.en.yaml`, `security-notifications.ko.yaml`, `security-notifications.ja.yaml`, `security-notifications.zh-CN.yaml`, `security-notifications.zh.yaml`  
  → en, ko, ja, zh-CN, zh 모두 동일 키 세트 필요.

새 언어를 추가할 때는 위 두 로케일 파일 계열에 모두 해당 언어 파일을 추가하고, 위 표에 나온 키를 빠짐없이 정의해야 합니다.

## 다국어·제품명 적용 검증 체크리스트

- **언어 결정**: 메일 언어는 요청의 `Accept-Language`(우선) → `X-Language` → `config.service.default_mail_locale` → `en` 순으로 결정됨.
- **프론트엔드**: ezauth-admin에서 flow API 호출 시 **반드시** `Accept-Language: <locale>` 헤더 전달 필요.  
  - 설정 > 보안(OTP 활성화/비활성화, OTP 코드 검증), 로그인, 비밀번호 찾기, 회원가입 등 모든 flow 요청에 `locale` 전달 확인.
- **제품명(ServiceName)**:  
  - `mail_template.yaml`이 로드되면 **제목·본문** 모두에서 `{{ .ServiceName }}`이 해당 로케일의 `product_name`으로 치환됨.  
  - 적용 경로: Security 알림(제목/본문), 패스코드 핸들러(로그인/복구 메일), flow API 패스코드·알림 메일.
- **지원 로케일**: en, ko, ja, zh. `zh`는 번들에서 `zh-CN`으로 fallback.
- **mail_template.yaml**: `config.service.mail_template_file`에 경로 지정. 비우면 로케일 YAML 번들만 사용(제품명은 config.service.name).
