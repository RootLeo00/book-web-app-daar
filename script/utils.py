from collections import Counter
import re
import io
from stop_words import safe_get_stop_words
import pandas as pd
 
def getWordsWithOcc(text, lang):
    stop_words = safe_get_stop_words(lang)
    txt = text.lower()
    words = re.findall('[a-zA-Z\u00C0-\u00FF]*', txt)
    lst = [x for x in words if x != '' and len(x) > 2 and not x in stop_words]
    wordsOcc = Counter(lst)
    return wordsOcc


def getListIndexBook(wordsOcc, idBook):
    list = []
    for word, occ in wordsOcc.items():
        list.append({
            "word": word,
            "occurrence": occ,
            "idBook": idBook,
        })
    return list


def saveTmpWords(wordsOcc):
    stream = io.open("tmpWords.txt", 'a+')
    for word, occ in wordsOcc.items():
        stream.write(word + ",\n")
    stream.close()


def distance_jaccard(d1, d2):
    a = 0
    b = 0
    list_word_d2 = list(d2)
    for key, value in d1.items():
        if key in d2:
            m = max(value, d2[key])
            a += m - min(value, d2[key])
            b += m
            list_word_d2.remove(key)
        # words only in d1
        else:
            a += value
            b += value
    # words only in d2
    for word in list_word_d2:
        a += d2[word]
        b += d2[word]
    try:
        return a/b
    except:
        return 1