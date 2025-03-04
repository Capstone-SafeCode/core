import ast
import json
import sys

def ast_to_json(node):
  """ Convertit un nœud AST en un dictionnaire JSON récursif. """
  if isinstance(node, ast.AST):
      result = {"_type": type(node).__name__}

      if hasattr(node, "lineno"):
        result["lineno"] = node.lineno
      if hasattr(node, "col_offset"):
          result["col_offset"] = node.col_offset

      for field, value in ast.iter_fields(node):
          result[field] = ast_to_json(value)
      return result
  elif isinstance(node, list):
      return [ast_to_json(item) for item in node]
  else:
      return node

def generate_ast(filename):
  """ Génère un AST pour le fichier Python spécifié. """
  f = open(filename, "r")
  tree = ast.parse(f.read())
  return ast_to_json(tree)

if __name__ == "__main__":
  if len(sys.argv) != 2:
    print("Usage: python ast.py <filename>")
    sys.exit(1)

  filename = sys.argv[1]
  ast_json = generate_ast(filename)
  print(json.dumps(ast_json, indent=2))
