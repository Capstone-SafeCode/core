import os

def delete_file(filename):
    os.remove(filename)

user_input = input("Enter the filename to delete: ")
delete_file(user_input)
