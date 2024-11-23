from sys import stdin

# O(n) в среднем на рандомных данных, O(n^2) в худшем случае. см. тест 999

data = stdin.read().splitlines()
n, H = map(int, data[0].split())
hw = sorted(zip(map(int, data[1].split()), map(int, data[2].split())))

def brute(n, H, hw):
   
    if n == 1:
        return 0

    mindelta = float('inf')
    i = length = 0
    
    while i < n - 1:
        length = hw[i][1]
        delta = 0        
        if  length >= H: return delta    
        next = i + 1
        for j in range(i+1, n):
            length += hw[j][1]
            if hw[j][0]-hw[j-1][0] >= delta:
                delta = hw[j][0]-hw[j-1][0]
                next = j
            if length >= H:
                mindelta = min(mindelta, delta)
                if mindelta == 0: return 0                
                break 
        i = next
    return mindelta

print(brute(n, H, hw))
