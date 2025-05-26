# 🎵 Music Store API – Dokumentacja Techniczna

Projekt REST API – Music Store (wersja podstawowa)

## ToDoList
- Dalsza rozbudowa do pełnego spełnienia wytycznych z poziomu I i II

## 📚 Opis funkcjonalny i technologiczny

### 1. Przeznaczenie
REST API dla sklepu muzycznego, umożliwiające:
-	przeglądanie dostępnych albumów
-	zarządzanie ofertą i stanami magazynowymi

### 2. Planowane Endpointy
Albumy (CRUD)
-	GET /albums – pobranie listy albumów
-	GET /albums/{id} – pobranie szczegółów albumu
-	POST /albums – dodanie nowego albumu
-	PATCH /albums/{id} – edycja albumu
-	DELETE /albums/{id} – usunięcie albumu

### 3. Zakres funkcjonalny
-	CRUD na kolekcji albumów w MongoDB (albums)
-	Dokumentacja API w Swagger UI

### 4. Projekt bazy danych (MongoDB)

Kolekcja albums:
```json
{
  "_id": ObjectId,
  "title": "Album Title",
  "artist": "Artist Name",
  "genre": "Genre",
  "price": 9.99,
  "quantity": 10
}
```

### 🚀 Technologie
- **Go**, **MongoDB**, **HTML5**, **CSS3**.

________________________________________

Wymagania projektu:

Poziom I
- projekt REST API (tylko opis: przeznaczenie, planowane endpointy, zakres funkcjonalny)
- dokumentacja API w Swaggerze
- projekt bazy danych (ERD albo opis kolekcji i dokumentów)
- działająca baza danych z kluczowymi elementami (tabelami, kolekcjami); nie trzeba realizować całego projektu bazy
- implementacja szkieletu aplikacji i podstawowych endpointów API (dowolna technologia)
- umożliwienie demonstracji tworzenia, odczytywania, modyfikowania i usuwania danych (CRUD) za pośrednictwem REST API

Poziom II
- wszystko z poziomu I
- sugerowana realizacja aplikacji w języku Go
- uzupełnienie implementacji bazy danych do pełnej postaci opisanej w projekcie
- przygotowanie skryptów inicjalizujacych bazę danych
- podstawowe testy jednostkowe i/lub REST API
- aktywne wykorzystywanie VCS (np. git), zdalne repozytorium do wglądu
- implementacja uwierzytelniania i autoryzacji użytkowników (JWT)