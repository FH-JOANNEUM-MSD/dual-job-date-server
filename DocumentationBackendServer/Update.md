## API Dokumentation: Ressourcen Aktualisieren (Partial Updates)

Unser Backend nutzt für Updates das **Partial Update Pattern** über die HTTP-Methode `PATCH`.

**Das wichtigste Prinzip:** Bitte schickt in euren Requests **immer nur die Felder mit, die der User auch wirklich geändert hat** (Delta-Update). Schickt **nicht** das komplette Objekt zurück. Das Backend erkennt automatisch, welche Felder im JSON fehlen und lässt diese in der Datenbank unangetastet.

Dadurch sparen wir Bandbreite, verhindern das Überschreiben von Daten durch gleichzeitige User (Race Conditions) und machen das API-Handling sicherer.

---

### Endpunkt: Firma aktualisieren (Beispiel)

Aktualisiert spezifische Felder einer bestehenden Firma anhand ihrer ID.

* **URL:** `/api/companies/{id}`
* **Methode:** `PATCH`
* **Header:** `Content-Type: application/json`

### 📦 Erlaubte Felder (Payload)
Alle Felder sind **optional**. Schickt nur das mit, was aktualisiert werden soll. Die `id`, `user_id` und `last_updated` werden vom Backend/der Datenbank verwaltet und dürfen nicht im Body mitgeschickt werden.

| Feldname | Typ | Beschreibung |
| :--- | :--- | :--- |
| `name` | String | Der Name der Firma |
| `description` | String | Beschreibungstext der Firma |
| `website` | String | URL zur Webseite der Firma |
| `logo_url` | String | URL zum hochgeladenen Logo |
| `active` | Boolean | Status, ob die Firma sichtbar/aktiv ist |

---

### 💡 Beispiele

**Szenario 1: Der User ändert nur den Firmennamen**
Der Request sollte nur das Feld `name` enthalten. Das Logo und die Beschreibung bleiben in der DB exakt so, wie sie waren.

**Request:**
```http
PATCH /api/companies/42
Content-Type: application/json

{
  "name": "TechNova Solutions GmbH"
}
```

Szenario 2: Der User deaktiviert die Firma und ändert die Website
Request:

PATCH /api/companies/42
Content-Type: application/json

```http
{
"active": false,
"website": "[https://technova.example.com](https://technova.example.com)"
}
```

### 📥 Responses (Antworten vom Server)
**✅ 200 OK (Erfolgreich)**
Wird zurückgegeben, wenn das Update erfolgreich in der Datenbank gespeichert wurde.
```http
{
  "message": "Firma erfolgreich aktualisiert",
  "status": "success"
}
```

**❌ 400 Bad Request**
Wird geworfen, wenn das JSON fehlerhaft ist, falsche Datentypen geschickt wurden oder die ID in der URL keine gültige Zahl ist.

**❌ 500 Internal Server Error**
Wird geworfen, wenn es ein Problem bei der Kommunikation mit der Datenbank (Supabase) gab.

### Wichtige Hinweise für das Frontend

Keine Leerstrings für unveränderte Felder! Wenn ihr ein Feld nicht ändern wollt, 
lasst den Key im JSON komplett weg. Wenn ihr z.B. "description": "" 
schickt, wird das Backend die Beschreibung in der Datenbank löschen!

Zeitstempel: Das Feld last_updated wird vom Backend bei jedem PATCH-Request 
automatisch auf die aktuelle Serverzeit gesetzt. 
Ihr müsst euch darum nicht kümmern.

---

### Endpunkt: Event aktualisieren

Aktualisiert spezifische Felder eines bestehenden Events anhand seiner ID. Nur für **Admins**.

* **URL:** `/api/events/{id}`
* **Methode:** `PATCH`
* **Header:** `Content-Type: application/json`

### 📦 Erlaubte Felder (Payload)
Alle Felder sind **optional**. Schickt nur das mit, was aktualisiert werden soll. Die `id` wird von der Datenbank verwaltet und darf nicht im Body mitgeschickt werden.

| Feldname | Typ | Beschreibung |
| :--- | :--- | :--- |
| `name` | String | Name des Events |
| `location` | String | Veranstaltungsort |
| `description` | String | Beschreibungstext |
| `event_date` | String | Datum des Events (z. B. `2026-07-01`) |
| `is_active` | Boolean | Aktiv-Status. Mehrere Events können gleichzeitig aktiv sein. |

### 💡 Beispiel

**Request:**
```http
PATCH /api/events/2
Content-Type: application/json

{
  "is_active": true
}
```

### 📥 Responses
**✅ 200 OK** – Event erfolgreich aktualisiert (gibt das aktualisierte Event zurück).
**❌ 400 Bad Request** – Ungültige ID, fehlerhaftes JSON oder kein einziges Feld im Body.
**❌ 404 Not Found** – Es existiert kein Event mit dieser ID.
**❌ 500 Internal Server Error** – Fehler bei der Kommunikation mit der Datenbank.

---

### Endpunkt: Slot aktualisieren

Aktualisiert die Zeiten eines bestehenden Zeitslots anhand seiner ID. Nur für **Admins**.

* **URL:** `/api/slots/{id}`
* **Methode:** `PATCH`
* **Header:** `Content-Type: application/json`

### 📦 Erlaubte Felder (Payload)
Alle Felder sind **optional**. Schickt nur das mit, was aktualisiert werden soll.

| Feldname | Typ | Beschreibung |
| :--- | :--- | :--- |
| `start_time` | String | Startzeit (z. B. `09:00:00`) |
| `end_time` | String | Endzeit (z. B. `09:15:00`) |

### 💡 Beispiel

**Request:**
```http
PATCH /api/slots/3
Content-Type: application/json

{
  "start_time": "10:00:00",
  "end_time": "10:15:00"
}
```

### 📥 Responses
**✅ 200 OK** – Slot erfolgreich aktualisiert (gibt den aktualisierten Slot zurück).
**❌ 400 Bad Request** – Ungültige ID, fehlerhaftes JSON oder kein einziges Feld im Body.
**❌ 404 Not Found** – Es existiert kein Slot mit dieser ID.
**❌ 500 Internal Server Error** – Fehler bei der Kommunikation mit der Datenbank.