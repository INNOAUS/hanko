# 메일 템플릿 ↔ 다국어 로케일 매핑

모든 HTML 메일은 `layout.html.tmpl`을 사용하며, 레이아웃의 헤더/푸터 번역은 **passcode** 로케일에서 가져옵니다.

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
