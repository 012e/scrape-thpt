import json
import re

def convert_string_to_json(input_string):
    data = {}

    # Use regular expression to find subject names and grades
    pattern = r'(\w+(?: \w+)*):\s+(\d+\.\d+)'
    matches = re.findall(pattern, input_string)

    # Populate the data dictionary
    for subject, grade in matches:
        data[subject] = float(grade)

    return data

# Input string containing subject names and grades
input_string = "Toán: 5.20 Ngữ văn: 5.50 Lịch sử: 3.50 Địa lí: 5.00 GDCD: 5.75 KHXH: 4.75 Tiếng Anh: 3.60"

# Call the function and print the result
print(convert_string_to_json(input_string))
