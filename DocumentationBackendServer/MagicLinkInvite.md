# 🛠 Architektur-Konzept: Supabase Invites & App-Routing

## 1. Das Problem: Die "Einmal-Ticket-Sackgasse"
Wenn wir über Supabase einen Invite-Link verschicken, ist dieser aus Sicherheitsgründen ein **Einmal-Ticket**.
Klickt ein Student in der Mail auf den Link, verifiziert Supabase den Token und **entwertet ihn sofort**.

Wenn unser Backend den User danach stur auf den nativen App-Link (`dualjob://setPassowort`) weiterleitet und der Student die App **nicht** installiert hat, knallt es:
* Der Browser wirft einen Fehler ("Adresse ungültig").
* Der Token ist verbrannt.
* Der Student steckt in einer Sackgasse und kommt nicht mehr an sein Passwort.

**Warum wir es nicht per Webseite & Zwischenablage lösen können:**
Man könnte eine Webseite dazwischenschalten, die den Token in die Zwischenablage (Copy & Paste) des Handys speichert, bevor der User in den App Store geht. **Das ist aber ein bekannter Dirty-Hack:** Wenn der User im App Store ist und während des Downloads eine WhatsApp-Nachricht kopiert, ist der Token überschrieben und der User ausgesperrt. Zudem werfen moderne iOS/Android-Versionen gruselige Warnmeldungen ("App möchte aus Safari einsetzen").

---

## 2. Die Lösung: "Deferred Deep Linking" (Der Industrie-Standard)
Wir nutzen einen fertigen Service (wie **Branch.io** oder **AppsFlyer**). Das ist der Standardweg, wie professionelle Apps das heute lösen.

* **Wie es funktioniert:** Branch.io generiert uns eine "Smart URL".
* **Der Router:** Klickt der User, landet er bei Branch. Branch erkennt, ob die App da ist. Wenn nicht, zeigt Branch **automatisch** eine gehostete Download-Seite an und leitet in den App Store.
* **Die Magie (Fingerprinting):** Branch merkt sich den Token auf *ihren* Servern ("iPhone 14 aus Graz hat gerade geklickt"). Wenn die App nach dem Download zum ersten Mal öffnet, fragt die App bei Branch nach und bekommt den Token sicher übergeben.

---

## 3. Wer macht was? (Aufgabenverteilung)

### ⚙️ Server Team (Backend / API)
* Unser Job ist bereits erledigt. Das Backend verschickt die Invites via Supabase und baut die Weiche (Student vs. Company).
* **Wir brauchen:** Nur die finalen URLs von euch. Sobald das App-Team den Branch-Link hat, tragen wir ihn im Backend ein.

### 📱 App Team (Students)
* **Step 1:** Richtet einen Deferred Deep Linking Service (z.B. Branch.io) ein.
* **Step 2:** Gebt dem Server-Team die generierte Smart-URL (z.B. `https://dualjob.app.link/invite`).
* **Step 3:** Baut das SDK in die App ein, damit die App beim ersten Start nach dem Download den Token abfragen und verarbeiten kann.

### 🌐 Web Team (Companies)
* **Studenten-Flow:** Gute Nachrichten – ihr seid hier komplett raus! Die Download-Landingpage für die Studenten wird vom Deep-Link-Service (Branch.io) automatisch gehostet.
* **Company-Flow (Euer Task):** Wenn eine Firma eingeladen wird, leiten wir sie auf euer Web-Portal weiter (z.B. `https://portal.dualjob.de/invite`). Supabase hängt den Token als Parameter an (`#access_token=...`). Ihr müsst diesen Token via JavaScript abgreifen und die Company einloggen.

---

## 4. WICHTIG: Der Fallback ("Neuen Link anfordern")
Selbst Branch.io hat eine Fehlerquote von ca. 5-10 %. Wenn ein Student **Apple Private Relay** aktiviert hat, den **Brave Browser** nutzt oder einen **Adblocker** (z.B. Pi-hole) hat, schlägt das Fingerprinting fehl. Der User lädt die App, aber die App bekommt keinen Token.

**Der Rettungsschirm:**
* **App Team:** Baut in der App einen Fallback-Screen ein: *"Einladung nicht gefunden? [E-Mail eingeben] ⮕ [Neuen Link anfordern]"*.
* **Server Team:** Wir bauen dafür einen Endpoint (z.B. `/api/resend-invite`), der eine frische Supabase-Mail verschickt.
* **Der Clou:** Da die App *jetzt* installiert ist, funktioniert der Klick in der neuen E-Mail zu 100 % sofort über den direkten Link!

---

## ⚡ TL;DR (Für Lesefaule)
1. E-Mail-Links sind Einmal-Tokens. Direkte App-Links (`dualjob://`) sperren User ohne App aus.
2. Das App-Team muss einen Service wie **Branch.io** nutzen, der den Token sicher "merkt" und durch den App Store schleust.
3. Das Web-Team baut keine Student-Landingpage, sondern kümmert sich nur um das Auslesen des Tokens beim Company-Login im Web.
4. Für den Fall, dass Branch blockiert wird (Adblocker etc.), muss die App einen **"Neuen Link anfordern"-Button** haben, den das Backend dann verarbeitet.
5. Das Server-Team wartet auf die Branch.io-URL vom App-Team, um sie im Backend scharfzuschalten.