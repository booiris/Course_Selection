import csv, random, time

with open("course.csv", "w", newline="") as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(["Name", "Cap"])
    for _ in range(300):
        name = "C" + str(round(time.time() * 1000 * 1000))
        cap = random.randint(40, 300)
        writer.writerow([name, cap])