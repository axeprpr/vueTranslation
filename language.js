import store from '../store'
import { language_list } from '@/lang/zh_en.js'

function myreplace(str, obj) {
    if (!obj) { return str }
    Object.keys(obj).forEach(function(key) {
        str = str.replace(new RegExp('#' + key + '#', 'g'), obj[key])
    })
    return str
}

export function __(str, object) {
    let language = store.getters.language
    if (language == 'zh' && language_list[str]) {
        return myreplace(language_list[str], object)
    } else {
        return myreplace(str, object)
    }
}