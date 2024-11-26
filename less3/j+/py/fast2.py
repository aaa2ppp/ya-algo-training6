#давай ломать дальше

from sys import stdin

data = stdin.read().splitlines()
n, H = map(int, data[0].split()) 
hw = sorted(zip(map(int, data[1].split()), map(int, data[2].split())), key = lambda x: (-x[0], -x[1])) 
 
def brute(n, H, hw): 
    
    if n == 1: 
        return 0 
     
      
    mindelta = float('inf') 
    length = 0 
    maxI = n - 1 
 
    for i in range(n-1, -1, -1): 
        length += hw[i][1] 
        if length >= H: 
            maxI = i + 1 
            break 
 
    i = 0 
    while i < maxI: 
        length = hw[i][1] 
        delta = 0         
        if  length >= H: return delta     
        next = i + 1 
        for j in range(i+1, n): 
            length += hw[j][1] 
            if -hw[j][0]+hw[j-1][0] >= delta: 
                delta = -hw[j][0]+hw[j-1][0] 
                next = j 
            if length >= H: 
                mindelta = min(mindelta, delta) 
                if mindelta == 0 or j == n - 1: return mindelta           
                break  
        i = next 
    return mindelta 
 
print(brute(n, H, hw))
