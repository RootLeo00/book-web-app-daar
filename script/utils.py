from collections import Counter
import re
from stop_words import safe_get_stop_words
 
def get_word_occurence(text, lang):
    stop_words = safe_get_stop_words(lang)
    txt = text.lower()
    words = re.findall('[a-zA-Z\u00C0-\u00FF]*', txt)
    lst = [x for x in words if x != '' and len(x) > 2 and not x in stop_words]
    wordsOcc = Counter(lst)
    return wordsOcc


def jaccard_distance(d1, d2):
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