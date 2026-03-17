# 🚀 API Update: Auth & Handshake

Dieses Update ersetzt die Sektion **"Auth"** in der ursprünglichen `api.md`. Da wir **Supabase** nutzen, verschieben wir die Authentifizierungs-Logik vom Go-Backend direkt zum Client/Supabase.

## 1. Authentifizierung (Supabase)
Der Login, Logout und Passwort-Reset werden **nicht** über das Go-Backend abgewickelt.

* **Login:** Erfolgt über das Supabase-SDK (Client-seitig).
* **Token:** Nach dem Login erhält die App ein JWT (`access_token`).
* **Validierung:** Dieses JWT muss bei jedem Request an das Go-Backend im Header mitgeschickt werden.

---

## 2. Neuer Identity-Endpunkt (Ersatz für /auth/login)

Statt `POST /auth/login` nutzt die App diesen Endpunkt, um sicherzustellen, dass die Session auch im Backend gültig ist.

### GET `/api/me`
Gibt die Identität des aktuell über das JWT eingeloggten Users zurück.

**Headers:** `Authorization: Bearer <token>`

**Response `200`:**
```json
{
  "user_id": "abc-123-uuid",
  "status": "authenticated"
}
```

## 3. Entfallende Endpunkte
Folgende Endpunkte aus der ursprünglichen Spezifikation werden **gestrichen**, da sie vollständig vom Supabase-SDK (Client-seitig) übernommen werden:

* ❌ **`POST /auth/login`**: Die App authentifiziert sich direkt bei Supabase. Das Go-Backend erhält danach nur das fertige JWT.
* ❌ **`POST /auth/logout`**: Ein Logout erfolgt durch das Löschen des Tokens im lokalen Speicher der App und einen Aufruf an das Supabase-SDK.
* ❌ **`PATCH /auth/password`**: Passwortänderungen werden über die eingebauten Supabase-Funktionen (z. B. Reset-Mail) abgewickelt.

---

## 4. Business Logik (Backend Responsibility)
Die eigentliche Logik der App bleibt in der Verantwortung des Go-Servers. Alle Anfragen hier müssen das JWT im Header mitschicken.

| Methode | Endpunkt | Beschreibung |
| :--- | :--- | :--- |
| **GET** | `/api/companies` | Liste aller aktiven Unternehmen für das Swiping. |
| **POST** | `/api/companies/{id}/vote` | Speichert `LIKE`, `DISLIKE` oder `NEUTRAL`. |
| **GET** | `/api/appointments/me` | Holt den persönlichen Terminplan (nach dem Matching). |

---

## 5. Security Check: Wer darf was?

| Datentyp | Verwaltung durch | Info |
| :--- | :--- | :--- |
| **User-Identität** | Supabase | Anmeldung, Verifizierung, Passwort-Sicherheit. |
| **Daten-Zugriff** | Go-Backend | Prüft pro Request: "Darf User X die Daten von Firma Y sehen?" |
| **Matching-Logik** | Go-Backend | Berechnet die Termine basierend auf den Votes. |