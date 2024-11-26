# генерирует тест, который ломает решение по TL

max_h = int(1e9)

n = 44720
h = 1
hh = []

while h <= max_h and n > 0:
    hh.append(h)
    h = h + n
    n -= 1

#print(n)
print(len(hh), len(hh))
print(*hh)
ww = [i for i in range(1, len(hh)+1)]
print(*ww)
