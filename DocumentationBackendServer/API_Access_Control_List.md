### 🛡️ API Access Control List (RBAC)

| Kategorie | Route | Methode | Erlaubte Rollen | Bedingung (Self-Service) |
| :--- | :--- | :--- | :--- | :--- |
| **System** | `/api/seed` | `GET` | Admin | - |
| **User** | `/api/invite` | `POST` | Admin | - |
| **User** | `/api/users/{id}` | `PATCH` | Admin, Student, Company | Nur eigene ID |
| **Auth** | `/api/me` | `GET` | Alle | Authentifizierter User |
| **Public-ish** | `/api/resend-invite` | `POST` | Alle (Public) | Kein Token nötig |
| **Companies** | `/api/companies/active` | `GET` | Admin, Student | - |
| **Companies** | `/api/companies/{id}/vote` | `POST` | Student | - |
| **Companies** | `/api/companies/{id}/logo` | `POST` | Admin, Company | Nur eigene Company-ID |
| **Companies** | `/api/companies/{id}/images`| `POST` | Admin, Company | Nur eigene Company-ID |
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

1.  **Admin-Check:** Wenn die Rolle im Token-Context `admin` ist, wird der Zugriff bei allen geschützten Routen gewährt (Dank `RequireRole` und `RequireSelfOrAdmin`).
2.  **Self-Service Check:** Bei Routen mit `{id}` prüft der Wrapper `RequireSelfOrAdmin()`, ob die `{id}` aus der URL mit der ID des eingeloggten Users aus dem JWT übereinstimmt. Ist das der Fall, darf der User seine eigenen Daten bearbeiten.
3.  **403 Forbidden vs. 401 Unauthorized:** * `401 Unauthorized`: Der Token fehlt, ist ungültig oder der User existiert nicht mehr in der DB (geprüft durch `JWTMiddleware`).
    * `403 Forbidden`: Der User ist eingeloggt, hat aber nicht die nötige Rolle für diesen Endpunkt (z. B. Student ruft `/api/students` auf) oder versucht, fremde IDs zu bearbeiten (geprüft durch die Wrapper).
4.  **Datenintegrität:** Bei Updates (PATCH/POST) auf `/companies/{id}` oder `/students/{id}` stellt der Wrapper sicher, dass kein User per "ID-Guessing" die Daten eines anderen überschreibt.