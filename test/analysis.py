# %%
import matplotlib.pyplot as plt
import matplotlib.ticker as ticker

times = []
with open("/home/booiris/Course_Selection/test/temp3.log", "r") as f:
    for line in f:
        temp = line.split("|")
        if len(temp) == 5:
            times.append(temp[2])

res = []
for time in times:
    if time.find("µs")!=-1:
        res.append(float(time[:-3])/ 1000)
    else:
        
        res.append(float(time[:-3]) )
print(len(res), max(res), min(res))
plt.plot(res)
plt.savefig("/home/booiris/Course_Selection/test/temp3.png")
plt.show()
