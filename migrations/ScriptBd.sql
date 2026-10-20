PRAGMA foreign_keys = ON;
--Cоздание таблиц и связей между ними
--Таблица "Товар"
CREATE TABLE Product IF NOT EXISTS (
    ID TEXT PRIMARY KEY,
    ProductName TEXT NOT NULL,
    ProductType TEXT,
    Gender TEXT,
    Color TEXT,
    ProductDescription TEXT,
    Brand TEXT,
    Season TEXT,
    ProductSize INT, 
    Price DECIMAL (10,2) NOT NULL
)

--Таблица "Фото"
CREATE TABLE Photo IF NOT EXISTS (
    ID TEXT PRIMARY KEY,
    ProductID TEXT FOREIGN KEY (ProductID) REFERENCES Product(ID) ON DELETE CASCADE,
    Filenames TEXT NOT NULL 
)

--Таблица "Роли"
CREATE TABLE Role IF NOT EXISTS (
    ID TEXT PRIMARY KEY,
    Names TEXT
)

--Таблица "Пользователь"
CREATE TABLE User IF NOT EXISTS (
    ID TEXT PRIMARY KEY
    NameUser TEXT,
    Email TEXT,
    Phone TEXT,
    UserRole TEXT FOREIGN KEY (UserRole) REFERENCES Role(ID),
    RegistrationDate DATE 
)

--Таблица "Корзина"
CREATE TABLE Cart IF NOT EXISTS (
    ID TEXT PRIMARY KEY,
    UserID TEXT FOREIGN KEY (UserID) REFERENCES User(ID),
    ProductID TEXT FOREIGN KEY (ProductID) REFERENCES Product(ID),
    Count INT
)

--Таблица Заказ
CREATE TABLE Order IF NOT EXISTS (
    ID TEXT PRIMARY KEY,
    UserID TEXT FOREIGN KEY (UserID) REFERENCES User(ID),
    DateTimes DATE,
    Status TEXT
)

--Таблица "Элементы заказа"
CREATE TABLE OrderItem IF NOT EXISTS (
    ID TEXT PRIMARY KEY,
    OrderID TEXT FOREIGN KEY (OrderID) REFERENCES Order(ID),
    ProductID TEXT FOREIGN KEY (ProductID) REFERENCES Product(ID),
)

--Запросы для создания ролей
INSERT INTO Role(ID, Names) VALUES ('514d27ae-e77d-4246-9d7a-8b6c820f11ae', 'Администратор')
INSERT INTO Role(ID, Names) VALUES ('7dbafc77-eaed-4ce7-94ae-34f046182350', 'Менеджер')
INSERT INTO Role(ID, Names) VALUES ('cdce0c51-aa68-4c32-ae31-8a5f106cbac1', 'Клиент')
INSERT INTO Role(ID, Names) VALUES ('11503e7a-a002-49bc-aac7-2d5a548a1c1a', 'Гость')

