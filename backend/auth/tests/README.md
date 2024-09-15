# для тестов нужно:
## создать venv и активировать
```sh
python3 -m venv venv
source venv/bin/activate
```
 
## установить библиотеки
```sh
pip install requirements.txt
```



## Запустить
```sh
pytest
```
*Примечание:*
    чтобы работал print достаточно добавить флаг `-s`
    но лучше делать через `assert False, "important print text"`

## Если устанавливаем новые библиотеки
```sh
pip install example
```

нужно обновить requirements.txt
```sh
pip freeze > requirements.txt
```