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

union value {
  const char *string;
  int integer;
};

typedef union value (*native_func)(size_t, union value *);

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

enum env_entry_type {
  ENV_ENTRY_FUNC,
  ENV_ENTRY_VAR,
};

struct env_entry {
  const char *sym;
  union {
    native_func func;
    union value var;
  };
  enum env_entry_type type;
};

struct env {
  struct env_entry *entries;
  const struct env *parent;
  size_t entry_count;
};

struct env_entry global_env_entries[] = {
    {.sym = "print-string",
     .func = native_print_string,
     .type = ENV_ENTRY_FUNC},
    {.sym = "print-int", .func = native_print_int, .type = ENV_ENTRY_FUNC},
    {.sym = "read-int", .func = native_read_int, .type = ENV_ENTRY_FUNC},
    {.sym = "+", .func = native_add, .type = ENV_ENTRY_FUNC},
    {.sym = "=", .func = native_equal, .type = ENV_ENTRY_FUNC},
    {.sym = "exit", .func = native_exit, .type = ENV_ENTRY_FUNC},
};

#define ARRAY_LENGTH(arr) (sizeof(arr) / sizeof(arr[0]))

struct env global_env = {global_env_entries, NULL,
                         ARRAY_LENGTH(global_env_entries)};

const struct env_entry *env_find(const struct env *env, const char *sym) {
  // try to find the env entry
  for (size_t i = 0; i < env->entry_count; ++i) {
    if (strcmp(env->entries[i].sym, sym) == 0) {
      return &env->entries[i];
    }
  }
  // entry not found
  return NULL;
}

union value eval(const struct env *env, struct node *node) {
  union value result;
  if (node->type == CATA_INTEGER) {
    result.integer = node->integer;
  } else if (node->type == CATA_STRING) {
    result.string = node->string;
  } else if (node->type == CATA_SYMBOL) {
    const struct env_entry *entry = env_find(env, node->symbol);
    if (entry->type != ENV_ENTRY_VAR) {
      printf("error: %s is not a variable\n", node->symbol);
      exit(1);
    }
    result = entry->var;
  } else if (node->type == CATA_LIST) {
    if (node->list.length == 0) {
      printf("error: empty list is invalid\n");
      exit(1);
    }
    struct node fn = node->list.nodes[0];
    if (fn.type != CATA_SYMBOL) {
      printf("%d\n", fn.type);
      printf("error: only symbol can be called as a function\n");
      exit(1);
    }
    if (strcmp(fn.symbol, "if") == 0) {
      // if special form
      if (node->list.length != 4) {
        printf("error: if needs exactly 3 arguments, but got %d\n",
               node->list.length - 1);
        exit(1);
      }
      union value condition = eval(env, &node->list.nodes[1]);
      if (condition.integer) {
        // condition is true
        result = eval(env, &node->list.nodes[2]);
      } else {
        // condition is false
        result = eval(env, &node->list.nodes[3]);
      }
    } else if (strcmp(fn.symbol, "let") == 0) {
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
      struct env new_env;
      // create a new environment
      new_env.entry_count = vars.list.length / 2;
      new_env.entries = malloc(sizeof(struct env_entry) * new_env.entry_count);
      new_env.parent = env;
      for (size_t i = 0; i + 1 < vars.list.length; i += 2) {
        struct node *var = &vars.list.nodes[i];
        struct node *expr = &vars.list.nodes[i + 1];
        if (var->type != CATA_SYMBOL) {
          printf("error: expected symbol in let\n");
          exit(1);
        }
        union value var_value = eval(env, expr);
        struct env_entry *entry = &new_env.entries[i / 2];
        entry->type = ENV_ENTRY_VAR;
        entry->sym = var->symbol;
        entry->var = var_value;
      }
      // evaluate the actual expressions and return the last one
      for (size_t i = 2; i < node->list.length; ++i) {
        struct node *expr = &node->list.nodes[i];
        result = eval(&new_env, expr);
      }
      // free the environment
      free(new_env.entries);
    } else {
      // evaluate all arguments
      size_t arg_count = node->list.length - 1;
      union value *args = malloc(arg_count * sizeof(union value));
      for (size_t i = 0; i < arg_count; ++i) {
        args[i] = eval(env, &node->list.nodes[i + 1]);
      }
      // find function
      const struct env_entry *entry = env_find(env, fn.symbol);
      if (entry == NULL) {
        printf("error: function %s does not exist\n", fn.symbol);
        exit(1);
      } else if (entry->type != ENV_ENTRY_FUNC) {
        printf("error: %s is not a function\n", fn.symbol);
        exit(1);
      }
      // call function
      result = entry->func(arg_count, args);
      // free argument array
      free(args);
    }
  }
  return result;
}

void run(const char *source) {
  struct list list = parse_list(&source);
  for (size_t i = 0; i < list.length; ++i) {
    eval(&global_env, &list.nodes[i]);
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
