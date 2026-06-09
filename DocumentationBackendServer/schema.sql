-- ============================================================================
-- Dual Job Dating — Datenbankschema (ER-Referenz)
-- ============================================================================
-- Kanonische Referenz des Postgres-Schemas der Supabase-Instanz.
-- Nur zur Dokumentation gedacht (Tabellen-Reihenfolge/Constraints ggf. für
-- direkte Ausführung anzupassen). Quelle: Supabase `public`-Schema.
--
-- Stand: Matching-Consolidation (event-bezogenes Matching + event-eigene Slots).
-- Gegenüber dem ursprünglichen Schema wurden ergänzt:
--   * slots.event_id          (NEU)  -> events(id) ON DELETE CASCADE
--   * meetings.event_id FK     (geändert) -> jetzt ON DELETE CASCADE
--   * Indizes idx_meetings_event_id, idx_slots_event_id
-- Das genaue Migrations-Delta steht am Ende dieser Datei.
-- ============================================================================

-- Enum für Studenten-Präferenzen (student bewertet ein Unternehmen)
CREATE TYPE public.preference_type AS ENUM ('like', 'dislike', 'neutral', 'none');

-- Benutzerkonto (verknüpft mit Supabase Auth über user_id)
CREATE TABLE public.users (
  id uuid NOT NULL,
  user_id uuid NOT NULL UNIQUE,          -- = auth.users.id (Supabase Auth UUID)
  role text NOT NULL,                    -- 'admin' | 'student' | 'company'
  first_name text,
  last_name text,
  CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- Studierendenprofil
CREATE TABLE public.students (
  id integer NOT NULL DEFAULT nextval('students_id_seq'::regclass),
  user_id uuid NOT NULL UNIQUE,
  study_program text,
  semester integer,
  CONSTRAINT students_pkey PRIMARY KEY (id),
  CONSTRAINT students_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

-- Unternehmensprofil
CREATE TABLE public.companies (
  id integer NOT NULL DEFAULT nextval('companies_id_seq'::regclass),
  user_id uuid NOT NULL UNIQUE,
  name text,
  description text,
  short_description text,
  website text,
  logo_url text,
  image_urls text,                       -- mehrere URLs, durch ';' getrennt
  active boolean DEFAULT true,
  last_updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT companies_pkey PRIMARY KEY (id),
  CONSTRAINT companies_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

-- Event (Job-Dating-Veranstaltung)
CREATE TABLE public.events (
  id integer NOT NULL DEFAULT nextval('events_id_seq'::regclass),
  name text NOT NULL,
  location text,
  description text,
  event_date date NOT NULL,
  is_active boolean DEFAULT true,
  CONSTRAINT events_pkey PRIMARY KEY (id)
);

-- Zeitslot. Slots sind seit der Matching-Consolidation event-eigen (event_id).
CREATE TABLE public.slots (
  id integer NOT NULL DEFAULT nextval('slots_id_seq'::regclass),
  start_time time without time zone NOT NULL,
  end_time time without time zone NOT NULL,
  event_id integer,                      -- NEU: Slot gehört zu einem Event
  CONSTRAINT slots_pkey PRIMARY KEY (id),
  CONSTRAINT slots_event_id_fkey FOREIGN KEY (event_id)
    REFERENCES public.events(id) ON DELETE CASCADE
);

-- Präferenz: Student bewertet Unternehmen (Grundlage für das Matching)
CREATE TABLE public.preferences (
  id integer NOT NULL DEFAULT nextval('preferences_id_seq'::regclass),
  student_id integer,
  company_id integer,
  preference_type public.preference_type NOT NULL,
  CONSTRAINT preferences_pkey PRIMARY KEY (id),
  CONSTRAINT preferences_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(id),
  CONSTRAINT preferences_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id)
);

-- Meeting = ein zugeteiltes "Date" (Student trifft Unternehmen in einem Slot
-- eines Events). event_id verknüpft das Meeting mit seinem Event; ON DELETE
-- CASCADE räumt beim Löschen eines Events dessen Meetings (und Slots) automatisch auf.
CREATE TABLE public.meetings (
  id integer NOT NULL DEFAULT nextval('meetings_id_seq'::regclass),
  event_id integer,
  slot_id integer,
  student_id integer,
  company_id integer,
  CONSTRAINT meetings_pkey PRIMARY KEY (id),
  CONSTRAINT meetings_event_id_fkey FOREIGN KEY (event_id)
    REFERENCES public.events(id) ON DELETE CASCADE,
  CONSTRAINT meetings_slot_id_fkey FOREIGN KEY (slot_id) REFERENCES public.slots(id),
  CONSTRAINT meetings_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(id),
  CONSTRAINT meetings_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id)
);

-- Hilfreich für das event-bezogene Matching (Filter nach event_id)
CREATE INDEX IF NOT EXISTS idx_meetings_event_id ON public.meetings(event_id);
CREATE INDEX IF NOT EXISTS idx_slots_event_id    ON public.slots(event_id);

-- Hinweis: Alle public-Tabellen haben RLS aktiviert ohne Policies. Der Zugriff
-- erfolgt ausschließlich über das Go-Backend mit dem service_role-Key (umgeht RLS).
-- Clients (Web/Mobile) sprechen die Daten-Tabellen nicht direkt über PostgREST an.

-- ============================================================================
-- Migration-Delta der Matching-Consolidation (auf bestehender DB angewendet)
-- ============================================================================
-- -- meetings.event_id existierte bereits; FK auf CASCADE umgestellt:
-- ALTER TABLE public.meetings DROP CONSTRAINT meetings_event_id_fkey;
-- ALTER TABLE public.meetings ADD CONSTRAINT meetings_event_id_fkey
--   FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;
-- -- slots.event_id neu hinzugefügt:
-- ALTER TABLE public.slots ADD COLUMN event_id integer
--   REFERENCES public.events(id) ON DELETE CASCADE;
-- -- Indizes:
-- CREATE INDEX IF NOT EXISTS idx_meetings_event_id ON public.meetings(event_id);
-- CREATE INDEX IF NOT EXISTS idx_slots_event_id    ON public.slots(event_id);
