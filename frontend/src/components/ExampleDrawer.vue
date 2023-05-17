<template>
  <v-navigation-drawer
    :model-value="props.modelValue"
    temporary
    width="400"
    @update:modelValue="(event: boolean) => emit('update:modelValue', event)"
  >
    <template v-for="({ color, ...example }) in options.examples" :key="example.name">
      <v-list>
        <v-list-group :value="example.name" :id="example.name.toLowerCase().replaceAll(' ', '-')">
          <template v-slot:activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-folder"
              :title="example.name"
              :class="color ? `text-${color}` : ''"
            ></v-list-item>
          </template>

          <template  v-if="example.subgroups">
            <template v-for="(subgroup, j) in example.subgroups" :key="j">
              <v-list-subheader>{{ subgroup.name }}</v-list-subheader>
              <v-list-item
                v-for="(policy, i) in subgroup.policies"
                :key="i"
                :title="policy.title"
                @click="() => loadExample(policy.url || subgroup.url || example.url, policy)"
              />
            </template>
          </template>

          <template v-else>
            <v-list-item
              v-for="(policy, i) in example.policies"
              :key="i"
              :title="policy.title"
              @click="() => loadExample(example.url, policy)"
            />
          </template>
        </v-list-group>
      </v-list>
    </template>

    <v-overlay :model-value="overlay" class="align-center justify-center">
      <v-progress-circular color="primary" indeterminate size="64"></v-progress-circular>
    </v-overlay>

    <template v-slot:append>
      <div class="pa-2">
        <v-btn flat color="primary" block @click="() => emit('update:modelValue', false)">Close</v-btn>
      </div>
    </template>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useConfig } from "../config";
import { Policy, loadPolicy } from "@/functions/github";
import { init } from "@/store";
const { options } = useConfig()

const props = defineProps({
  modelValue: { type: Boolean, default: false },
});

const emit = defineEmits(["update:modelValue", "select:example"]);
const overlay = ref<boolean>(false);


const loadExample = async (url: string, policy: Policy) => {
  loadPolicy(url, policy, {
    start: () => { overlay.value = true },
    success: ([policy, resource, context]) => {
      init({ policy, resource, context })
      emit("update:modelValue", false);
    },
    error: (err) => console.error(err),
    finished: () => { overlay.value = false; }
  })
};
</script>
