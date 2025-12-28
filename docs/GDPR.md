# GDPR Implementation Guide

## Overview

This document describes the iRegistro system's GDPR compliance implementation, covering encryption, data protection, user rights, and legal requirements for Italian schools.

## Architecture

The GDPR implementation consists of:
- **Encryption Service**: AES-256-GCM for sensitive data at rest
- **Audit Logger**: Comprehensive access logging (Art. 30)
- **Compliance Service**: User rights implementation (Art. 12-22)
- **Middleware**: Security headers and TLS enforcement
- **Database Schema**: Consent, audit logs, deletion requests, data exports

## Encryption

### At Rest

Sensitive fields are encrypted using AES-256-GCM:
- Tax codes (Codice Fiscale)
- Phone numbers
- Home addresses

**Implementation**:
```go
encSvc, _ := gdpr.NewEncryptionService()
encrypted, _ := encSvc.EncryptTaxCode("RSSMRA80A01H501U")
```

**Environment Variable**:
```bash
ENCRYPTION_KEY=your-32-byte-encryption-key-here
```

⚠️ **IMPORTANT**: Use a strong, randomly generated 32-byte key and store it securely (e.g., AWS Secrets Manager, HashiCorp Vault).

### In Transit

- **TLS 1.3 minimum** enforced in production
- **HSTS headers** force HTTPS
- **WebSocket Secure (WSS)** only
- Certificate pinning recommended for mobile apps

**Middleware**:
```go
router.Use(middleware.GDPRComplianceMiddleware())
router.Use(middleware.TLSRequiredMiddleware())
```

## User Rights (GDPR Art. 12-22)

### 1. Right to Access (Art. 15)

Users can download all their data.

**Endpoint**: `GET /users/:userId/data?format=json|csv|xml`

**Example**:
```bash
curl -X GET https://api.registro.it/users/123/data?format=json \
  -H "Authorization: Bearer TOKEN"
```

**Response**:
```json
{
  "export_id": "uuid",
  "status": "PENDING",
  "message": "Data export initiated"
}
```

### 2. Right to Rectification (Art. 16)

Users can update their personal data.

**Authorization**:
- Parents/Students: Can modify their own profile data
- Teachers: Cannot modify personal data (only school admin/director)

### 3. Right to Erasure (Art. 17)

Users can request account deletion with 30-day grace period.

**Endpoint**: `POST /users/:userId/data-deletion-request`

**Request**:
```json
{
  "reason": "No longer attending this school"
}
```

**Response**:
```json
{
  "request_id": "uuid",
  "scheduled_deletion_at": "2024-02-28T12:00:00Z",
  "message": "Deletion request submitted. Data will be deleted after 30-day grace period."
}
```

**Exceptions** (data retained by law):
- Academic marks: 5 years (Italian school law)
- Student absences: 5 years
- Official documents: 10 years (Italian archive law)

### 4. Right to Data Portability (Art. 20)

Users receive data in structured, machine-readable format.

**Endpoint**: `POST /users/:userId/data-export`

**Request**:
```json
{
  "format": "JSON"
}
```

**Supported Formats**:
- JSON (default)
- CSV
- XML

**Export Contents**:
- User profile
- Marks (voti)
- Absences (assenze)
- Messages
- Documents
- Colloquium bookings

### 5. Right to Restriction (Art. 18)

Parents can temporarily block data access (implementation pending).

## Consent Management (Art. 7)

### Consent Types

1. **COMMUNICATIONS**: School communications (email, SMS, push notifications)
2. **PHOTOS**: Photos and videos of student
3. **BIOMETRIC**: Biometric data (Face ID, Touch ID)
4. **DATA_PROCESSING**: General data processing consent

### Grant Consent

**Endpoint**: `PUT /users/:userId/consent/:consentType`

**Example**:
```bash
curl -X PUT https://api.registro.it/users/123/consent/PHOTOS \
  -H "Authorization: Bearer TOKEN"
```

### Revoke Consent

**Endpoint**: `DELETE /users/:userId/consent/:consentType`

**Request**:
```json
{
  "reason": "No longer comfortable with photos"
}
```

### View Consents

**Endpoint**: `GET /users/:userId/consent`

**Response**:
```json
[
  {
    "consent_type": "COMMUNICATIONS",
    "granted": true,
    "given_at": "2024-01-15T10:30:00Z",
    "ip_address": "192.168.1.100",
    "revoked_at": null
  }
]
```

## Audit Logging (Art. 30)

All data access is logged:

**Endpoint**: `GET /audit/my-accesses`

**Response**:
```json
{
  "accesses": [
    {
      "accessed_by_user_id": "teacher-uuid",
      "resource_type": "MARKS",
      "accessed_at": "2024-01-20T14:30:00Z",
      "purpose": "Grade entry",
      "action": "WRITE"
    }
  ]
}
```

**Logged Information**:
- Who accessed the data
- What data was accessed
- When
- From which IP address
- Purpose/action

## Data Retention Policies

| Data Type | Retention Period | Legal Basis |
|-----------|------------------|-------------|
| Marks (Voti) | 5 years | Italian School Law |
| Absences (Assenze) | 5 years | Italian School Law |
| Audit Logs | 1 year | GDPR Art. 32 |
| Messages | 2 years | School Policy |
| Temporary Files | 7 days | Data Minimization |
| Official Documents | 10 years | Italian Archive Law |

**Automatic Cleanup**:
A cron job should run daily to delete expired data:
```bash
0 2 * * * /app/bin/cleanup-expired-data
```

## Security Headers

Automatically added by `GDPRComplianceMiddleware`:

```http
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Content-Security-Policy: default-src 'self'
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

## Privacy by Design

1. **Data Minimization**: Collect only necessary data
2. **Pseudonymization**: Reports use anonymized IDs
3. **Encryption**: Sensitive fields encrypted at rest
4. **Access Control**: Users see only their own data + role permissions
5. **Audit Trail**: All accesses logged
6. **Secure by Default**: HTTPS only, TLS 1.3 minimum

## Data Breach Response

**Procedure** (GDPR Art. 33, 34):

1. **Detection**: Monitor audit logs and security alerts
2. **Assessment** (within 24h):
   - Severity
   - Affected users
   - Data types compromised
3. **Notification**:
   - Supervisory authority: within 72 hours
   - Affected users: if high risk
4. **Mitigation**:
   - Revoke compromised credentials
   - Force password reset
   - Audit system access
5. **Documentation**: Log all actions in incident report

**Contact**: privacy@registro.it

## DPIA (Data Protection Impact Assessment)

Required for high-risk processing (Art. 35).

**Database Table**: `privacy_impact_assessments`

**When Required**:
- New biometric features
- Large-scale processing of sensitive data
- Systematic monitoring

**Process**:
1. Assess necessity and proportionality
2. Identify risks
3. Define mitigation measures
4. Seek DPO opinion
5. Document assessment
6. Review regularly

## Testing

**Run GDPR tests**:
```bash
go test ./internal/gdpr/...
```

**Coverage**:
- Encryption/decryption
- Consent management
- Data export (all formats)
- Deletion requests
- Audit logging

## Production Checklist

- [ ] Set `ENCRYPTION_KEY` environment variable (32 bytes)
- [ ] Enable TLS 1.3 on load balancer/reverse proxy
- [ ] Configure HSTS preload
- [ ] Set up certificate pinning for mobile apps
- [ ] Configure daily cleanup cron job
- [ ] Enable database encryption at rest (PostgreSQL)
- [ ] Set up Slack/email alerts for data access violations
- [ ] Train staff on GDPR procedures
- [ ] Appoint Data Protection Officer (DPO)
- [ ] Review privacy policy annually

## Legal References

- **GDPR**: Regulation (EU) 2016/679
- **Italian Implementation**: D.Lgs. 196/2003 (as amended)
- **School Data**: DPR 122/2009 (student assessment)
- **Archive Law**: D.Lgs. 42/2004

## Support

For GDPR compliance questions:
- **DPO**: dpo@registro.it
- **Privacy Team**: privacy@registro.it
- **Legal**: legal@registro.it
