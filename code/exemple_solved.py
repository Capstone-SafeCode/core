import os

SAFE_DIRECTORY = "/src"

def delete_file(filename):
  filepath = os.path.join(SAFE_DIRECTORY, filename)

  filepath = os.path.normpath(filepath)

  if not filepath.startswith(SAFE_DIRECTORY):
    print("Unauthorized file path")
    return

  os.remove(filepath)

user_input = input("Enter the filename to delete: ")
delete_file(user_input)
