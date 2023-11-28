// Copyright (c) 2023 Justus Zorn

#include <stdint.h>
#include <stdio.h>

struct u8_slice {
    uint8_t *data;
    size_t length;
};

void print(struct u8_slice s) {
    printf("%.*s", s.length, s.data);
}
