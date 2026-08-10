--Схема нужна для того, чтобы отделить таблицы нашего приложения от других таблиц в базе данных. 
--Это позволяет лучше организовать данные и избежать конфликтов имен.
CREATE SCHEMA todoapp 

CREATE TABLE todoapp.users (
    id           SERIAL                PRIMARY KEY,
    --version нужен для оптимистической блокировки
    version      BIGINT       NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL CHECK(char_length(full_name) BETWEEN 3 AND 100),
    phone_number VARCHAR(15)           CHECK(
        phone_number ~ '^\+[0-9]+$'
        AND 
        char_length(phone_number) BETWEEN 10 AND 15
    )
); 
--'~' - это оператор регулярного выражения в PostgreSQL.
-- ^ - начало строки
-- \+ - символ '+' (\ чтобы + воспринялся как литерал, а не как оператор)
-- [0-9]+ - одна или более цифр (могут повторяться)
-- $ - конец строки
--
--CHECK не замедляет работу, так как выполняется только при вставке или обновлении данных, а не при каждом запросе.

CREATE TABLE todoapp.tasks (
    id             SERIAL                   PRIMARY KEY,
    version        BIGINT                   DEFAULT 1,
    title          VARCHAR(100)    NOT NULL CHECK(char_length(title) BETWEEN 1 AND 100),
    description    VARCHAR(1000)            CHECK(char_length(description) <= 1000),
    completed      BOOLEAN         NOT NULL,
    created_at     TIMESTAMPTZ     NOT NULL,
    completed_at   TIMESTAMPTZ,

    CHECK(
    (completed = FALSE AND completed_at IS NULL) --Если задача не выполнена, то completed_at должно быть NULL
    OR
    (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)--Если задача выполнена, то completed_at должно быть не NULL и не раньше created_at
    ),
    
    author_user_id INTEGER         NOT NULL REFERENCES todoapp.users(id)
)