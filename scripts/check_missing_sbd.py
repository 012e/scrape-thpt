import sqlite3
import sys

con = sqlite3.connect("student.db")
start = int(sys.argv[1])
end = int(sys.argv[2])
cur = con.cursor()
for i in range(start, end+1):
    if cur.execute(f"select * from students where sbd={i}").fetchone() is None:
        print(i)

