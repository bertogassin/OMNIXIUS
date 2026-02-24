# OMNIXIUS — Privacy & PII

Политика применима с первого дня работы платформы. Контакт для запросов — на сайте (страница Contact).

## What we store (PII)

- **Account:** email, name, password hash (Argon2id), optional avatar.
- **Activity:** products you list, orders (buyer/seller), conversations and messages you send.

Data is stored to run the platform (auth, marketplace, internal mail). We do not sell or share PII with third parties for marketing.

## Retention

- **Account data:** Kept until you delete your account.
- **Messages:** Stored as long as the conversation exists; deleted when you delete your account (your messages are removed).
- **Orders:** Order records referencing you as buyer/seller are removed when you delete your account.

No automated retention period beyond that; deletion is on request via account deletion.

## Your rights (GDPR-style)

- **Access:** Use the API (e.g. GET `/api/users/me`, your orders, your conversations) to see your data.
- **Correction:** PATCH `/api/users/me` to update name; use profile/avatar endpoints as documented.
- **Deletion:** DELETE `/api/users/me` (with valid auth) to permanently delete your account and associated PII (see above). This removes your user row, your products, your orders, your conversation participations and messages you sent.

## Contact

For privacy requests or questions, use the contact channel listed on the site (e.g. contact page or support email).
