# 🚀 Guide: Invite Flow & Passwort-Setup (Web & App)

Dieses Dokument beschreibt den Ablauf, wie neu eingeladene User (Studenten & Firmen) ihr initiales Passwort setzen. Das Backend übernimmt die Erstellung der Accounts und den E-Mail-Versand. Frontend (Web) und App müssen lediglich den Redirect abfangen.

## 🔄 Der generelle Ablauf
1. **Backend:** Admin lädt User hoch -> Backend triggert Supabase Invite.
2. **E-Mail:** User bekommt Mail mit Link: `"Einladung annehmen"`.
3. **Redirect:** User klickt -> Supabase hängt einen Session-Token an die URL an und leitet weiter an Web oder App.
4. **Frontend/App:** Fängt den Token ab, lässt den User ein Passwort eintippen und speichert es in Supabase.

---

## 🛤️ Die zwei Wege (Routing)
Da Firmen das Webportal nutzen und Studenten die App, nutzt das Backend beim Versenden der Einladung unterschiedliche `redirect_to` Parameter:

- **Firmen (Web):** Wir leiten weiter auf z.B. `https://unsere-web-app.com/set-password`
- **Studenten (Mobile App):** Wir nutzen Deep Linking! Wir leiten weiter auf z.B. `jobdateapp://set-password` (Bitte gebt mir Bescheid, wie euer Deep-Link-Scheme genau heißt!)

---

## 🛠️ Was müsst ihr tun? (Frontend & App Team)

### Schritt 1: Die URL / den Deep Link abfangen
Wenn der User auf eurer Seite oder in eurer App landet, sieht die URL so aus:
`[EURE_URL]#access_token=eyJhb...&refresh_token=...&type=invite`

*Wichtig:* Die Supabase Client SDKs (JS/TS oder Flutter/Swift) erkennen diesen Token in der URL oft schon **automatisch** und loggen den User temporär ein!

### Schritt 2: Die UI bauen
Baut einen simplen Screen `SetPassword`:
- Eingabefeld 1: `Neues Passwort`
- Eingabefeld 2: `Passwort bestätigen`
- Button: `Passwort speichern & Einloggen`

### Schritt 3: Das Passwort an Supabase senden
Wenn der User den Button klickt, ruft ihr einfach die Standard-Update-Funktion des Supabase SDKs auf.

**Beispiel für Web (JavaScript / TypeScript):**
```javascript
const handlePasswordSet = async (newPassword) => {
  const { data, error } = await supabase.auth.updateUser({
    password: newPassword
  });

  if (error) {
    console.error("Fehler beim Speichern:", error.message);
  } else {
    // Erfolg! User in den geschützten Bereich weiterleiten (z.B. Dashboard)
    window.location.href = "/dashboard";
  }
}
```

### ⚠️ Wichtige To-Dos für unser Team
Frontend-Team: Baut die /set-password Route.

App-Team: Definiert euer Deep-Link-Scheme (z.B. app://...) und teilt es mir mit.

Supabase-Config: Sobald eure URLs/Deep-Links stehen, müssen wir diese im Supabase Dashboard unter Authentication -> URL Configuration -> Redirect URLs auf die Whitelist setzen.