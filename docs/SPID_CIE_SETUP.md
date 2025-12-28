# SPID & CIE Authentication Setup Guide

## Overview

This guide explains how to set up SPID (Sistema Pubblico di Identità Digitale) and CIE (Carta d'Identità Elettronica) authentication for the iRegistro system.

## Prerequisites

### SPID Registration

1. **Register as Service Provider** with AgID (Agenzia per l'Italia Digitale)
   - Visit: https://www.agid.gov.it/it/piattaforme/spid
   - Submit SP metadata
   - Obtain AgID certification

2. **Production Requirements**:
   - Public HTTPS endpoint
   - Valid SSL certificate
   - Registered callback URL: `https://your-domain.com/auth/spid/callback`
   - Service Provider metadata XML

### CIE Registration

1. **Register with Ministero dell'Interno**
   - Visit: https://www.cartaidentita.interno.gov.it/
   - Apply for OIDC client credentials
   - Complete integration testing

2. **Production Requirements**:
   - Client ID and Client Secret
   - Registered redirect URI: `https://your-domain.com/auth/cie/callback`
   - Public HTTPS endpoint

## Environment Configuration

### Required Environment Variables

```bash
# General
FRONTEND_URL=https://your-frontend.com
API_URL=https://api.your-domain.com

# SPID Configuration
SPID_METADATA_URL=https://registry.spid.gov.it/metadata/idp/spid-entities-idps.xml
SPID_ENTITY_ID=https://api.your-domain.com/auth/spid/metadata
SPID_ACS_URL=https://api.your-domain.com/auth/spid/callback
SPID_CERTIFICATE_PATH=/path/to/spid-cert.pem
SPID_PRIVATE_KEY_PATH=/path/to/spid-key.pem

# CIE Configuration
CIE_ISSUER=https://idserver.servizicie.interno.gov.it/idp/profile/oidc
CIE_CLIENT_ID=your-client-id
CIE_CLIENT_SECRET=your-client-secret
CIE_REDIRECT_URI=https://api.your-domain.com/auth/cie/callback
```

### Example `.env` file

```env
# Add to your existing .env file

# SPID
SPID_METADATA_URL=https://registry.spid.gov.it/metadata/idp/spid-entities-idps.xml
SPID_ENTITY_ID=https://registro.example.com/auth/spid/metadata
SPID_ACS_URL=https://registro.example.com/auth/spid/callback
SPID_CERTIFICATE_PATH=./certs/spid-cert.pem
SPID_PRIVATE_KEY_PATH=./certs/spid-key.pem

# CIE
CIE_ISSUER=https://idserver.servizicie.interno.gov.it/idp/profile/oidc
CIE_CLIENT_ID=your_cie_client_id
CIE_CLIENT_SECRET=your_cie_client_secret
CIE_REDIRECT_URI=https://registro.example.com/auth/cie/callback

# Frontend
FRONTEND_URL=https://registro.example.com
```

## SPID Setup Steps

### 1. Generate Certificates

```bash
# Generate private key
openssl genrsa -out spid-key.pem 2048

# Generate certificate signing request
openssl req -new -key spid-key.pem -out spid-csr.pem \
  -subj "/C=IT/O=Your School/CN=registro.example.com"

# Generate self-signed certificate (for testing)
openssl x509 -req -in spid-csr.pem -signkey spid-key.pem \
  -out spid-cert.pem -days 365
```

### 2. Generate Service Provider Metadata

Start the API server and visit:
```
https://api.your-domain.com/auth/spid/metadata
```

Download the XML metadata and submit to AgID.

### 3. SPID Test Environment

For testing, use SPID Validator:
- URL: https://validator.spid.gov.it
- Test credentials provided by validator
- No production registration needed

## CIE Setup Steps

### 1. Register Application

1. Access CIE developer portal
2. Create new application
3. Obtain `CLIENT_ID` and `CLIENT_SECRET`
4. Register redirect URI

### 2. CIE Test Environment

CIE provides test environment with test cards:
- Test Issuer: `https://preproduzione.idserver.servizicie.interno.gov.it`
- Test cards available from Ministero dell'Interno

## Integration

### Backend Routes

The following routes are automatically registered:

**SPID:**
```
GET  /auth/spid/login?school_id=<uuid>&redirect_uri=<url>
POST /auth/spid/callback
GET  /auth/spid/metadata
```

**CIE:**
```
GET /auth/cie/login?school_id=<uuid>&redirect_uri=<url>
GET /auth/cie/callback?code=<code>&state=<state>
```

### Frontend Integration

Import the login components:

```vue
<template>
  <div class="login-page">
    <SPIDLogin :schoolId="schoolId" />
    <CIELogin :schoolId="schoolId" />
    <EmailLogin />
  </div>
</template>

<script setup>
import SPIDLogin from '@/components/auth/SPIDLogin.vue'
import CIELogin from '@/components/auth/CIELogin.vue'
import EmailLogin from '@/components/auth/EmailLogin.vue'

const schoolId = 'your-school-uuid'
</script>
```

## Authentication Flow

### SPID Flow

1. User clicks "Accedi con SPID"
2. Frontend redirects to `/auth/spid/login?school_id=<uuid>`
3. Backend creates SAML request, redirects to SPID aggregator
4. User selects identity provider (Poste, Aruba, etc.)
5. User authenticates with selected provider
6. Provider returns SAML assertion to `/auth/spid/callback`
7. Backend validates assertion, extracts attributes
8. Backend creates/updates user, generates JWT
9. Redirect to frontend with JWT token
10. Frontend stores token, user logged in

### CIE Flow

1. User clicks "Accedi con CIE"
2. Frontend redirects to `/auth/cie/login?school_id=<uuid>`
3. Backend creates OAuth2 request, redirects to CIE
4. User authenticates with CIE (card + PIN or app)
5. CIE returns authorization code to `/auth/cie/callback`
6. Backend exchanges code for ID token
7. Backend validates token, extracts claims
8. Backend creates/updates user, generates JWT
9. Redirect to frontend with JWT token
10. Frontend stores token, user logged in

## Troubleshooting

### SPID Issues

**Problem**: "SAML validation failed"
- Check certificate validity
- Verify metadata is up to date
- Ensure clock synchronization (SAML is time-sensitive)

**Problem**: "Missing required attribute"
- SPID providers must return: `spidCode`, `name`, `familyName`
- Check attribute mapping in provider metadata

### CIE Issues

**Problem**: "Token validation failed"
- Verify `CIE_CLIENT_ID` and `CIE_CLIENT_SECRET`
- Check ID token signature
- Ensure HTTPS is enabled

**Problem**: "State mismatch"
- CSRF protection - ensure cookies are enabled
- Check domain matches for cookie
- Verify HTTPS (secure cookies required)

## Security Considerations

1. **HTTPS Required**: Both SPID and CIE require HTTPS in production
2. **Certificate Management**: Rotate certificates annually
3. **Token Security**: JWT tokens in httpOnly cookies
4. **CSRF Protection**: State parameter validated for CIE
5. **Audit Logging**: All authentication attempts logged in `auth_audit_logs`

## Testing

### Unit Tests

```bash
go test ./internal/application/auth/...
go test ./internal/presentation/http/handlers/...
```

### Integration Tests

```bash
# Requires test SPID/CIE setup
go test ./tests/integration/spid_test.go
go test ./tests/integration/cie_test.go
```

### Manual Testing

1. **SPID Validator**: https://validator.spid.gov.it
2. **CIE Test Environment**: Use preprod issuer
3. **Check logs**: `auth_audit_logs` table

## Production Checklist

- [ ] SPID metadata submitted to AgID
- [ ] SPID certificates installed and valid
- [ ] CIE client registered with Ministero
- [ ] Environment variables configured
- [ ] HTTPS enabled with valid certificate
- [ ] Callback URLs whitelisted
- [ ] Frontend integration tested
- [ ] Audit logging verified
- [ ] Error handling tested
- [ ] Documentation updated

## Support

- **SPID**: support@agid.gov.it
- **CIE**: assistenzacie@interno.it
- **Technical**: Consult developer portals for each service

## References

- [SPID Technical Rules](https://docs.italia.it/italia/spid/spid-regole-tecniche)
- [CIE OIDC Documentation](https://www.cartaidentita.interno.gov.it/identificazione-digitale/entra-con-cie/)
- [AgID Developer Portal](https://developers.italia.it/)
