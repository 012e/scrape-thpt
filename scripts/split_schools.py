"""
Search where new batches start based on students' name.
We also have to manually check for some names because there
is a small amount of inconsistencies in the database (I don't even know how).

First argument: start id
Second argument: end id
"""
import sqlite3
import sys


def convert_to_ascii(str):
    import re
    """
    This function replaces Vietnamese accented characters with their non-accented counterparts while keeping the casing unchanged.
    Parameters:
    str (str): The input string to be converted.
    Returns:
    str: The converted string with all accented characters replaced by their non-accented counterparts, while keeping the original casing.
    """
    str = re.sub(r'[AÁÀÃẠÂẤẦẪẬĂẮẰẴẶ]', 'A', str)
    str = re.sub(r'[àáạảãâầấậẩẫăằắặẳẵ]', 'a', str)
    str = re.sub(r'[EÉÈẼẸÊẾỀỄỆ]', 'E', str)
    str = re.sub(r'[èéẹẻẽêềếệểễ]', 'e', str)
    str = re.sub(r'[IÍÌĨỊ]', 'I', str)
    str = re.sub(r'[ìíịỉĩ]', 'i', str)
    str = re.sub(r'[OÓÒÕỌÔỐỒỖỘƠỚỜỠỢ]', 'O', str)
    str = re.sub(r'[òóọỏõôồốộổỗơờớợởỡ]', 'o', str)
    str = re.sub(r'[UÚÙŨỤƯỨỪỮỰ]', 'U', str)
    str = re.sub(r'[ùúụủũưừứựửữ]', 'u', str)
    str = re.sub(r'[YÝỲỸỴ]', 'Y', str)
    str = re.sub(r'[ỳýỵỷỹ]', 'a', str)
    str = re.sub(r'[Đ]', 'D', str)
    str = re.sub(r'[đ]', 'd', str)

    # Some system encode vietnamese combining accent as individual utf-8 characters
    str = re.sub(r'\u0300|\u0301|\u0303|\u0309|\u0323', '', str)
    str = re.sub(r'\u02C6|\u0306|\u031B', '', str)
    return str


def last_name_char(name):
    return convert_to_ascii(name.split()[-1][0].upper())

con = sqlite3.connect("student.db")
start = int(sys.argv[1])
end = int(sys.argv[2])
cur = con.cursor()

just_ended = False
result = []

prev_name = 'Z'

for i in range(start, end+1):
    if just_ended:
        just_ended = False
        result.append(i)
    student = cur.execute(f"select name from students where sbd={i}").fetchone()
    if student is None:
        continue

    name = student[0]
    
    # comparator = ord(last_name_char())
    if ord(last_name_char(prev_name)) > ord(last_name_char(name)):
        result.append(i)

    prev_name = name

# ask for confirmation
confirmed_result = []

for i in result:
    print("-"*20)
    students = cur.execute("select sbd, name from students where sbd between (:critical-3) and (:critical+3);", {"critical": i}).fetchall()
    if len(students) == 0:
        continue
    for student in students:
        if student[0] == i:
            print(f"{student}<----")
            continue
        print(student)
    ignore = input(f"Ignore {i}? (y/n): ")
    if ignore == '' or ignore.lower() == 'y':
        print("Ignored")
    else:
        print("Confirmed!!!!!")
        confirmed_result.append(i)

for i in confirmed_result:
    print(i)




