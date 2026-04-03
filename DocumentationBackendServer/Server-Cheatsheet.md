# 🔄 Quick-Redeploy Guide (Main Branch Update)

Nutze diesen Workflow, wenn du Änderungen im `main`-Branch fertig hast und die Live-API auf dem FH-Server aktualisieren willst.

## Schritt 1: Lokal bauen & hochladen (auf deinem Mac)

Damit Kubernetes merkt, dass es ein neues Image gibt, erhöhen wir die Versionsnummer (z.B. von `v1` auf `v2`).

```bash
# 1. Neues Image bauen (Tag erhöhen! v1 -> v2 -> v3...) 
docker buildx build --platform linux/amd64 -t ghcr.io/jakemoes/dual-job-dating-backend:v3 --push .
```

# Schritt 2: Auf dem Server aktualisieren (SSH)
**Einlogen:**
```bash
cd ~/msd-dual-jobdating/dual-job-date-server
```
**Yaml anbassen:**
```bash
nano deployment.yaml
```
Suche die Zeile image: ghcr.io/jakemoes/...:v1 und ändere sie auf das neue Tag (z.B. :v2).

**Änderungen übernehmen**
```bash
kubectl apply -f deployment.yaml
```

# Schritt 3: Erfolgskontrolle (Wichtig!)

```bash
# 1. Status checken (Sollte 'Running' sein)
kubectl get pods -n msd21-dual-dating-prod

# 2. Logs live verfolgen
kubectl logs -f deployment/dual-job-date-api -n msd21-dual-dating-prod
```