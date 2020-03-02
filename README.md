# Olanzapine

Olanzapine - CLI todo приложение, работающее с базой данных PostgreSQL.

## Первый запуск
+ Скачайте репозиторий:
        
        git clone https://github.com/qrlzvrn/Olanzapine
        
+ Дальше необходимо, что бы на вашем компьютере был установлен **postgres**.

+ Если его нет, то установите.

+ Создайте базу данных **olnaza**.
       
       createdb olanza
        
+ Запустите программу с командой **init**.
       
       olanza init
        
+ Теперь вы можете создать вашу первую задачу.

## Использование

**add** - добавить новую задачу

Ключи:

   + --content/-C - задает название задачи
   + --category/-c - задает категорию
   + --deadline/-d - задает дедллайн
   
**complete** - выполнить задачу

Данная команда меняет состояние задачи на выполненное, задача не удаляется из базы

        olanza complete <id задачи>
        
**delete** - удалить задачу

        olanza delete <id задачи>
        
**list / ls** - просмотреть список задач

Что бы просмотреть только какую-то определенную категорию укажите ее в качестве аргумента:

        olanza ls <название категории>
        
**reDead** - изменить дедлайн задачи

        olanza reDead <id задачи> <новый дедлайн>