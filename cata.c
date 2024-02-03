/* Copyright (c) 2024 Justus Zorn */

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define LIST_INITIAL_SIZE 8
#define LIST_GROW_FACTOR 2
#define READ_CHUNK_SIZE 4096

struct list {
  struct node *nodes;
  size_t length;
  size_t capacity;
};

struct node {
  union {
    struct list list;
    char *symbol;
    char *string;
    int integer;
  };
  enum { CATA_LIST, CATA_SYMBOL, CATA_STRING, CATA_INTEGER } type;
};

struct list list_create() {
  struct list list;
  list.nodes = NULL;
  list.length = 0;
  list.capacity = 0;
  return list;
}

void list_append(struct list *list, struct node node) {
  if (list->length == list->capacity) {
    // grow list
    if (list->capacity == 0) {
      list->capacity = LIST_INITIAL_SIZE;
    } else {
      list->capacity *= LIST_GROW_FACTOR;
    }
    list->nodes = realloc(list->nodes, list->capacity * sizeof(struct node));
  }
  // add node
  list->nodes[list->length++] = node;
}

void node_free(struct node *node);

void list_free(struct list *list) {
  for (size_t i = 0; i < list->length; ++i) {
    node_free(&list->nodes[i]);
  }
  free(list->nodes);
}

void node_print(const struct node *node) {
  if (node->type == CATA_LIST) {
    printf("(");
    for (size_t i = 0; i < node->list.length; ++i) {
      node_print(&node->list.nodes[i]);
      if (i + 1 < node->list.length) {
        // separate elements by spaces
        printf(" ");
      }
    }
    printf(")");
  } else if (node->type == CATA_SYMBOL) {
    printf("%s", node->symbol);
  } else if (node->type == CATA_STRING) {
    printf("\"%s\"", node->string);
  } else if (node->type == CATA_INTEGER) {
    printf("%d", node->integer);
  } else {
    printf("internal error: invalid node type\n");
    exit(1);
  }
}

void node_free(struct node *node) {
  if (node->type == CATA_LIST) {
    list_free(&node->list);
  } else if (node->type == CATA_SYMBOL) {
    free(node->symbol);
  } else if (node->type == CATA_STRING) {
    free(node->string);
  } else if (node->type == CATA_INTEGER) {
    // don't do anything
  } else {
    printf("internal error: invalid node type\n");
    exit(1);
  }
}

struct node parse(const char **input);

int is_ws(char c) { return c == ' ' || c == '\n' || c == '\r' || c == '\t'; }
int is_valid(char c) { return c != '\0' && c != '(' && c != ')' && !is_ws(c); }

void skip_ws(const char **input) {
  // skip all whitespace characters
  while (is_ws(**input)) {
    ++*input;
  }
}

struct list parse_list(const char **input) {
  skip_ws(input);
  struct list list = list_create();
  while (**input != ')' && **input != '\0') {
    struct node node = parse(input);
    list_append(&list, node);
    skip_ws(input);
  }
  return list;
}

char uppercase(char c) {
  if (c >= 'a' && c <= 'z') {
    return c - ('a' - 'A');
  } else {
    return c;
  }
}

char *read_sym(const char **input) {
  // save the beginning of the symbol
  const char *begin = *input;
  // skip all symbol characters
  while (is_valid(**input)) {
    ++*input;
  }
  // calculate length
  size_t length = *input - begin;
  // copy and upcase symbol
  char *sym = malloc(length + 1);
  for (size_t i = 0; i < length; ++i) {
    sym[i] = uppercase(begin[i]);
  }
  sym[length] = '\0';
  return sym;
}

char *read_str(const char **input) {
  // skip "
  ++*input;
  const char *begin = *input;
  bool escaped = false;
  // the length the result string will have
  // e.g. \\ in the source code will have only
  // length 1 in the result string
  size_t actual_length = 0;
  // skip chars until a matching " is found or EOF is reached
  // if the current char is escaped, don't break on "
  while ((**input != '"' || escaped) && **input != '\0') {
    if (!escaped) {
      // only count this character when it's not escaped
      ++actual_length;
    }
    // the next char is escaped when this is a backslash
    // if this character is escaped, the backslash is ignored
    escaped = !escaped && **input == '\\';
    if (**input == '\n') {
      printf("error: string literal cannot contain newline\n");
      exit(1);
    }
    ++*input;
  }
  if (**input != '"') {
    // reached EOF
    printf("error: unterminated string literal\n");
    exit(1);
  }
  // copy string and replace escape sequences
  char *str = malloc(actual_length + 1);
  escaped = false;
  size_t index = 0;
  for (const char *p = begin; p < *input; ++p) {
    if (escaped) {
      // not escaped anymore
      escaped = false;
      // we are in an escape sequence
      if (*p == 'n') {
        // newline escape
        str[index++] = '\n';
      } else if (*p == '"') {
        // double-quote escape
        str[index++] = '"';
      } else {
        printf("error: invalid escape sequence\n");
        exit(1);
      }
    } else {
      // check whether the next character will be escaped
      escaped = *p == '\\';
      if (!escaped) {
        // copy only when this is not an
        // escaping char
        str[index++] = *p;
      }
    }
  }
  str[actual_length] = '\0';
  // skip "
  ++*input;
  return str;
}

int to_int(const char *s, int *result) {
  *result = 0;
  size_t length = strlen(s);
  for (size_t i = 0; i < length; ++i) {
    char c = s[i];
    if (c < '0' || c > '9') {
      // not a valid integer
      return 1;
    }
    *result = *result * 10 + (s[i] - '0');
  }
  return 0;
}

struct node parse(const char **input) {
  struct node result;
  if (**input == '(') {
    // skip (
    ++*input;
    result.type = CATA_LIST;
    result.list = parse_list(input);
    // skip )
    if (**input != ')') {
      printf("error: missing closing parenthesis\n");
      exit(1);
    }
    ++*input;
  } else if (**input == '"') {
    result.type = CATA_STRING;
    result.string = read_str(input);
  } else {
    char *symbol = read_sym(input);
    if (to_int(symbol, &result.integer) == 0) {
      result.type = CATA_INTEGER;
    } else {
      result.type = CATA_SYMBOL;
      result.symbol = symbol;
    }
  }
  return result;
}

typedef union value (*native_func_t)(size_t, union value *);

union value {
  const char *string;
  int integer;
  native_func_t native_func;
};

union value native_print_string(size_t arg_count, union value *args) {
  union value result;
  printf("%s\n", args[0].string);
  result.integer = 0;
  return result;
}

union value native_print_int(size_t arg_count, union value *args) {
  union value result;
  printf("%d\n", args[0].integer);
  result.integer = 0;
  return result;
}

union value native_read_int(size_t arg_count, union value *args) {
  union value result;
  scanf("%d", &result.integer);
  return result;
}

union value native_add(size_t arg_count, union value *args) {
  union value result;
  result.integer = args[0].integer + args[1].integer;
  return result;
}

union value native_equal(size_t arg_count, union value *args) {
  union value result;
  result.integer = args[0].integer == args[1].integer;
  return result;
}

union value native_exit(size_t arg_count, union value *args) { exit(0); }

struct scope_entry {
  const char *sym;
  union value var;
};

struct scope {
  struct scope_entry *entries;
  const struct env *parent;
  size_t entry_count;
};

struct scope_entry global_scope_entries[] = {
    {.sym = "PRINT-STRING", .var.native_func = native_print_string},
    {.sym = "PRINT-INT", .var.native_func = native_print_int},
    {.sym = "READ-INT", .var.native_func = native_read_int},
    {.sym = "+", .var.native_func = native_add},
    {.sym = "=", .var.native_func = native_equal},
    {.sym = "EXIT", .var.native_func = native_exit},
};

#define ARRAY_LENGTH(arr) (sizeof(arr) / sizeof(arr[0]))

struct scope global_scope = {
    .entries = global_scope_entries, NULL, ARRAY_LENGTH(global_scope_entries)};

// find a variable in the given scope
union value *scope_find(const struct scope *scope, const char *sym) {
  // try to find the scope entry
  for (size_t i = 0; i < scope->entry_count; ++i) {
    if (strcmp(scope->entries[i].sym, sym) == 0) {
      return &scope->entries[i].var;
    }
  }
  // try the parent scope
  if (scope->parent) {
    return scope_find(scope->parent, sym);
  }
  // entry not found
  return NULL;
}

union value eval(const struct scope *scope, struct node *node) {
  union value result;
  if (node->type == CATA_INTEGER) {
    result.integer = node->integer;
  } else if (node->type == CATA_STRING) {
    result.string = node->string;
  } else if (node->type == CATA_SYMBOL) {
    const union value *var = scope_find(scope, node->symbol);
    if (!var) {
      printf("error: '%s' does not exist\n", node->symbol);
      exit(1);
    }
    result = *var;
  } else if (node->type == CATA_LIST) {
    if (node->list.length == 0) {
      printf("error: an empty list cannot be evaluated\n");
      exit(1);
    }
    struct node form = node->list.nodes[0];
    if (form.type != CATA_SYMBOL) {
      printf("error: first element of list must be symbol\n");
      exit(1);
    }
    if (strcmp(form.symbol, "IF") == 0) {
      // if special form
      if (node->list.length != 4) {
        printf("error: if needs exactly 3 arguments, but got %d\n",
               node->list.length - 1);
        exit(1);
      }
      union value condition = eval(scope, &node->list.nodes[1]);
      if (condition.integer) {
        // condition is true
        result = eval(scope, &node->list.nodes[2]);
      } else {
        // condition is false
        result = eval(scope, &node->list.nodes[3]);
      }
    } else if (strcmp(form.symbol, "LET") == 0) {
      // let special form
      if (node->list.length < 2) {
        printf("error: let needs at least 2 arguments, but got %d\n",
               node->list.length);
        exit(1);
      }
      struct node vars = node->list.nodes[1];
      if (vars.type != CATA_LIST) {
        printf("error: second argument of let must be a list\n");
        exit(1);
      }
      struct scope new_scope;
      // create a new environment
      new_scope.entry_count = vars.list.length / 2;
      new_scope.entries =
          malloc(sizeof(struct scope_entry) * new_scope.entry_count);
      new_scope.parent = scope;
      for (size_t i = 0; i + 1 < vars.list.length; i += 2) {
        struct node *var = &vars.list.nodes[i];
        struct node *expr = &vars.list.nodes[i + 1];
        if (var->type != CATA_SYMBOL) {
          printf("error: expected symbol in let\n");
          exit(1);
        }
        union value var_value = eval(scope, expr);
        struct scope_entry *entry = &new_scope.entries[i / 2];
        entry->sym = var->symbol;
        entry->var = var_value;
      }
      // evaluate the actual expressions and return the last one
      for (size_t i = 2; i < node->list.length; ++i) {
        struct node *expr = &node->list.nodes[i];
        result = eval(&new_scope, expr);
      }
      // free the environment
      free(new_scope.entries);
    } else if (strcmp(form.symbol, "FNCALL") == 0) {
      if (node->list.length < 2) {
        printf("error: fncall needs at least 1 argument\n");
        exit(1);
      }
      // evaluate the function
      union value function = eval(scope, &node->list.nodes[1]);
      // evaluate all arguments
      size_t arg_count = node->list.length - 2;
      union value *args = malloc(arg_count * sizeof(union value));
      for (size_t i = 0; i < arg_count; ++i) {
        args[i] = eval(scope, &node->list.nodes[i + 2]);
      }
      // call function
      result = function.native_func(arg_count, args);
      // free argument array
      free(args);
    } else {
      printf("error: invalid form %s\n", form.symbol);
      exit(1);
    }
  }
  return result;
}

// runs the given source code
void run(const char *source) {
  // parse the source as multiple objects
  struct list list = parse_list(&source);
  for (size_t i = 0; i < list.length; ++i) {
    eval(&global_scope, &list.nodes[i]);
  }
  list_free(&list);
}

char *read_file(const char *filename) {
  // try to open the file
  FILE *f = fopen(filename, "r");
  if (!f) {
    return NULL;
  }
  // allocate a buffer
  char *buffer = malloc(READ_CHUNK_SIZE);
  size_t read = 0, total = 0, offset = 0;
  while ((read = fread(buffer + total, 1, READ_CHUNK_SIZE, f)) ==
         READ_CHUNK_SIZE) {
    // when the buffer was filled completely, reallocate it
    total += read;
    buffer = realloc(buffer, total + READ_CHUNK_SIZE);
  }
  // count the total number of bytes written
  total += read;
  // add '\0' to buffer
  buffer = realloc(buffer, total + 1);
  buffer[total] = '\0';
  return buffer;
}

int main(int argc, char *argv[]) {
  if (argc < 2) {
    printf("usage: cata <file>\n");
    exit(0);
  }
  for (int i = 1; i < argc; ++i) {
    const char *filename = argv[i];
    char *content = read_file(filename);
    if (!content) {
      printf("error: could not read file %s\n", filename);
      continue;
    }
    run(content);
    free(content);
  }
  return 0;
}
