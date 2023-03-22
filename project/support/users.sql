-- tworzenie sekwencji
CREATE SEQUENCE public.user_id_seq
-- wartość początkowa sekwencji
START WITH 1
-- krok inkrementacji sekwencji
INCREMENT BY 1
-- brak minimalnej wartości 
NO MINVALUE
-- brak maksymalnej wartości
NO MAXVALUE
-- rozmiar cache dla wartości sekwencji
CACHE 1;

-- zmiana właściciela sekwencji na użytkownika postgres
ALTER TABLE public.user_id_seq OWNER TO postgres;

-- ustawienie domyślnej przestrzeni tabel na pustą wartość
SET default_tablespace = '';

-- ustawienie domyślnej metody dostępu do tabel na 'heap'
SET default_table_access_method = heap;

-- tworzenie tablicy users
CREATE TABLE public.users (
-- wartość domyślna dla kolumny id pobierana z sekwencji user_id_seq
id integer DEFAULT nextval('public.user_id_seq'::regclass) NOT NULL,
-- adres e-mail użytkownika
email character varying(255),
-- imię użytkownika
first_name character varying(255),
-- nazwisko użytkownika
last_name character varying(255),
-- hasło użytkownika
password character varying(60),
-- flaga informująca o aktywności użytkownika
user_active integer DEFAULT 0,
-- data i czas utworzenia użytkownika
created_at timestamp without time zone,
-- data i czas ostatniej aktualizacji użytkownika
updated_at timestamp without time zone
);

-- zmiana właściciela tabeli users na użytkownika postgres
ALTER TABLE public.users OWNER TO postgres;

-- ustawienie wartości aktualnej sekwencji na 1
SELECT pg_catalog.setval('public.user_id_seq', 1, true);

-- dodanie klucza głównego dla kolumny id tabeli users
ALTER TABLE ONLY public.users
ADD CONSTRAINT users_pkey PRIMARY KEY (id);

-- wstawienie wiersza z przykładowymi danymi do tabeli users
INSERT INTO "public"."users"("email","first_name","last_name","password","user_active","created_at","updated_at")
VALUES (E'hubert@example.com',E'Admin',E'User',E'$2a$12$1zGLuYDDNvATh4RA4avbKuheAMpb1svexSzrQm7up.bnpwQHs0jNe',1,
E'2023-03-18 00:00:00',E'2023-03-18 00:00:00');
