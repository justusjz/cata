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
  };
  enum { CATA_LIST, CATA_SYMBOL, CATA_STRING } type;
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

void list_free(struct list *list) { free(list->nodes); }

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
  } else {
    printf("internal error: invalid node type");
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
  } else {
    printf("internal error: invalid node type");
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
    printf("error: unterminated string literal");
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

struct node parse(const char **input) {
  struct node result;
  if (**input == '(') {
    // skip (
    ++*input;
    result.type = CATA_LIST;
    result.list = parse_list(input);
    // skip )
    if (**input != ')') {
      printf("error: missing closing parenthesis");
      exit(1);
    }
    ++*input;
  } else if (**input == '"') {
    result.type = CATA_STRING;
    result.string = parse_string(input);
  } else {
    result.type = CATA_SYMBOL;
    result.symbol = parse_symbol(input);
  }
  return result;
}

int main() {
  const char *input = "(print \"Hello, world!\") (print (+ 3 4))";
  struct list list = parse_list(&input);
  node_print(&list.nodes[0]);
  return 0;
}
