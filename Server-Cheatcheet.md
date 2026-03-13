# 🚀 Hetzner Server Cheat Sheet - MSD Dual Jobdating

Dieses Dokument enthält alle wichtigen Befehle zur Verwaltung des Go-Backends auf dem Hetzner-Server.

---

## 📂 Verzeichnisse & Dateien
Der Code liegt unter: `~/msd-dual-jobdating/dual-job-date-server`

- `ls` : Zeigt alle Dateien (inkl. versteckter `.env`)
- `nano .env` : Bearbeitet die Umgebungsvariablen (Datenbank-Keys, Ports)
- `cat .env` : Schneller Blick in die Konfiguration

---

## 🐳 Docker Management

### Status prüfen
- `docker ps` : Zeigt laufende Container
- `docker ps -a` : Zeigt alle Container (auch gestoppte/gecrashte)
- `docker logs -f go-server` : Live-Logs der App (Beenden mit `Ctrl+C`)
- `docker stats go-server` : CPU & RAM Verbrauch prüfen

### Steuerung
- `docker stop go-server` : App anhalten
- `docker start go-server` : App wieder starten
- `docker rm -f go-server` : Container komplett löschen (erzwingen)

---

## 🔄 Deployment Workflow (Updates einspielen)

Wenn du neuen Code auf GitHub gepusht hast, führe diese Befehle nacheinander aus:

1. **Code aktualisieren:**
   ```bash
   git pull

2. **Image neu bauen:**
    ```bash
    docker build -t msd23-backend .

3. **Container neu starten (Löschen & Run):**
    ```bash
   docker rm -f go-server
    docker run -d -p 80:80 
   --name go-server 
   --restart always 
   --env-file .env msd23-backend

## Verbindungsinfos
* Base URL: http://49.13.22.97:8080

* SSH Login: ssh root@49.13.22.97

