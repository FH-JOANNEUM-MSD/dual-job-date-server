### 🛡️ API Access Control List (RBAC)

| Kategorie | Route | Methode | Erlaubte Rollen | Bedingung (Self-Service) |
| :--- | :--- | :--- | :--- | :--- |
| **System** | `/` | `GET` | Alle (Public) | Health-Check; kein Token nötig |
| **System** | `/api/seed` | `GET` | Admin | - |
| **User** | `/api/invite` | `POST` | Admin | - |
| **User** | `/api/users/{id}` | `PATCH` | Admin, Student, Company | Nur eigene ID |
| **Auth** | `/api/me` | `GET` | Alle | Authentifizierter User |
| **Public-ish** | `/api/resend-invite` | `POST` | Alle (Public) | Kein Token nötig |
| **Companies** | `/api/companies` | `GET` | Admin | Alle Einträge inkl. inaktiver Firmen |
| **Companies** | `/api/companies/active` | `GET` | Admin, Student | Nur aktive Firmen |
| **Companies** | `/api/companies/{id}/vote` | `POST` | Student | - |
| **Companies** | `/api/companies/{id}/logo` | `POST` | Admin, Company | Nur eigene Company-ID |
| **Companies** | `/api/companies/{id}/images`| `POST` | Admin, Company | Nur eigene Company-ID |
| **Companies** | `/api/companies/{id}/meetings` | `GET` | Admin, Company | Nur eigene Company-ID; enthält Slot-Zeiten + Studentenname |
| **Companies** | `/api/companies/{id}` | `PATCH` | Admin, Company | Nur eigene Company-ID |
| **Companies** | `/api/companies/{id}` | `GET` | Admin, Company | Nur eigene Company-ID |
| **Students** | `/api/students` | `GET` | Admin | - |
| **Students** | `/api/students/{id}` | `PATCH` | Admin, Student | Nur eigene ID |
| **Students** | `/api/students/{id}` | `DELETE` | Admin, Student | Nur eigene ID |
| **Students** | `/api/students/{id}/preferences` | `GET` | Admin, Student | Nur eigene ID |
| **Students** | `/api/students/{id}/meetings` | `GET` | Admin, Student | Nur eigene ID |
| **Matching** | `/api/meetings/assign` | `POST` | Admin | Event-bezogen: `event_id` (Pflicht); optional `slot_ids`, `student_ids`, `dry_run`, `replace_existing`, `include_inactive_companies` |
| **Matching** | `/api/events/{id}/meetings` | `PUT` | Admin | Setzt/synchronisiert den kompletten Zeitplan eines Events; Body: `{ meetings: [{ slot_id, student_id, company_id }, ...] }` |
| **Matching** | `/api/meetings/{id}` | `PATCH` | Admin | `slot_id`, `student_id`, `company_id` (mind. eines) |
| **Matching** | `/api/allMeetings` | `GET` | Admin | Alle Meetings inkl. Slot-Zeiten + Studentenname |
| **Matching** | `/api/allMeetings/{id}` | `GET` | Admin, Company | Nur eigene Company-ID; enthält Slot-Zeiten + Studentenname |
| **Matching** | `/api/allPrefences` | `GET` | Admin | Alle Präferenzen aller Studierenden |
| **Events** | `/api/events` | `GET` | Admin | Alle Einträge inkl. inaktiver Events |
| **Events** | `/api/events` | `POST` | Admin | - |
| **Events** | `/api/events/{id}` | `PATCH` | Admin | Partial Update (`name`, `location`, `description`, `event_date`, `is_active`) |
| **Events** | `/api/events/{id}` | `DELETE` | Admin | - |
| **Events** | `/api/events/active` | `GET` | Alle | Authentifizierter User |
| **Slots** | `/api/slots` | `GET` | Alle | Authentifizierter User; `?event_id=<id>` filtert auf die Slots eines Events |
| **Slots** | `/api/slots` | `POST` | Admin | `event_id` (Pflicht); Slots sind event-eigen |
| **Slots** | `/api/slots/{id}` | `PATCH` | Admin | Partial Update (`start_time`, `end_time`) |
| **Slots** | `/api/slots/{id}` | `DELETE` | Admin | - |

---

### 🔑 Logik-Hinweise für die Implementierung

1.  **Admin-Check:** Wenn die Rolle im Token-Context `admin` ist, wird der Zugriff bei allen geschützten Routen gewährt (Dank `RequireRole` und `RequireSelfOrAdmin`).
2.  **Self-Service Check:** Bei Routen mit `{id}` prüft der Wrapper `RequireSelfOrAdmin()`, ob die `{id}` aus der URL mit der ID des eingeloggten Users aus dem JWT übereinstimmt. Ist das der Fall, darf der User seine eigenen Daten bearbeiten.
3.  **403 Forbidden vs. 401 Unauthorized:** * `401 Unauthorized`: Der Token fehlt, ist ungültig oder der User existiert nicht mehr in der DB (geprüft durch `JWTMiddleware`).
    * `403 Forbidden`: Der User ist eingeloggt, hat aber nicht die nötige Rolle für diesen Endpunkt (z. B. Student ruft `/api/students` auf) oder versucht, fremde IDs zu bearbeiten (geprüft durch die Wrapper).
4.  **Datenintegrität:** Bei Updates (PATCH/POST) auf `/companies/{id}` oder `/students/{id}` stellt der Wrapper sicher, dass kein User per "ID-Guessing" die Daten eines anderen überschreibt.
5.  **Event-bezogenes Matching:** `POST /api/meetings/assign` matcht innerhalb der Slots des angegebenen `event_id`; sind `student_ids` gesetzt, werden nur diese Studierenden berücksichtigt. Erzeugte Meetings werden mit der `event_id` markiert. `replace_existing` löscht ausschließlich die Meetings dieses Events. Schutz: Bei einem echten Commit (`dry_run: false`) darf `replace_existing` NICHT mit einer Teilmenge aus `slot_ids`/`student_ids` kombiniert werden (-> `400`); in der Vorschau (`dry_run: true`) ist diese Kombination erlaubt.
6.  **Zeitplan-Synchronisation:** `PUT /api/events/{id}/meetings` setzt den kompletten Zeitplan eines Events: weggefallene Meetings werden entfernt, neue hinzugefügt; zurückgegeben werden die angereicherten Meetings des Events. Wird vom Web zum Speichern der Termin-Matrix genutzt.
7.  **`event_id` & Kaskaden:** Slots und Meetings (sowie die DTOs `MeetingDetail`/`CompanyMeeting`) tragen eine `event_id`; die Lese-Endpunkte liefern sie mit. Die `event_id`-Fremdschlüssel sind `ON DELETE CASCADE` – das Löschen eines Events entfernt dessen Slots und Meetings. Kanonisches Schema: `DocumentationBackendServer/schema.sql`.