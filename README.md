# Tages Test
## ТЗ
Необходимо написать сервис на Golang работающий по gRPC.
### 1) Принимать бинарные файлы (изображения) от клиента и сохранять их на жесткий диск.
За сохранение изображений отвечает методы:
~~~bash
rpc LoadImage (ImageRequest) returns (Empty) {}
rpc LoadImageStream (stream ImageRequest) returns (Empty) {}
~~~
Входные параметры следующие:
~~~bash
message ImageRequest{
    bytes Data = 1; // файл, представленный массивом байтов
    string Name = 2; // имя файла
}
~~~
Т.к. массив байтов подразумевает с собой только внутреннее содержимое файла(без метаинформации (н.р имя файла)), нужно передавать имя файла(с расширением) параметром функции.
Все файлы сохраняются на сервере в папке `images`
#### Возникшие проблемы и их решение
Возник вопрос в том, что при одном подключении к серверу, клиент сможет отправить только одну картинку. Что если клиенту нужно отправить сразу несколько картинок? Будут случаи (если сервер пользуется популярностью), когда некоторые картинки будут загружены, а некоторые нет из-за ограничения к-ва подключений к серверу.
Возможные решения:
1) Принимать массив бинарных файлов
2) Сделать возможность 'стримить' со стороны клиента изображения
3) Гибрид?
Первый способ прост в реализации, но имеет, по моему мнению, недостатки. Если суммарный вес, отправляемых файлов, будет слишком большой, то это ударит по производительности сервера, особенно если таких запросов несколько.
Второй вариант тоже прост в реализации и избавляет от недостатка первого способа, потому что он будет передавать файлы не 'разом', а постепенно. Есть конечно вопросы по поводу к-ва возможных 'стримов' с одного подключения, но пока не могу ответить на этот вопрос из-за недостатка опыта.
Можно смешать первый и второй вариант, тем самым мы избавимся от количества 'стримов', если с ними возникнет проблема. Но при этом должна немного упать производительность.
#### Вывод
Я выбрал второй вариант решения проблемы, потому что он показался мне более разумным.
Оставил первый вариант в качестве примера.
### 2) Иметь возможность просмотра списка всех загруженных файлов в формате: Имя файла | Дата создания | Дата обновления
За вывод всех файлов отвечает метод:
~~~bash
rpc GetImagesInfo (Empty) returns (ImagesInfo) {}
~~~
Выходные параметры следующие:
~~~bash
message ImagesInfo{
    repeated ImageInfo images = 1;
}

message ImageInfo{
    string Name = 1;
    string CreateAt = 2;
    string UpdateAt = 3;
}
~~~
#### Возникшие проблемы и их решение
Проблема этого метода состоит в том, что выводится сразу весь список изображений. В дальнейшем (При росте популярности приложения) нужно будет реализовать пагинацию и возвращать клиенту данные, разделенные по стримам. Хоть эти данные весят очень мало по сравнению с картинками, лучше всего их отправлять частично.