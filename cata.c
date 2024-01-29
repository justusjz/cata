/* Copyright (c) 2024 Justus Zorn */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define LIST_INITIAL_SIZE 8
#define LIST_GROW_FACTOR 2

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

struct env_entry {
  const char *sym;
  native_func func;
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

union value native_add(size_t arg_count, union value *args) {
  union value result;
  result.integer = args[0].integer + args[1].integer;
  return result;
}

union value native_exit(size_t arg_count, union value *args) { exit(0); }

struct env_entry env[] = {
    {"print-string", native_print_string},
    {"print-int", native_print_int},
    {"+", native_add},
    {"exit", native_exit},
    {NULL, NULL},
};

native_func env_find(const char *sym) {
  // try to find the function
  for (size_t i = 0; env[i].sym != NULL; ++i) {
    if (strcmp(env[i].sym, sym) == 0) {
      return env[i].func;
    }
  }
  // function does not exist
  return NULL;
}

union value eval(struct node *node) {
  union value result;
  if (node->type == CATA_INTEGER) {
    result.integer = node->integer;
  } else if (node->type == CATA_STRING) {
    result.string = node->string;
  } else if (node->type == CATA_SYMBOL) {
    printf("error: variables are not yet supported\n");
    exit(1);
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
    // evaluate all arguments
    size_t arg_count = node->list.length - 1;
    union value *args = malloc(arg_count * sizeof(union value));
    for (size_t i = 0; i < arg_count; ++i) {
      args[i] = eval(&node->list.nodes[i + 1]);
    }
    // find function
    native_func func = env_find(fn.symbol);
    if (func == NULL) {
      printf("error: function %s does not exist\n", fn.symbol);
      exit(1);
    }
    // call function
    result = func(arg_count, args);
    // free argument array
    free(args);
  }
  return result;
}

void run(const char *source) {
  struct list list = parse_list(&source);
  for (size_t i = 0; i < list.length; ++i) {
    eval(&list.nodes[i]);
  }
  list_free(&list);
}

int main() {
  const char *source = "(print-string \"Hello, world!\")";
  run(source);
  return 0;
}
