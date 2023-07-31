import sqlite3
import re
import threading

start = 51000001
end = 51019942
thread_count = 5
def chunks(xs, n):
    n = max(1, n)
    return (xs[i:i+n] for i in range(0, len(xs), n))

def efficient_task_split(start, end, con):
    pack_count = (end - start + 1) // con
    remaining = (end - start + 1) % con
    share = [pack_count + 1 if i < remaining else pack_count for i in range(con)]

    split = [[0, 0] for _ in range(con)]
    current_index = start
    for i in range(con):
        split[i][0] = current_index
        split[i][1] = current_index + share[i] - 1
        current_index += share[i]

    return split


def convert_string_to_dict(input_string):
    data = {}

    # Use regular expression to find subject names and grades
    pattern = r'(\w+(?: \w+)*):\s+(\d+\.\d+)'
    matches = re.findall(pattern, input_string)

    # Populate the data dictionary
    for subject, grade in matches:
        data[subject] = float(grade)

    return data

default_dict = {
    "Toán": None,
    "Ngữ văn": None,
    "Lịch sử": None,
    "Vật lí": None,
    "Địa lí": None,
    "Sinh học": None,
    "Hóa học": None,
    "KHTN": None,
    "KHXH": None,
    "Tiếng Anh": None,
    "GDCD": None,
    "Tiếng Trung": None,
    "Tiếng Nhật": None,
    "Tiếng Pháp": None,
}
to_ascii = {
    "Toán": "toan",
    "Ngữ văn": "van",
    "Lịch sử": "su",
    "Vật lí": "ly",
    "Địa lí": "dia",
    "Sinh học": "sinh",
    "Hóa học": "hoa",
    "KHTN": "khtn",
    "KHXH": "khxh",
    "Tiếng Anh": "anh",
    "GDCD": "gdcd",
    "Tiếng Trung": "trung",
    "Tiếng Nhật": "nhat",
    "Tiếng Pháp": "phap",
}

con = sqlite3.connect("student.db", check_same_thread=False)
cur = con.cursor()

def apply_between(start: int, end: int, thread_count: int = 10):
    print(f"starting thread (begin = {start})")
    for id in range(start, end + 1):
        result = cur.execute(
            f"select unparsed_result from students where id={id}"
        ).fetchone()
        if result is None:
            print("failed: ", id)
            continue
        result = convert_string_to_dict(result[0])
        temp = default_dict.copy()
        temp.update(result)
        result = temp
        for subject, score in result.items():
            if score is None:
                score = "NULL"
            cur.execute(sql_set_score(id, to_ascii[subject], score))
        print(id)

def sql_set_score(id, subject, score):
    return """
    UPDATE students
    SET {} = {}
    WHERE
        id={}
    """.format(subject, score, id)

apply_between(start, end, thread_count)
# threads = []
# for split in efficient_task_split(start, end, thread_count):
    # thread = threading.Thread(target=apply_between, args=(split[0], split[1], thread_count))
    # apply_between(split[0], split[1], thread_count)
    # threads.append(thread)
    # thread.start()

# execute threads1
# for i, thread in enumerate(threads):
#     thread.join()
#     print(f"thread[{i}] finished")
con.commit()
