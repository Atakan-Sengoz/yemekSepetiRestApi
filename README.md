
# Yemeksepeti REST-API

In memory key-value store olarak çalışan bir REST-API servisidir.
Key değerini için get ve set metotlarını uygulayabilir.
Ayrıca dakikada bir o anki veriyi bir JSON dosyasına yazmaktadır.
API ayağa kalkarken eğer kayıtlı bir veri var ise value'nun değeri kayıtlı
verinin son verisine eşit olur. Kayıtlı olan tüm verileri gösterme imkanına sahiptir. 
İçerisinde Swagger kuruludur isterseniz test edebilirsiniz.


## API Kullanımı

#### O anki key değerini çağırmak için bir endpoint

```http
  GET /get
```

#### Tüm kaydedilmiş verileri çağırmak için bir için bir endpoint

```http
  GET /flush
```

#### Key değerini değiştirmek için bir endpoint

```http
  POST /post
```

| Parametre | Tip     | Açıklama                       |
| :-------- | :------- | :-------------------------------- |
| `value`      | `string` | **Zorunlu alandır**. Set edilmek istenen veriyi JSON olarak post etmeye yarar |

Gönderilmesi gereken örnek JSON
```json
{
  "value": "string"
}
```


## Swagger ile Demo

http://localhost:9090/swagger/index.html#/ adresine giderek swagger ile test edilebilir.


## Yazar

- Atakan Şengöz

  