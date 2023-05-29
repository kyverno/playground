<template>
  <v-menu open-on-hover class="py-0">
    <template v-slot:activator="{ props }">
      <v-btn color="secondary" v-bind="props">Legend</v-btn>
    </template>

    <v-list class="my-0 py-0">
      <v-list-item v-for="(item, index) in items" :key="index" class="border">
        <template v-slot:prepend>
          <div style="width: 130px">
            {{ item.shortcut }}
          </div>
        </template>

        <v-list-item-title class="ml-2">{{ item.description }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
type Command = { shortcut: string; description: string }

const items: Command[] = [
  { shortcut: 'OPTION/ALT + Z', description: 'Toggle Word Wrap' },
  { shortcut: 'CMD/CTRL + S', description: 'Save changes' },
  { shortcut: 'ESCAPE', description: 'Discard changes and close dialog' }
]
</script>
