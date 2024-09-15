# для тестов нужно:
## создать venv и активировать
```sh
python3 -m venv venv
source venv/bin/activate
```
 
## установить библиотеки
```sh
pip install pytest

pip install jwt
pip install requests
```

## Запустить
```sh
pytest
```
*Примечание:*
    чтобы работал print достаточно добавить флаг `-s`
    но лучше делать через `assert False, "important print text"`