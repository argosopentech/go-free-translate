# GoFreeTranslate
A translation application written in Golang using the [Gio UI graphics library](https://gioui.org/). GoFreeTranslate uses [LibreTranslate](https://libretranslate.com) for translations.

![Screenshot](https://community.libretranslate.com/uploads/default/original/2X/5/51c09026e28b3c603b64d4bf83ec121b216bd874.png)

### Quickstart

```
# First run a LibreTranslate instance at localhost:5000
git clone https://github.com/LibreTranslate/LibreTranslate.git
cd LibreTranslate
virtualenv env
source env/bin/activate
pip install -e .
argospm install translate-en_es
argospm install translate-es_en
libretranslate

# Then
git clone https://github.com/argosopentech/go-free-translate.git
cd go-free-translate
go run .
```
