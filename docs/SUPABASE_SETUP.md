# Supabase Cloud Setup

If you prefer using Supabase Managed Cloud instead of self-hosting.

## 1. Create Project
- Go to [Supabase Dashboard](https://supabase.com/dashboard)
- Create a new project.
- Note the `Reference ID`, `URL`, and `service_role` Key.

## 2. Database
- Navigate to **SQL Editor**.
- Open `schema.sql` from this repository.
- Copy/Paste content and run to create tables.
- Run `009_rls_policies.up.sql` to enable security policies.

## 3. Storage
- Go to **Storage** section.
- Create new private buckets:
    - `documents`
    - `attachments`
    - `pdfs`
- Set Policies (ACL) for buckets:
    - Create policy "Authenticated can view" -> `SELECT` for `auth.role() = 'authenticated'`.
    - Create policy "Uploads" -> `INSERT` for valid users.

## 4. Environment Variables
Update your application `.env`:
```env
DB_HOST=aws-0-eu-central-1.pooler.supabase.com
DB_USER=postgres.[reference_id]
DB_PASSWORD=[your_password]

SUPABASE_URL=https://[reference_id].supabase.co
SUPABASE_KEY=[service_role_key]
```

## 5. Auth
- Since we use **Custom Auth** in Go (JWT), you don't need Supabase Auth enabled for login, but you do need strict RLS.
- Ensure your Go backend generates JWTs signed with the **Supabase JWT Secret** so Supabase (PostgREST) recognizes them if you use Supabase Client features.
- If only using Postgres driver directly from Go, RLS policies must be applied using `SET app.current_user_id = ...` in the transaction (see `009_rls_policies.up.sql`).
