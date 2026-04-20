### 🛡️ API Access Control List (RBAC)

| Kategorie | Route | Methode | Erlaubte Rollen | Bedingung (Self-Service) |
| :--- | :--- | :--- | :--- | :--- |
| **System** | `/api/seed` | `GET` | Admin | - |
| **User** | `/api/invite` | `POST` | Admin | - |
| **User** | `/api/users/{id}` | `PATCH` | Admin, Student| Nur eigene ID |
| **Auth** | `/api/me` | `GET` | Alle | Authentifizierter User |
| **Public-ish** | `/api/resend-invite` | `POST` | Alle (Public) | Kein Token nötig |
| **Companies** | `/api/companies/active` | `GET` | Student, Admin | - |
| **Companies** | `/api/companies/{id}/vote` | `POST` | Student | - |
| **Companies** | `/api/companies/{id}/logo` | `POST` | Admin, Company | Nur eigene Company-ID |
| **Companies** | `/api/companies/{id}` | `PATCH` | Admin, Company | Nur eigene Company-ID |
| **Students** | `/api/students` | `GET` | Admin | - |
| **Students** | `/api/students/{id}` | `PATCH` | Admin, Student | Nur eigene ID |
| **Students** | `/api/students/{id}` | `DELETE` | Admin, Student | Nur eigene ID |
| **Students** | `/api/students/{id}/preferences` | `GET` | Admin, Student | Nur eigene ID |
| **Students** | `/api/students/{id}/meetings` | `GET` | Admin, Student | Nur eigene ID |
| **Matching** | `/api/meetings/assign` | `POST` | Admin | - |
| **Events** | `/api/events/active` | `GET` | Alle | Authentifizierter User |
| **Slots** | `/api/slots` | `GET` | Alle | Authentifizierter User |
| **Slots** | `/api/slots/{id}` | `DELETE` | Admin | - |

---

### 🔑 Logik-Hinweise für die Implementierung

1.  **Admin-Check:** Wenn `user_role == 'admin'`, wird der Zugriff immer gewährt.
2.  **Self-Service Check:** Bei Routen mit `{id}` muss im Handler geprüft werden, ob die `id` aus der URL mit der `id` (oder `company_id`/`student_id`) aus dem JWT übereinstimmt.
3.  **403 Forbidden:** Wenn ein Student versucht, auf `/api/students` (Liste aller Studenten) zuzugreifen, sollte der Server strikt mit einem **403** antworten.
4.  **Datenintegrität:** Bei `PATCH /companies/{id}` sollte im SQL-Query zusätzlich sichergestellt werden, dass die `user_id` der Firma zum Token passt, um "ID-Guessing" zu verhindern.