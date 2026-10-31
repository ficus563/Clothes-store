PRAGMA foreign_keys = ON;

-- Создание таблиц и связей между ними
-- Таблица "Товар"
CREATE TABLE IF NOT EXISTS Product (
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
);

-- Таблица "Фото"
CREATE TABLE IF NOT EXISTS Photo (
    ID TEXT PRIMARY KEY,
    ProductID TEXT,
    Filename TEXT NOT NULL,
    FOREIGN KEY (ProductID) REFERENCES Product(ID) ON DELETE CASCADE
);

-- Таблица "Роли"
CREATE TABLE IF NOT EXISTS Role (
    ID TEXT PRIMARY KEY,
    Name TEXT NOT NULL
);

-- Таблица "Пользователь"
CREATE TABLE IF NOT EXISTS "User" (
    ID TEXT PRIMARY KEY,
    NameUser TEXT,
    Email TEXT,
    Phone TEXT,
    UserRole TEXT,
    RegistrationDate DATE,
    FOREIGN KEY (UserRole) REFERENCES Role(ID)
);

-- Таблица "Корзина"
CREATE TABLE IF NOT EXISTS Cart (
    ID TEXT PRIMARY KEY,
    UserID TEXT,
    ProductID TEXT,
    Count INT,
    FOREIGN KEY (UserID) REFERENCES "User"(ID),
    FOREIGN KEY (ProductID) REFERENCES Product(ID)
);

-- Таблица Заказ
CREATE TABLE IF NOT EXISTS "Order" (
    ID TEXT PRIMARY KEY,
    UserID TEXT,
    DateTimes DATE,
    Status TEXT,
    FOREIGN KEY (UserID) REFERENCES "User"(ID)
);

-- Таблица "Элементы заказа"
CREATE TABLE IF NOT EXISTS OrderItem (
    ID TEXT PRIMARY KEY,
    OrderID TEXT,
    ProductID TEXT,
    Price DECIMAL(10,2),
    FOREIGN KEY (OrderID) REFERENCES "Order"(ID),
    FOREIGN KEY (ProductID) REFERENCES Product(ID)
);
