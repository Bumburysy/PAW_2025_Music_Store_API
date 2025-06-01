# 🎵 Music Store REST API 🎵 – Dokumentacja Techniczna

http://193.28.226.78:25565/

## ToDoList
- Dodać dane testowe dla orders i reviews - problem z połączeniem z danymi od users i albums
- Poprawić dokumentację Swaggera - usunąć niepotrzbne pola generowane automatycznie (może DTO dodać???)
- Dodać testy Postman

## 📚 Opis funkcjonalny i technologiczny

### 1. Przeznaczenie API

Music Store API służy jako backendowy serwis RESTful, umożliwiający zarządzanie zasobami związanymi z internetowym sklepem muzycznym. Umożliwia zarządzanie danymi takimi jak albumy muzyczne, użytkownicy, recenzje, zamówienia, a także procesy związane z uwierzytelnianiem i autoryzacją użytkowników. API zostało zaprojektowane do współpracy z frontendową częścią aplikacji lub zewnętrznymi systemami integrującymi.

### 2. Zakres funkcjonalny

System umożliwia:
- Obsługę albumów muzycznych – dodawanie, przeglądanie, aktualizowanie i usuwanie albumów muzycznych, w tym masowe dodawanie albumów.
- Rejestrację i zarządzanie użytkownikami – tworzenie kont, przeglądanie danych użytkowników, edycję, usuwanie i zarządzanie rolami użytkowników.
- Zarządzanie zamówieniami – tworzenie zamówień, przeglądanie historii, edytowanie statusów oraz zarządzanie danymi wysyłki.
- Moderację recenzji – dodawanie, przeglądanie, edytowanie i usuwanie recenzji użytkowników.
- Uwierzytelnianie i autoryzację – umożliwiające logowanie użytkowników oraz kontrolę dostępu na podstawie przypisanej roli (np. administrator, użytkownik).
- Wczytywanie danych testowych – szybkie ładowanie przykładowych danych (np. albumów, użytkowników) z plików JSON do bazy danych.
- Interaktywną dokumentację – generowaną na podstawie definicji OpenAPI (Swagger), umożliwiającą testowanie endpointów bezpośrednio z poziomu przeglądarki.

### 3. Autoryzacja użytkowników

System wykorzystuje autoryzację opartą na tokenach JWT (JSON Web Token). Użytkownicy muszą się zalogować, podając prawidłowy adres e-mail i hasło, aby uzyskać token autoryzacyjny, który następnie przesyłają w nagłówkach Authorization podczas wywoływania chronionych endpointów.

#### Proces autoryzacji:
- Użytkownik wysyła zapytanie HTTP POST na endpoint /login, podając swoje dane uwierzytelniające (email i hasło).
- Serwer weryfikuje dane i w przypadku powodzenia generuje i zwraca token JWT.
- Użytkownik używa otrzymanego tokena JWT do uzyskiwania dostępu do zasobów chronionych (np. albumy, zamówienia, recenzje).

#### Przykład zapytania logowania:
`
curl -X 'POST' \
  'http://localhost:25565/login' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "user@example.com",
  "password": "strongpassword"}
`
#### Zawartość zapytania:
`
{
  "email": "user@example.com",
  "password": "strongpassword"
}
`
#### Odpowiedź serwera (przykładowy token):
`
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
`
####  Użycie tokena w dalszych zapytaniach:

Token JWT należy umieszczać w nagłówku Authorization z prefiksem Bearer:
`
Authorization: Bearer <token>
`

### 4. Planowane endpointy

#### Uwierzytelnianie:
- POST /login – logowanie i generowanie tokena JWT

#### Obsługa albumów (/albums):
- GET /albums - pobranie listy albumów
- GET /albums/:id - pobranie danych konkretnego albumu
- POST /albums – dodanie nowego albumu
- POST /albums/bulk – masowe dodanie albumów
- PATCH /albums/:id – aktualizacja danych albumu
- DELETE /albums/:id – usunięcie albumu

#### Obsługa użytkowników (/users):
- GET /users – pobranie listy użytkowników
- GET /users/:id – pobranie danych konkretnego użytkownika
- POST /users – utworzenie nowego użytkownika
- PATCH /users/:id – aktualizacja danych użytkownika
- DELETE /users/:id – usunięcie użytkownika

#### Obsługa zamówień (/orders):
- GET /orders – pobranie wszystkich zamówień
- GET /orders/:id – pobranie zamówienia o podanym ID
- GET /orders/user/:userID – pobranie zamówień użytkownika o podanym ID
- POST /orders – utworzenie nowego zamówienia
- PUT /orders/:id – aktualizacja zamówienia
- PATCH /orders/:id/status – zmiana statusu zamówienia
- PUT /orders/:id/shipping – aktualizacja danych wysyłki zamówienia
- DELETE /orders/:id – usunięcie zamówienia

#### Obsługa recenzji (/reviews):
- GET /reviews – pobranie wszystkich recenzji
- GET /reviews/:id – pobranie recenzji o podanym ID
- GET /reviews/album/:albumID – pobranie recenzji dla danego albumu
- GET /reviews/user/:userID – pobranie recenzji użytkownika
- POST /reviews – utworzenie nowej recenzji
- PUT /reviews/:id – aktualizacja recenzji
- DELETE /reviews/:id – usunięcie recenzji

#### Dane testowe (/data):
- POST /data/load – wczytanie danych testowych (np. albumów, użytkowników)

#### Dokumentacja (Swagger):
- GET /swagger/*any – interaktywna dokumentacja REST API

### 5. Projekt bazy danych

System wykorzystuje MongoDB jako bazę danych NoSQL, która umożliwia elastyczne i szybkie przechowywanie danych w formacie dokumentów BSON (Binary JSON). Poniżej przedstawiono modele danych wykorzystywane w aplikacji.

#### Kluczowe założenia projektu bazy danych:
- MongoDB zapewnia elastyczność w przechowywaniu danych o różnej strukturze i rozmiarze.
- Dane o albumach, zamówieniach, recenzjach i użytkownikach są przechowywane w osobnych kolekcjach.
- Powiązania między kolekcjami realizowane są poprzez referencje (np. UserID, AlbumID).
- Dane wrażliwe (hasła) przechowywane są w postaci haszowanej za pomocą algorytmu bcrypt, a nie w postaci jawnej.

#### Modele danych:

#### Album:
Kolekcja albums przechowuje dane dotyczące albumów muzycznych:
- ID (_id): Unikalny identyfikator albumu.
- Title: Tytuł albumu.
- Artist: Wykonawca albumu.
- Genre: Gatunek muzyczny.
- Description: Opcjonalny opis albumu.
- ReleaseDate: Data wydania albumu.
- Tracks: Lista utworów w albumie.
- Price: Cena albumu.
- Quantity: Ilość dostępnych egzemplarzy.
- CoverURL: URL do okładki albumu.
- CreatedAt, UpdatedAt: Daty utworzenia i modyfikacji wpisu.

#### Order:
Kolekcja orders przechowuje informacje o zamówieniach użytkowników:
- ID (_id): Unikalny identyfikator zamówienia.
- UserID: Identyfikator użytkownika, który złożył zamówienie.
- Items: Lista pozycji zamówienia (OrderItem), zawierająca ID albumu, ilość i cenę jednostkową.
- Total: Łączna wartość zamówienia.
- Status: Status zamówienia (pending, processing, shipped, completed, cancelled).
- Shipping: Dane do wysyłki (ShippingDetails).
- CreatedAt, UpdatedAt: Daty utworzenia i aktualizacji zamówienia.

#### Review:
Kolekcja reviews zawiera recenzje użytkowników dotyczące albumów:
- ID (_id): Unikalny identyfikator recenzji.
- AlbumID: ID albumu, którego dotyczy recenzja.
- UserID: ID użytkownika wystawiającego recenzję.
- Rating: Ocena albumu w skali (np. 1-5).
- Comment: Komentarz do recenzji.
- CreatedAt: Data utworzenia recenzji.

#### User:
Kolekcja users przechowuje dane użytkowników systemu:
- ID (_id): Unikalny identyfikator użytkownika.
- FirstName, LastName: Imię i nazwisko użytkownika.
- Email: Adres e-mail użytkownika.
- PhoneNumber: Numer telefonu.
- PasswordHash: Zabezpieczony hash hasła użytkownika (pole Password używane tylko przy tworzeniu).
- Role: Rola użytkownika (np. admin, customer, employee).
- IsActive: Status aktywności konta.
- ShippingDetails: Dane adresowe użytkownika (ShippingDetails).
- CreatedAt, UpdatedAt: Daty utworzenia i aktualizacji konta.

#### ShippingDetails:
Podstruktura wykorzystywana w zamówieniach i danych użytkownika:
- Address: Adres dostawy.
- City: Miasto dostawy.
- PostalCode: Kod pocztowy.
- Country: Kraj dostawy.
- PhoneNumber: Numer telefonu kontaktowego.

### 6. 🚀 Technologie

- **Go**, **GIN**, **MongoDB**, **HTML5**, **CSS3**.

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