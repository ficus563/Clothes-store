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
CREATE TABLE Photo (
    ID TEXT PRIMARY KEY,
    ProductID TEXT FOREIGN KEY (ProductID) REFERENCES Product(ID),
    Size INT,
    Price DECIMAL (10,2)
)

--Таблица "Роли"
CREATE TABLE Role (
    ID TEXT PRIMARY KEY,
    Name TEXT
)

--Таблица "Пользователь"
CREATE TABLE User (
    ID TEXT PRIMARY KEY
    Name TEXT,
    Email TEXT,
    Phone TEXT,
    UserRole TEXT FOREIGN KEY (UserRole) REFERENCES Role(ID),
    RegistrationDate DATE 
)

--Таблица "Корзина"
CREATE TABLE Cart (
    ID TEXT PRIMARY KEY,
    UserID TEXT FOREIGN KEY (UserID) REFERENCES User(ID),
    ProductID TEXT FOREIGN KEY (ProductID) REFERENCES Product(ID),
    Count INT
)

--Таблица Заказ
CREATE TABLE Order (
    ID TEXT PRIMARY KEY,
    UserID TEXT FOREIGN KEY (UserID) REFERENCES User(ID),
    DateTimes DATE,
    Status TEXT
)

--Таблица "Элементы заказа"
CREATE TABLE OrderItem (
    ID TEXT PRIMARY KEY,
    OrderID TEXT FOREIGN KEY (OrderID) REFERENCES Order(ID),
    ProductID TEXT FOREIGN KEY (ProductID) REFERENCES Product(ID),
)

