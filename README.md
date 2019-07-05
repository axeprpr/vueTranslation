#### vue翻译
类似django
```
支持以下两种引用
_("xxxx")
_("xxxxx #abc#", {"abc": "xxxxxxxx"})
```
check_translation_tools是配套的检查翻译脚本。golang写的，需要自己build下。
```
➜  check_translation_tools git:(master) ✗ ../../../axe/WorkGit/class/cloudclass-dashboard/src/lang/check_translation_tools/check_translation --help
Usage of ../../../axe/WorkGit/class/cloudclass-dashboard/src/lang/check_translation_tools/check_translation:
  -f string
        File format to search. (default "vue")
  -m string
        File to match with json string in it. (default "zh_en.js")
  -p string
        Path to search. Default is pwd. (default "/Users/axe/WorkGit/class/cloudclass-dashboard/src/lang/check_translation_tools")
  -s string
        show file path.
  -t string
        auto translate.
```