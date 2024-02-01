/* Copyright (c) 2024 Justus Zorn */

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

char *parse_symbol(const char **input) {
  // save the beginning of the symbol
  const char *begin = *input;
  // skip all symbol characters
  while (is_valid(**input)) {
    ++*input;
  }
  // calculate length
  size_t length = *input - begin;
  // copy string
  char *symbol = malloc(length + 1);
  strncpy(symbol, begin, length);
  symbol[length] = '\0';
  return symbol;
}

char *parse_string(const char **input) {
  // skip "
  ++*input;
  const char *begin = *input;
  while (**input != '"' && **input != '\0') {
    ++*input;
  }
  if (**input != '"') {
    printf("error: unterminated string literal\n");
    exit(1);
  }
  size_t length = *input - begin;
  char *string = malloc(length + 1);
  strncpy(string, begin, length);
  string[length] = '\0';
  // skip "
  ++*input;
  return string;
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
    result.string = parse_string(input);
  } else {
    char *symbol = parse_symbol(input);
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
    {.sym = "print-string", .var.native_func = native_print_string},
    {.sym = "print-int", .var.native_func = native_print_int},
    {.sym = "read-int", .var.native_func = native_read_int},
    {.sym = "+", .var.native_func = native_add},
    {.sym = "=", .var.native_func = native_equal},
    {.sym = "exit", .var.native_func = native_exit},
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
    if (strcmp(form.symbol, "if") == 0) {
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
    } else if (strcmp(form.symbol, "let") == 0) {
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
    } else if (strcmp(form.symbol, "fncall") == 0) {
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

void run(const char *source) {
  struct list list = parse_list(&source);
  for (size_t i = 0; i < list.length; ++i) {
    eval(&global_scope, &list.nodes[i]);
  }
  list_free(&list);
}

char *read_file(const char *filename) {
  FILE *f = fopen(filename, "r");
  if (!f) {
    return NULL;
  }
  char *buffer = malloc(READ_CHUNK_SIZE);
  size_t read = 0, total = 0, offset = 0;
  while ((read = fread(buffer + total, 1, READ_CHUNK_SIZE, f)) ==
         READ_CHUNK_SIZE) {
    total += read;
    buffer = realloc(buffer, total + READ_CHUNK_SIZE);
  }
  total += read;
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
