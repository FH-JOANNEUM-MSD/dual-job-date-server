# E2E Testing

Diese Datei erklärt die End-to-End-Tests des Projekts und den manuellen E2E-Runner in `tests/e2e`.

## Was wird getestet?

Es gibt zwei E2E-nahe Testpfade:

1. Der Go-Test in [tests/e2e/e2e_test.go](../tests/e2e/e2e_test.go) prüft den Health-Endpoint `/` gegen eine laufende Instanz.
2. Der manuelle Runner in [tests/e2e/main.go](../tests/e2e/main.go) fährt mehrere API-Endpunkte gegen eine echte Server-Instanz und schreibt einen Markdown-Report nach `tests/e2e/reports/`.

Der Runner prüft vor allem:

- ob die Server- und Rollen-Logik auf echten HTTP-Aufrufen reagiert
- ob geschützte Endpunkte mit den erwarteten Rollen antworten
- ob zentrale Flows wie Invite, Company-, Student- und Slot-Endpoints erreichbar sind
- ob `/api/me` eine echte Datenbank-ID liefert, damit Folge-Requests mit einer realen User-ID laufen können

## Voraussetzungen

Für den manuellen Runner brauchst du eine laufende Server-Instanz und passende Testdaten in `.env.test`.

Typische Variablen sind:

- `API_URL`
- `SUPABASE_URL`
- `SUPABASE_ANON_KEY`
- `ADMIN_EMAIL`
- `ADMIN_PASS`
- `STUDENT_EMAIL`
- `STUDENT_PASS`
- `COMPANY_EMAIL`
- `COMPANY_PASS`
- `TEST_STUDENT_ID`
- `OTHER_STUDENT_ID`
- `TEST_COMPANY_ID`
- `OTHER_USER_ID` optional, sonst wird ein Fallback verwendet

## Den Runner starten

Aus dem Repository-Root:

```bash
go run ./tests/e2e/main.go
```

Alternativ kannst du auch direkt im Ordner starten:

```bash
cd tests/e2e
go run .
```

Der Runner:

- lädt `.env.test`
- meldet sich als Admin, Student und Company-User an
- ruft `/api/me` auf, um die echte DB-ID des Student-Users zu holen
- führt dann eine Liste von GET, POST, PATCH und DELETE Requests gegen die API aus
- erzeugt danach einen Report unter `tests/e2e/reports/report_<timestamp>.md`

## Was wird im Detail geprüft?

Der manuelle Runner testet unter anderem:

- die Root-Route `/`
- Invite- und Resend-Invite-Endpunkte
- Meeting-Zuordnung
- Company-Listen, Votes, Logo- und Bild-Endpunkte
- Student-Endpunkte wie Profile, Preferences, Meetings und Updates
- User-Updates über eine echte User-ID
- Event- und Slot-Endpunkte
- einige DELETE-Fälle, um Sperr- und Löschverhalten zu prüfen

Wichtig ist dabei nicht nur der Statuscode, sondern auch die Rollenprüfung. Der Runner ruft viele Endpunkte mit mehreren Rollen auf und dokumentiert, ob die API den Zugriff erlaubt, blockiert oder einen Treffer wie `404` zurückgibt.

## Der kleine Go-Test

[tests/e2e/e2e_test.go](../tests/e2e/e2e_test.go) ist ein kleiner Health-Test. Er ist nur aktiv, wenn `E2E_BASE_URL` gesetzt ist.

Beispiel:

```bash
export E2E_BASE_URL="http://localhost:8080"
go test ./tests/e2e -v
```

Dieser Test erwartet auf `/` den Text `Server läuft!`.

## Kurz gesagt

- `go test ./tests/e2e` prüft den einfachen Health-Check
- `go run ./tests/e2e/main.go` führt den größeren manuellen API-Run aus
- beide sind auf echte Server-Antworten ausgelegt, nicht auf Mock-Daten
