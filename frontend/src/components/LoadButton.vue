<template>
  <v-menu location="bottom" :close-on-content-click="false" v-model="menu">
    <template v-slot:activator="{ props }">
      <v-btn
        prepend-icon="mdi-folder"
        :class="btnClass"
        v-bind="props"
        :block="block"
        :variant="variant"
        >Load</v-btn
      >
    </template>
    <v-list class="py-0">
      <v-list-item class="py-0 pl-0">
        <v-btn
          prepend-icon="mdi-open-in-app"
          variant="text"
          block
          @click="loadDefault"
          class="mr-2 text-left justify-start"
          >Default</v-btn
        >
      </v-list-item>
      <v-divider v-if="list.length" />
      <template v-for="(item, i) in list" :key="item">
        <v-list-item class="py-0 pl-0">
          <v-btn
            prepend-icon="mdi-open-in-app"
            variant="text"
            block
            @click="load(item)"
            class="mr-2 text-left justify-start"
            >{{ item }}</v-btn
          >
          <template #append>
            <v-btn
              small
              class="my-1"
              variant="flat"
              @click="remove(item)"
              icon="mdi-close"
              title="remove entry"
            />
          </template>
        </v-list-item>
        <v-divider v-if="i < list.length - 1" />
      </template>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { ContextTemplate, PolicyTemplate, ResourceTemplate } from "@/assets/templates";
import { createLocalInput, removeLocalInput, getPersisted, useConfig } from "@/config";
import { ref } from "vue";

defineProps({
  btnClass: { type: String, default: "" },
  policy: { type: String, default: "" },
  resource: { type: String, default: "" },
  context: { type: String, default: "" },
  variant: { type: String },
  block: { type: Boolean, default: false },
});

const menu = ref<boolean>(false);

const list = getPersisted();

const emit = defineEmits(["update:policy", "update:resource", "update:context"]);

const config = useConfig();

const loadDefault = () => {
  emit("update:policy", PolicyTemplate);
  config.policy.value = PolicyTemplate;

  emit("update:resource", ResourceTemplate);
  config.resource.value = ResourceTemplate;

  emit("update:context", ContextTemplate);
  config.context.value = ContextTemplate;

  menu.value = false;
  config.state.value = '';
};

const load = (name: string) => {
  const input = createLocalInput(name);

  if (input.policy.value) {
    emit("update:policy", input.policy.value);
    config.policy.value = input.policy.value;
  }
  if (input.resource.value) {
    emit("update:resource", input.resource.value);
    config.resource.value = input.resource.value;
  }
  if (input.context.value) {
    emit("update:context", input.context.value);
    config.context.value = input.context.value;
  }

  menu.value = false;
  config.state.value = input.name;
};

const remove = (name: string) => removeLocalInput(name);
</script>
