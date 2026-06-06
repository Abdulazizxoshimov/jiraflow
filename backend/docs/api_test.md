| Holat | Method | URL | Izoh |
|-------|--------|-----|------|
| ✅ 200 | GET | /auth/me | Auth: me |
| ✅ 200 | POST | /auth/refresh | Auth: refresh |
| ✅ 200 | GET | /users | Users: list |
| ✅ 200 | GET | /users/a0000000-0000-0000-0000-000000000001 | Users: get |
| ❌ --- | POST | /users | Users: create — {"code":"CONFLICT","message":"email already exists"} |
| ✅ 200 | PUT | /users/a0000000-0000-0000-0000-000000000001 | Users: update |
| ✅ 204 | PUT | /users/a0000000-0000-0000-0000-000000000001/password | Users: change password |
| ✅ 201 | POST | /invites | Invites: create → 9550f6fe-e875-44f9-81fc-3217ed60c716 |
| ✅ 200 | GET | /invites | Invites: list |
| ✅ 204 | DELETE | /invites/9550f6fe-e875-44f9-81fc-3217ed60c716 | Invites: revoke |
| ❌ --- | POST | /workflows | Workflows: create — {"code":"INTERNAL_ERROR","message":"ERROR: invalid input syntax for type uuid: \"\" (SQLSTATE 22P02)"} |
| ✅ 200 | GET | /workflows | Workflows: list |
| ❌ --- | POST | /projects | Projects: create — {"code":"NOT_FOUND","message":"project.Create get default workflow: default workflow not found"} |
| ✅ 200 | GET | /projects | Projects: list |
| ✅ 200 | GET | /issues | Issues: list |
| ❌ --- | POST | /spaces | Spaces: create — {"code":"INTERNAL_ERROR","message":"spaceRepo.Create: ERROR: new row for relation \"spaces\" violates check constraint \"spaces_type_check\" (SQLSTATE 23514)"} |
| ✅ 200 | GET | /spaces | Spaces: list |
| ✅ 200 | GET | /pages | Pages: list |
| ✅ 200 | GET | /notifications | Notifications: list |
| ✅ 200 | GET | /notifications/unread-count | Notifications: unread count |
| ✅ 200 | GET | /notifications/preferences | Notifications: get prefs |
| ✅ 200 | PUT | /notifications/preferences | Notifications: update prefs |
| ✅ 204 | POST | /notifications/mark-all-read | Notifications: mark all read |
| ✅ 200 | GET | /search?q=test | Search: q=test |
| ✅ 200 | GET | /audit-logs | Audit logs |
| ✅ 200 | GET | /files/presign?path=test/image.png | Files: presign |

---
**Test tugadi: Sat May 23 01:10:39 PM +05 2026**
