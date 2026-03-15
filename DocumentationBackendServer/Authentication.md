# 🔐 JWT & Auth Guide für Frontend/Mobile

Dieses Dokument beschreibt, wie ihr euch gegenüber dem Go-Backend (API) ausweist.

## 1. Das Wichtigste auf einen Blick
* **Keine Passwörter an die API:** Passwörter werden nur einmalig beim Login an Supabase geschickt.
* **Bearer Token:** Die API identifiziert User ausschließlich über das **JWT (JSON Web Token)** im Header.
* **Base URL:** `http://49.13.22.97:443` (Aktuell über HTTP erreichbar).

---

## 2. Der Login-Ablauf (Authentifizierung)

Ihr nutzt das Supabase Frontend-SDK (JS/Flutter), um den User einzuloggen.

1.  **Frontend -> Supabase:** Login mit Email & Passwort.
2.  **Supabase -> Frontend:** Antwortet mit einem User-Objekt, das den `access_token` enthält.
3.  **Frontend -> Go-API:** Verwendet diesen `access_token` für alle weiteren Anfragen.

---

## 3. JWT Verwendung (Autorisierung)

Jeder Request an die API (außer der Health-Check `/`) benötigt den Token im **HTTP-Header**.

### Der Authorization Header
Der Token muss zwingend mit dem Präfix `Bearer` gesendet werden:

| Key             | Value                         |
| :-------------- | :---------------------------- |
| `Authorization` | `Bearer <DEIN_ACCESS_TOKEN>`  |

### Beispiel mit `fetch`:
```javascript
const response = await fetch('[http://49.13.22.97:443/api/students](http://49.13.22.97:443/api/students)', {
    method: 'GET',
    headers: {
        'Authorization': `Bearer ${session.access_token}`,
        'Content-Type': 'application/json'
    }
});
```
## 4. Geheim vs. Öffentlich (Wer darf was wissen?)

In der JWT-Architektur gibt es klare Regeln, welche Daten wo liegen dürfen:

| Element | Status | Beschreibung |
| :--- | :--- | :--- |
| **Passwort** | 🔴 **STRENG GEHEIM** | Wird nur für den Login-Request genutzt. Niemals speichern! |
| **JWT Access Token** | 🟡 **TEMPORÄR** | Euer digitaler Ausweis. Sicher im App-State/SessionStorage halten. |
| **Supabase Anon Key** | 🟢 **ÖFFENTLICH** | Darf im Frontend-Code stehen (ist nur zur Identifikation der App). |
| **JWT Secret (Go)** | 🔴 **TOP SECRET** | **Nur im Backend (`.env`)**. Wer das hat, kann eigene Token fälschen! |

---

## 5. Troubleshooting (404 & 401)

Falls die API nicht so antwortet wie erwartet, prüft folgende Punkte:

* **401 Unauthorized:**
    * Ist der Token im Header vorhanden?
    * Steht das Wort `Bearer ` (mit Leerzeichen!) vor dem Token?
    * Ist der Token abgelaufen? (Supabase Tokens halten oft nur 1 Stunde -> Refresh Token nutzen).
* **404 Not Found:**
    * Prüft die URL.

---

## 6. Beispiel-Test mit cURL

Ihr könnt eure Verbindung und euren Token direkt im Terminal testen, ohne die App zu starten:

```bash
curl -i -X GET \
  -H "Authorization: Bearer <DEIN_ACCESS_TOKEN>" \
  [http://49.13.22.97:443/api/students](http://49.13.22.97:443/api/students)
```

## 7. Den JWT erhalten (Frontend Workflow)

Um den `access_token` zu bekommen, nutzt ihr die Standard `signIn`-Methoden von Supabase. Hier ist der Ablauf in JavaScript/TypeScript:

### Beispiel: Login & Token-Extraktion
```javascript
import { createClient } from '@supabase/supabase-js'

const supabase = createClient('https://dtgigetmxmrqniibjsal.supabase.co', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImR0Z2lnZXRteG1ycW5paWJqc2FsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzI2MTQxNDMsImV4cCI6MjA4ODE5MDE0M30.DV_NZ7trFcTIcm8Hb_pdkv3cDvJSaJpd7KnF3twtL1s')

async function loginUser(email, password) {
  // 1. Authentifizierung bei Supabase
  const { data, error } = await supabase.auth.signInWithPassword({
    email: email,
    password: password,
  })

  if (error) {
    console.error("Login fehlgeschlagen:", error.message)
    return
  }

  // 2. Den JWT (Access Token) extrahieren
  // Dieser Token muss in den Authorization Header!
  const jwt = data.session.access_token 
  
  console.log("Dein JWT für das Backend:", jwt)
  
  // Tipp: Speichert den jwt im App-State oder SessionStorage
  return jwt
}
```

### 🔑 JWT direkt via Bruno holen (Supabase Auth)

Falls ihr den Token nicht aus der App kopieren wollt, könnt ihr euch mit Bruno direkt einen Token bei Supabase generieren.

#### 1. POST Request erstellen
* **Methode:** `POST`
* **URL:** `https://dtgigetmxmrqniibjsal.supabase.co/auth/v1/token?grant_type=password`
#### 2. Header setzen
Ihr müsst zwei wichtige Header mitschicken, damit Supabase weiß, wer anfragt:
* **apikey:** `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImR0Z2lnZXRteG1ycW5paWJqc2FsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzI2MTQxNDMsImV4cCI6MjA4ODE5MDE0M30.DV_NZ7trFcTIcm8Hb_pdkv3cDvJSaJpd7KnF3twtL1s`
* **Content-Type:** `application/json`

#### 3. Body (Login-Daten)
Wählt im Reiter **Body** -> **JSON** und gebt die User-Daten ein:
```json
{
  "email": "jon.snow@winterfell.com",
  "password": "12345678"
}
```

### 🚀 API & Auth Quick-Ref

| Ressource | Wert / URL |
| :--- | :--- |
| **Backend API (Base)** | `http://49.13.22.97:443` |
| **Supabase URL** | `https://dtgigetmxmrqniibjsal.supabase.co` |
| **Supabase Anon Key** | `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImR0Z2lnZXRteG1ycW5paWJqc2FsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzI2MTQxNDMsImV4cCI6MjA4ODE5MDE0M30.DV_NZ7trFcTIcm8Hb_pdkv3cDvJSaJpd7KnF3twtL1s` |
| **Auth-Endpunkt (Bruno)** | `https://dtgigetmxmrqniibjsal.supabase.co/auth/v1/token?grant_type=password` |
| **Header Key** | `Authorization` |
| **Header Value** | `Bearer <DEIN_JWT_TOKEN>` |

---

### 📍 Wichtige Endpunkte
* **Health Check:** `GET /`
* **Studenten:** `GET /api/students`
* **Aktive Firmen:** `GET /api/companies/active`
* **Studenten-Meetings:** `GET /api/students/{id}/meetings`