# 👋 Willkommen bei Dual Job Dating!

Schön, dass du da bist! Dieses Dokument soll dir einen schnellen Überblick geben, worum es in diesem Projekt geht und vor allem, wo du was findest.

## 🚀 Worum geht's hier?

> **Kurz gesagt:** Dies ist das zentrale Backend unseres Projekts. Es ist ein performanter Server, geschrieben in **Go (Golang)**, der sämtliche APIs bereitstellt, die von unseren Frontends / Apps / externen Clients konsumiert werden. Er kümmert sich um die Business-Logik, die Datenbankkommunikation und die Datenbereitstellung.

---

## 🗺️ Projektstruktur: Wo finde ich was?

Damit du dich nicht im Ordner-Dschungel verirrst, hier eine kleine Übersicht der wichtigsten Dateien und Verzeichnisse.
```text
dual-job-date-server/
├─ .env                              # Lokale Umgebungsvariablen
├─ .gitignore                        # Git-Ignore-Regeln
├─ Dockerfile                        # Container-Build für den Server
├─ Makefile                          # Build/Test/Run Kurzbefehle
├─ go.mod                            # Go-Modul + direkte Dependencies
├─ go.sum                            # Checksums der Go-Abhängigkeiten
├─ README.md                         # Hauptdoku / Projektübersicht
├─ README_TEST.md                    # Test-Dokumentation
├─ coverage.html                     # HTML Test-Coverage Report
│
├─ cmd/                              # Entry Points (ausführbare Programme)
│  └─ server/
│     ├─ main.go                     # API-Server Startpunkt
│     └─ main_test.go                # Tests für Server-Entry
│
├─ DocumentationBackendServer/       # API-Dokumentation (fachlich)
│  ├─ API_Access_Control_List.md     # ACL / Rollen / Rechte
│  ├─ Authentication.md              # Auth-Flows
│  ├─ Invite.md                      # Invite-Prozess
│  ├─ Login.md                       # Login-Doku
│  ├─ MagicLinkInvite.md             # Magic-Link Einladung
│  ├─ Server-Cheatsheet.md           # Schnellreferenz
│  ├─ Update.md                      # Update-Endpunkte/-Flows
│  └─ welcome.md                     # Einstieg in die Doku
│
├─ internal/                         # Interne App-Logik (nicht public API)
│  ├─ auth/                          # Auth-Kontext, Middleware, Permissions
│  │  ├─ context.go
│  │  ├─ context_test.go
│  │  ├─ keys.go
│  │  ├─ middleware.go
│  │  ├─ permissions.go
│  │  └─ permissions_test.go
│  │
│  ├─ database/                      # DB-Anbindung/Tests
│  │  ├─ superbase.go                # Supabase-DB Zugriff
│  │  └─ database_test.go
│  │
│  ├─ handlers/                      # HTTP-Handler je Endpoint/Feature
│  │  ├─ auth handler.go
│  │  ├─ comany_handler.go
│  │  ├─ companies_handler.go
│  │  ├─ company_images_handler.go
│  │  ├─ company_logo_handler.go
│  │  ├─ event_handler.go
│  │  ├─ invite_handler.go
│  │  ├─ me_handler_test.go
│  │  ├─ meeting_assignment_handler.go
│  │  ├─ meeting_handler.go
│  │  ├─ mock_handlers_test.go
│  │  ├─ prefernces_handler.go
│  │  ├─ reinvite_handler.go
│  │  ├─ seed_handler.go
│  │  ├─ slot_handler.go
│  │  ├─ slots_delete_handler.go
│  │  ├─ student_handler.go
│  │  ├─ student_update_handler.go
│  │  ├─ students_delete_handler.go
│  │  ├─ updateCompany_handler.go
│  │  ├─ userUpdate_handler.go
│  │  └─ validation_test.go
│  │
│  ├─ models/                        # Domain-Modelle/DTOs
│  │  ├─ UpdateCompany.go
│  │  ├─ company.go
│  │  ├─ event.go
│  │  ├─ inviteRequest.go
│  │  ├─ meeting.go
│  │  ├─ models_test.go
│  │  ├─ prefernces.go
│  │  ├─ reinvite.go
│  │  ├─ slot.go
│  │  ├─ student.go
│  │  ├─ student_update.go
│  │  ├─ user.go
│  │  ├─ userAuthContext.go
│  │  └─ user_update.go
│  │
│  ├─ repository/                    # DB-Zugriff pro Handler
│  │  ├─ auth_invite.go
│  │  ├─ auth_reinvite.go
│  │  ├─ auth_repo.go
│  │  ├─ company_invite.go
│  │  ├─ company_logo_helpers_test.go
│  │  ├─ company_logo_repo.go
│  │  ├─ company_repo.go
│  │  ├─ event_repo.go
│  │  ├─ meeting_assignment_repo.go
│  │  ├─ meeting_repo.go
│  │  ├─ prefernces_slot.go
│  │  ├─ singleCompany_repo.go
│  │  ├─ slot_repo.go
│  │  ├─ slots_delete_repo.go
│  │  ├─ student_delete_repo.go
│  │  ├─ student_invite.go
│  │  ├─ student_repo.go
│  │  ├─ updateCompany_repo.go
│  │  └─ userUpdate_repo.go
│  │
│  └─ routes/                        # Routing + Route-Tests
│     ├─ routes.go
│     └─ routes_test.go
│
├─ tests/                            # Integration/E2E Tests
│  └─ e2e/
│     ├─ .env.test                   # Test-Umgebungsvariablen
│     ├─ main.go                     # E2E Test-Runner Setup
│     ├─ e2e_test.go                 # E2E Testfälle
│     └─ reports/                    # Gespeicherte Testreports
│
├─ coverage/                         # Raw Coverage-Artefakte (ausgeblendet)
└─ .git/                             # Git-Historie/Objekte (ausgeblendet)
```
---

## 📚 Weiterführende Dokumentation

> 🔦 **Interaktive API-Referenz:** [Stoplight Doku](https://jobdatingbackend.stoplight.io/docs/dualjobdating/a54e0e5192a6d-dual-job-dating) — Alle Endpunkte direkt im Browser durchstöbern und ausprobieren.

- **API Access Control List:** [API_Access_Control_List.md](./API_Access_Control_List.md) — Übersicht der API-Routen, Rollen und RBAC-Logik (wer darf welche Endpunkte nutzen; Self-Service-Regeln).
- **Authentication:** [Authentication.md](./Authentication.md) — JWT- & Supabase-Auth-Guide für Frontend/Mobile (Login-Flow, Authorization-Header, Troubleshooting und Beispiele).
- **Invite:** [Invite.md](./Invite.md) — Ablauf für Einladungen und initiales Passwort-Setup; Redirects und Deep-Link-Handling für Web/App.
- **Login:** [Login.md](./Login.md) — Beschreibung des aktualisierten Auth-/Handshake-Workflows und des `/api/me` Endpunkts; was vom Client (Supabase) vs. Backend gehandhabt wird.
- **MagicLinkInvite:** [MagicLinkInvite.md](./MagicLinkInvite.md) — Architektur und Lösungsvorschläge für Deferred Deep Linking (Branch.io, Fallbacks, App-/Web-Handling).
- **Server-Cheatsheet:** [Server-Cheatsheet.md](./Server-Cheatsheet.md) — Kurzanleitung für Build & Redeploy (Docker, kubectl) des Servers auf dem Produktionssystem.
- **E2E Testing:** [E2E.md](./E2E.md) — Start des E2E-Runners, benötigte Umgebungsvariablen und welche Endpunkte bzw. Rechte geprüft werden.
- **Update:** [Update.md](./Update.md) — Guidelines für Partial-Updates (PATCH) der API und Beispiele für korrekte Payloads.

## 🛠️ Schneller Start (Setup)

Willst du das Projekt lokal bei dir ausführen? So geht's:

1. **Klone das Repository:**
   ```bash
   # Mit SSH (Empfohlen)
   git clone git@github.com:FH-JOANNEUM-MSD/dual-job-date-server.git
   # Mit HTTPS
   git clone https://github.com/FH-JOANNEUM-MSD/dual-job-date-server.git
   ```

2. **Installiere die Go-Dependencies:**
   ```bash
   go mod download
   ```

3. **Starte den Server:**
   ```bash
   go run ./cmd/server
   ```

4. **Umgebungsvariablen einrichten:**
   Stelle sicher, dass du eine `.env`-Datei im Hauptverzeichnis hast.
   (Um diese zu erhalten, wende dich bitte an [Andreas Öffl](mailto:andreas.oeffl@fh-joanneum.at).)


## 🙋‍♂️ Noch Fragen?

Melde dich bei [Andreas Öffl](mailto:andreas.oeffl@fh-joanneum.at)

Viel Spaß beim Coden! 💻