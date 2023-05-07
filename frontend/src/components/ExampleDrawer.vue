<template>
  <v-navigation-drawer
    :model-value="props.modelValue"
    temporary
    width="400"
    @update:modelValue="(event: boolean) => emit('update:modelValue', event)"
  >
    <template v-for="(value, group) in config.examples" :key="group">
      <v-list>
        <v-list-group :value="group">
          <template v-slot:activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-folder"
              :title="group"
            ></v-list-item>
          </template>

          <v-list-item
            v-for="policy in value.policies"
            :key="policy"
            :title="policy.replaceAll('-', ' ')"
            @click="() => loadExample(value.path, policy)"
          ></v-list-item>
        </v-list-group>
      </v-list>
    </template>

    <v-overlay :model-value="overlay" class="align-center justify-center">
      <v-progress-circular color="primary" indeterminate size="64"></v-progress-circular>
    </v-overlay>

    <template v-slot:append>
      <div class="pa-2">
        <v-btn flat color="primary" block @click="() => emit('update:modelValue', false)"
          >Close</v-btn
        >
      </div>
    </template>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { config } from "../config";

const props = defineProps({
  modelValue: { type: Boolean, default: false },
});

const emit = defineEmits(["update:modelValue", "select:example"]);
const overlay = ref<boolean>(false);

const loadExample = async (path: string, name: string) => {
  try {
    overlay.value = true;
    const values = await Promise.all([
      fetch(`${path}/${name}/${name}.yaml`).then((resp) => resp.text()),
      fetch(`${path}/${name}/resource.yaml`).then((resp) => {
        if (resp.status === 404) {
          return fetch(`${path}/${name}/resources.yaml`)
        }

        return resp
      }).then((resp) => resp.text()),
    ]);

    emit("select:example", values);
    emit("update:modelValue", false);
  } catch (err) {
    console.error(err);
  } finally {
    overlay.value = false;
  }
};
</script>
