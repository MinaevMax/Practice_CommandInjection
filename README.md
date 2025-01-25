# Practice_CommandInjection

## Краткое описание уязвимости
OS Command Injection - уязвимость, которая позволяет злоумышленнику выполнять его команды в консоли сервера, что не предусмотрено разработчиком.
Возникает, когда данные, принятые сервером не валидируются, не экранируются и даже не проверяются на безопасность. После чего они передаются в функцию выполняющую серверную консольную команду и срабатывают так, как задумывал злоумышленник. Отличительным знаком использования такой уязвимости является использование специальных символов в запросе, например ";", "%20"(пробел) и тд.
Данный сервис позволяет получить все штрафы пользователя по имени человека, добавить новый штраф и посмотреть статистику данных.


## Способы защиты от уязвимости
- Корректная проверка и фильтрация входных данных.
- По возможности избегание исполнения консольных команд(замена на возможности языка)
- Экранирование ввода.

## Запуск приложения
Для запуска приложения нужно:
1. Склонировать репозиторий:
```bash
git clone https://github.com/MinaevMax/Practice_CommandInjection.git
cd Practice_CommandInjection/
```
2. Запуск приложения через docker-compose:
```bash
docker-compose up --build
```
3. Приложение будет запущено по адресу: [http://localhost:8080](http://localhost:8080)

## Описание уязвимости конкретно в данном приложении
Флаг находится в файле admin.txt, лежащем в папке admin, но защещенной динамическим паролем, который равен размеру папки tmp. Сама папка лежит в директории bills, вместе с папками штрафов других людей.
Есть хэндлер, позволяющий получить штрафы по имени; хэндлер, позволяющий добавить штраф с полями имя и размер; хэндлер, позволяющий получить количество записей в хранилище. Для того, чтобы получить флаг, нужно передать вместо имени команду, который изменяет размер папки tmp на нужное значение.

## POC-эксплоит для получения флага
Для получения флага необходимо получить файл admin.txt из папки admin, защищенную паролем. 

Для начала, стоит попробовать вбить в поиск штрафов по имени `admin`. Заметим, что нам сообщают о неверном пароле, подсказка: пароль равен количеству файлов в папке tmp. Значит нам необходимо либо получить размер этой директории, либо изменить его на нужное нам значение. Рассмотрим второй вариант. Нам необходимо удалить все файлы из папки tmp, и тогда пароль к папке станет "0". Попробуем это сделать с помощью ввода: 
```
Name1;rm -rf ./tmp/*
```
Этой командой мы не вызовем ошибки программы, тк она получила имя для поиска и в то же время удалили содержимое нужной нам папки. После этого пробуем имя: `admin` и пароль `0`. Программа должна вывести искомый флаг.

Выполнение запросов лучше и удобнее всего проводить через любой удобный для Вас браузер.